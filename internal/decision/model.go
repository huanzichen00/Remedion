package decision

type Cause string

const (
	CauseResourceExhaustion Cause = "resource_exhaustion"
	CauseApplicationFailure Cause = "application_failure"
	CauseUnknown            Cause = "unknown"
)

type Decision struct {
	Cause Cause

	CauseConfidence float64

	RestartHelpfulProbability   float64
	NeedsHumanReviewProbability float64
}
