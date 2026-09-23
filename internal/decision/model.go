package decision

type Cause string

const (
	CauseResourceExhaustion Cause = "resource_exhaustion"
	CauseApplicationFailure Cause = "application_failue"
	CauseUnknown            Cause = "unknown"
)

type Decision struct {
	Cause            Cause
	RestartHelpful   bool
	NeedsHumanReview bool
}
