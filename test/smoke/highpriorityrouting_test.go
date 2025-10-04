package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

func TestHighPriorityRouting(t *testing.T) {
	t.Parallel()

	const sql = `
    insert into users (id, name, email)
    values ({{uuid "u1"}}, 'bob', 'joe');
    insert into user_contact_methods (id, user_id, name, type, value)
    values
        ({{uuid "cms"}}, {{uuid "u1"}}, 'sms', 'SMS', {{phone "1"}}),
        ({{uuid "cmv"}}, {{uuid "u1"}}, 'voice', 'VOICE', {{phone "1"}});
    insert into user_notification_rules (user_id, contact_method_id, delay_minutes)
    values ({{uuid "u1"}}, {{uuid "cms"}}, 0);
    insert into escalation_policies (id, name)
    values ({{uuid "ep"}}, 'esc');
    insert into escalation_policy_steps (id, escalation_policy_id)
    values ({{uuid "step"}}, {{uuid "ep"}});
    insert into escalation_policy_actions (escalation_policy_step_id, user_id)
    values ({{uuid "step"}}, {{uuid "u1"}});
    insert into services (id, escalation_policy_id, name)
    values ({{uuid "svc"}}, {{uuid "ep"}}, 'svc');
    `

	h := harness.NewHarness(t, sql, "ids-to-uuids")
	defer h.Close()

	h.SetConfigValue("Alerts.HighPriorityLabelKey", "alerts/priority")
	h.SetConfigValue("Alerts.HighPriorityLabelValue", "high")

	dev := h.Twilio(t).Device(h.Phone("1"))

	h.CreateAlert(h.UUID("svc"), "normal alert")
	dev.ExpectSMS("normal alert")

	h.GraphQLQuery2(`mutation{createAlert(input:{serviceID:"` + h.UUID("svc") + `", summary:"priority alert", meta:[{key:"alerts/priority", value:"high"}]}){id}}`)
	dev.ExpectVoice("priority alert")
}
