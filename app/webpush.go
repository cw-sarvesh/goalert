package app

import "github.com/target/goalert/oncall"

// filterPrimaryStepUserIDs returns only the user IDs belonging to the
// escalation step with the lowest step number present.
func filterPrimaryStepUserIDs(users []oncall.ServiceOnCallUser) ([]string, int) {
	minStep := -1
	result := make([]string, 0, len(users))

	for _, u := range users {
		if minStep == -1 || u.StepNumber < minStep {
			minStep = u.StepNumber
			result = result[:0]
		}
		if u.StepNumber == minStep {
			result = append(result, u.UserID)
		}
	}

	return append([]string(nil), result...), minStep
}
