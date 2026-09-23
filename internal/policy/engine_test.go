package policy

import (
	"testing"

	"github.com/huanzichen00/remedion/internal/decision"
)

func TestEngineEvaluate(t *testing.T) {
	engine := NewEngine(0.7)

	tests := []struct {
		name         string
		probability  float64
		wantAction   Action
		wantApproval bool
	}{
		{
			name:         "below threshold",
			probability:  0.29,
			wantAction:   ActionNone,
			wantApproval: false,
		},
		{
			name:         "equal threshold",
			probability:  0.70,
			wantAction:   ActionRestartContainer,
			wantApproval: true,
		},
		{
			name:         "above threshold",
			probability:  0.90,
			wantAction:   ActionRestartContainer,
			wantApproval: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := decision.Decision{
				RestartHelpfulProbability: tt.probability,
			}

			got := engine.Evaluate(dec)

			if got.Action != tt.wantAction {
				t.Errorf(
					"Action = %q, want %q",
					got.Action,
					tt.wantAction,
				)
			}

			if got.RequiresHumanApproval != tt.wantApproval {
				t.Errorf(
					"RequireHumanApproval = %v, want %v",
					got.RequiresHumanApproval,
					tt.wantApproval,
				)
			}
		})
	}
}
