package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/target/goalert/alert/alertlog"
	"github.com/target/goalert/engine/message"
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/twilio"
	"github.com/target/goalert/permission"
	"github.com/target/goalert/user/contactmethod"
	"github.com/target/goalert/util/log"
)

type contactMethodFinder interface {
	FindAll(ctx context.Context, dbtx gadb.DBTX, userID string) ([]contactmethod.ContactMethod, error)
}

// applyHighPriorityOverride returns true when the current notification should be suppressed
// (e.g., a non-priority alert targeting a voice contact method).
func applyHighPriorityOverride(ctx context.Context, db gadb.DBTX, store contactMethodFinder, msg *message.Message, meta map[string]string, key, val string) bool {
	if key == "" || val == "" {
		return false
	}

	isHigh := meta[key] == val
	log.Logf(ctx, "high-priority override: evaluating meta[%q]=%q (required=%q)", key, meta[key], val)
	if isHigh {
		log.Logf(ctx, "high-priority override: alert tagged high priority; ensuring voice delivery")
		// Already targeting voice, nothing to promote.
		if msg.Dest.Type == twilio.DestTypeTwilioVoice {
			log.Logf(ctx, "high-priority override: already targeting voice; no change")
			return false
		}

		// Promote the notification to voice if the user has a voice contact method.
		log.Logf(ctx, "high-priority override: scanning contact methods for voice destination")
		cms, err := store.FindAll(ctx, db, msg.UserID)
		if err != nil {
			log.Log(ctx, errors.Wrap(err, "lookup contact methods for high priority alert"))
			return false
		}
		for i := range cms {
			cm := &cms[i]
			if cm.Dest.Type != twilio.DestTypeTwilioVoice {
				continue
			}
			msg.Dest = cm.Dest
			msg.DestID = notification.DestID{CMID: uuid.NullUUID{UUID: cm.ID, Valid: true}}
			log.Logf(ctx, "high-priority override: promoted to voice contact method %s", cm.ID)
			return false
		}
		log.Logf(ctx, "high-priority override: no voice contact method found; keep existing destination")
		return false
	}

	// Not high priority: suppress voice notifications and allow subsequent rules to run.
	if msg.Dest.Type == twilio.DestTypeTwilioVoice {
		log.Logf(ctx, "high-priority override: alert not high priority; suppressing voice notification for user %s", msg.UserID)
		return true
	}

	return false
}

