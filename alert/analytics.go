package alert

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/target/goalert/permission"
	"github.com/target/goalert/validation"
	"github.com/target/goalert/validation/validate"
)

// MetaPriorityCount represents the number of alerts for a given metadata value.
type MetaPriorityCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// MetaAckLevelCount represents the number of acknowledgements at a given escalation level.
type MetaAckLevelCount struct {
	EscalationLevel int `json:"escalationLevel"`
	Count           int `json:"count"`
}

// MetaAnalytics aggregates alert priority counts and acknowledgement distribution for a service.
type MetaAnalytics struct {
	PriorityCounts      []MetaPriorityCount `json:"priorityCounts"`
	AcknowledgedByLevel []MetaAckLevelCount `json:"acknowledgedByLevel"`
}

const metaPrioritySQL = `
WITH priority_alerts AS (
    SELECT
        jsonb_extract_path_text(ad.metadata, 'alertMetaV1', $4) AS meta_value
    FROM alerts a
    LEFT JOIN alert_data ad ON ad.alert_id = a.id
    WHERE a.service_id = $1
      AND a.created_at >= $2
      AND a.created_at < $3
)
SELECT meta_value, COUNT(*)
FROM priority_alerts
WHERE meta_value IS NOT NULL AND meta_value <> ''
GROUP BY meta_value
ORDER BY meta_value;
`

const metaAckSQL = `
WITH acked AS (
    SELECT
        a.escalation_level
    FROM alert_logs l
    JOIN alerts a ON a.id = l.alert_id
    LEFT JOIN alert_data ad ON ad.alert_id = a.id
    WHERE a.service_id = $1
      AND l.event = 'acknowledged'
      AND l.timestamp >= $2
      AND l.timestamp < $3
      AND jsonb_extract_path_text(ad.metadata, 'alertMetaV1', $4) IS NOT NULL
      AND jsonb_extract_path_text(ad.metadata, 'alertMetaV1', $4) <> ''
)
SELECT escalation_level, COUNT(*)
FROM acked
GROUP BY escalation_level
ORDER BY escalation_level;
`

// MetaKeyAnalytics returns aggregated statistics for alerts that have the provided metadata key.
func (s *Store) MetaKeyAnalytics(ctx context.Context, serviceID, metaKey string, start, end time.Time) (*MetaAnalytics, error) {
	if err := permission.LimitCheckAny(ctx, permission.Admin); err != nil {
		return nil, err
	}

	if err := validate.UUID("ServiceID", serviceID); err != nil {
		return nil, err
	}

	if metaKey == "" {
		return nil, validation.NewFieldError("MetaKey", "must be provided")
	}

	if !start.Before(end) {
		return nil, validation.NewFieldError("TimeRange", "start must be before end")
	}

	svcID := uuid.MustParse(serviceID)

	res := &MetaAnalytics{}

	rows, err := s.db.QueryContext(ctx, metaPrioritySQL, svcID, start, end, metaKey)
	if err != nil {
		return nil, errors.Wrap(err, "query metadata priorities")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			value sql.NullString
			count int
		)
		if err := rows.Scan(&value, &count); err != nil {
			return nil, errors.Wrap(err, "scan metadata priority row")
		}
		res.PriorityCounts = append(res.PriorityCounts, MetaPriorityCount{
			Value: value.String,
			Count: count,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate metadata priorities")
	}

	rows, err = s.db.QueryContext(ctx, metaAckSQL, svcID, start, end, metaKey)
	if err != nil {
		return nil, errors.Wrap(err, "query acknowledgement distribution")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			level sql.NullInt64
			count int
		)
		if err := rows.Scan(&level, &count); err != nil {
			return nil, errors.Wrap(err, "scan acknowledgement row")
		}
		lvl := -1
		if level.Valid {
			lvl = int(level.Int64)
		}
		res.AcknowledgedByLevel = append(res.AcknowledgedByLevel, MetaAckLevelCount{
			EscalationLevel: lvl,
			Count:           count,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate acknowledgement rows")
	}

	return res, nil
}
