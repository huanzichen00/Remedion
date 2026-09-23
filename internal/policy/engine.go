package policy

import "github.com/huanzichen00/remedion/internal/decision"

type Engine struct {
	restartThreshold float64
}

func NewEngine(restartThreshold float64) *Engine {
	return &Engine{
		restartThreshold: restartThreshold,
	}
}

func (e *Engine) Evaluate(dec decision.Decision) Recommendation {
	if dec.RestartHelpfulProbability < e.restartThreshold {
		return Recommendation{
			Action: ActionNone,
		}
	}

	return Recommendation{
		Action:               ActionRestartContainer,
		RequireHumanApproval: true,
	}
}