func (p *Engine) sendMessage(ctx context.Context, msg *message.Message) (*notification.SendResult, error) {
	ctx = log.WithField(ctx, "CallbackID", msg.ID)
	log.Logf(ctx, "sendMessage: start processing message type=%s destType=%s", msg.Type, msg.Dest.Type)

	if msg.DestID.IsUserCM() {
		ctx = permission.UserSourceContext(ctx, msg.UserID, permission.RoleUser, &permission.SourceInfo{
			Type: permission.SourceTypeContactMethod,
			ID:   msg.DestID.String(),
		})
		log.Logf(ctx, "sendMessage: permission context set for contact method %s", msg.DestID.String())
	} else {
		ctx = permission.SystemContext(ctx, "SendMessage")
		ctx = permission.SourceContext(ctx, &permission.SourceInfo{
			Type: permission.SourceTypeNotificationChannel,
			ID:   msg.DestID.String(),
		})
		log.Logf(ctx, "sendMessage: permission context set for notification channel %s", msg.DestID.String())
	}

	var notifMsg notification.Message
	var isFirstAlertMessage bool
	switch msg.Type {
	case notification.MessageTypeAlertBundle:
		log.Logf(ctx, "sendMessage: building alert bundle payload for service=%s", msg.ServiceID)
		name, count, err := p.a.ServiceInfo(ctx, msg.ServiceID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup service info")
		}
		if count == 0 {
			log.Logf(ctx, "sendMessage: bundle resolved to zero alerts; skipping send")
			// already acked/closed, don't send bundled notification
			return &notification.SendResult{
				ID: msg.ID,
				Status: notification.Status{
					Details: "alerts acked/closed before message sent",
					State:   notification.StateFailedPerm,
				},
			}, nil
		}
		notifMsg = notification.AlertBundle{
			Base:        msg.Base(),
			ServiceID:   msg.ServiceID,
			ServiceName: name,
			Count:       count,
		}
	case notification.MessageTypeAlert:
		log.Logf(ctx, "sendMessage: building alert payload alertID=%d", msg.AlertID)
		name, _, err := p.a.ServiceInfo(ctx, msg.ServiceID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup service info")
		}
		a, err := p.a.FindOne(ctx, msg.AlertID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup alert")
		}
		meta, err := p.a.Metadata(ctx, p.b.db, msg.AlertID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup alert metadata")
		}
		log.Logf(ctx, "sendMessage: fetched alert metadata keys=%d", len(meta))
		suppress := applyHighPriorityOverride(ctx, p.b.db, p.cfg.ContactMethodStore, msg, meta, p.cfg.ConfigSource.Config().Alerts.HighPriorityLabelKey, p.cfg.ConfigSource.Config().Alerts.HighPriorityLabelValue)
		if suppress {
			log.Logf(ctx, "sendMessage: voice notification suppressed by high-priority override")
			return &notification.SendResult{
				ID: msg.ID,
				Status: notification.Status{
					Details: "voice notification suppressed for non-priority alert",
					State:   notification.StateFailedPerm,
				},
			}, nil
		}
		stat, err := p.cfg.NotificationStore.OriginalMessageStatus(ctx, msg.AlertID, msg.DestID)
		if err != nil {
			return nil, fmt.Errorf("lookup original message: %w", err)
		}
		if stat != nil && stat.ID == msg.ID {
			log.Logf(ctx, "sendMessage: original message status matched current message; clearing stat reference")
			// set to nil if it's the current message
			stat = nil
		}
		notifMsg = notification.Alert{
			Base:        msg.Base(),
			AlertID:     msg.AlertID,
			Summary:     a.Summary,
			Details:     a.Details,
			ServiceID:   a.ServiceID,
			ServiceName: name,
			Meta:        meta,

			OriginalStatus: stat,
		}
		isFirstAlertMessage = stat == nil
	case notification.MessageTypeAlertStatus:
		log.Logf(ctx, "sendMessage: building alert status payload logID=%d", msg.AlertLogID)
		e, err := p.cfg.AlertLogStore.FindOne(ctx, msg.AlertLogID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup alert log entry")
		}
		a, err := p.cfg.AlertStore.FindOne(ctx, msg.AlertID)
		if err != nil {
			return nil, fmt.Errorf("lookup original alert: %w", err)
		}
		stat, err := p.cfg.NotificationStore.OriginalMessageStatus(ctx, msg.AlertID, msg.DestID)
		if err != nil {
			return nil, fmt.Errorf("lookup original message: %w", err)
		}
		if stat == nil {
			return nil, fmt.Errorf("could not find original notification for alert %d to %s", msg.AlertID, msg.Dest.String())
		}

		var status notification.AlertState
		switch e.Type() {
		case alertlog.TypeAcknowledged:
			status = notification.AlertStateAcknowledged
		case alertlog.TypeEscalated:
			status = notification.AlertStateUnacknowledged
		case alertlog.TypeClosed:
			status = notification.AlertStateClosed
		}

		notifMsg = notification.AlertStatus{
			Base:           msg.Base(),
			AlertID:        e.AlertID(),
			ServiceID:      a.ServiceID,
			LogEntry:       e.String(ctx),
			Summary:        a.Summary,
			Details:        a.Details,
			NewAlertState:  status,
			OriginalStatus: *stat,
		}
	case notification.MessageTypeTest:
		log.Logf(ctx, "sendMessage: building test notification payload")
		notifMsg = notification.Test{
			Base: msg.Base(),
		}
	case notification.MessageTypeVerification:
		log.Logf(ctx, "sendMessage: building verification payload verifyID=%s", msg.VerifyID)
		code, err := p.cfg.NotificationStore.Code(ctx, msg.VerifyID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup verification code")
		}
		notifMsg = notification.Verification{
			Base: msg.Base(),
			Code: fmt.Sprintf("%06d", code),
		}
	case notification.MessageTypeScheduleOnCallUsers:
		log.Logf(ctx, "sendMessage: building on-call schedule payload scheduleID=%s", msg.ScheduleID)
		users, err := p.cfg.OnCallStore.OnCallUsersBySchedule(ctx, msg.ScheduleID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup on call users by schedule")
		}
		sched, err := p.cfg.ScheduleStore.FindOne(ctx, msg.ScheduleID)
		if err != nil {
			return nil, errors.Wrap(err, "lookup schedule by id")
		}

		var onCallUsers []notification.User
		for _, u := range users {
			onCallUsers = append(onCallUsers, notification.User{
				Name: u.Name,
				ID:   u.ID,
				URL:  p.cfg.ConfigSource.Config().CallbackURL("/users/" + u.ID),
			})
		}

		notifMsg = notification.ScheduleOnCallUsers{
			Base:         msg.Base(),
			ScheduleName: sched.Name,
			ScheduleURL:  p.cfg.ConfigSource.Config().CallbackURL("/schedules/" + msg.ScheduleID),
			ScheduleID:   msg.ScheduleID,
			Users:        onCallUsers,
		}
	case notification.MessageTypeSignalMessage:
		log.Logf(ctx, "sendMessage: building signal payload messageID=%s", msg.ID)
		id, err := uuid.Parse(msg.ID)
		if err != nil {
			return nil, errors.Wrap(err, "parse signal message id")
		}
		rawParams, err := gadb.New(p.b.db).EngineGetSignalParams(ctx, uuid.NullUUID{Valid: true, UUID: id})
		if err != nil {
			return nil, errors.Wrap(err, "get signal message params")
		}

		var params map[string]string
		err = json.Unmarshal(rawParams, &params)
		if err != nil {
			return nil, errors.Wrap(err, "parse signal message params")
		}

		notifMsg = notification.SignalMessage{
			Base:   msg.Base(),
			Params: params,
		}
	default:
		log.Log(ctx, errors.New("SEND NOT IMPLEMENTED FOR MESSAGE TYPE "+string(msg.Type)))
		return &notification.SendResult{ID: msg.ID, Status: notification.Status{State: notification.StateFailedPerm}}, nil
	}

	meta := alertlog.NotificationMetaData{
		MessageID: msg.ID,
	}
	log.Logf(ctx, "sendMessage: dispatching via notification manager destType=%s", notifMsg.DestType())

	res, err := p.cfg.NotificationManager.SendMessage(ctx, notifMsg)
	if err != nil {
		return nil, err
	}
	log.Logf(ctx, "sendMessage: provider result state=%s details=%s", res.State, res.Details)

	switch msg.Type {
	case notification.MessageTypeAlert:
		log.Logf(ctx, "sendMessage: recording alert log entry for alertID=%d", msg.AlertID)
		p.cfg.AlertLogStore.MustLog(ctx, msg.AlertID, alertlog.TypeNotificationSent, meta)
	case notification.MessageTypeAlertBundle:
		err = p.cfg.AlertLogStore.LogServiceTx(ctx, nil, msg.ServiceID, alertlog.TypeNotificationSent, meta)
		if err != nil {
			log.Log(ctx, errors.Wrap(err, "append alert log"))
		}
	}

	if isFirstAlertMessage && res.State.IsOK() {
		log.Logf(ctx, "sendMessage: tracking status for alertID=%d dest=%s", msg.AlertID, msg.Dest.String())
		_, err = p.b.trackStatus.ExecContext(ctx, msg.DestID.NCID, msg.DestID.CMID, msg.AlertID)
		if err != nil {
			// non-fatal, but log because it means status updates will not work for that alert/dest.
			log.Log(ctx, fmt.Errorf("track status updates for alert #%d for %s: %w", msg.AlertID, msg.Dest.String(), err))
		}
	}

	log.Logf(ctx, "sendMessage: completed processing for messageID=%s", msg.ID)
	return res, nil
}
