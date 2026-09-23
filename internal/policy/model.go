package policy

type Action string

const (
	ActionNone             Action = "none"
	ActionRestartContainer Action = "restart_container"
)

type Recommendation struct {
	Action                Action
	RequiresHumanApproval bool
}
