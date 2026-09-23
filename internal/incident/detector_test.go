package incident

import (
	"testing"

	"github.com/huanzichen00/remedion/internal/observe"
)

func TestDetectorDetect(t *testing.T) {
	detector := NewDetector(90, 90)

	tests := []struct {
		name        string
		observation observe.Observation
		wantOK      bool
		wantSignals []SignalType
	}{
		{
			name: "healthy container",
			observation: observe.Observation{
				Container: observe.ContainerInfo{
					State:     "running",
					Health:    "healthy",
					OOMKilled: false,
				},
				Metrics: observe.Metrics{
					CPUPercent:    20,
					MemoryPercent: 30,
				},
			},
			wantOK: false,
		},
		{
			name: "oom killed container",
			observation: observe.Observation{
				Container: observe.ContainerInfo{
					State:     "exited",
					Health:    "none",
					OOMKilled: true,
				},
				Metrics: observe.Metrics{
					CPUPercent:    10,
					MemoryPercent: 20,
				},
			},
			wantOK: true,
			wantSignals: []SignalType{
				SignalNotRunning,
				SignalOOMKilled,
			},
		},
		{
			name: "high cpu",
			observation: observe.Observation{
				Container: observe.ContainerInfo{
					State:  "running",
					Health: "healthy",
				},
				Metrics: observe.Metrics{
					CPUPercent:    95,
					MemoryPercent: 20,
				},
			},
			wantOK: true,
			wantSignals: []SignalType{
				SignalHighCPU,
			},
		},
		{
			name: "high memory",
			observation: observe.Observation{
				Container: observe.ContainerInfo{
					State:  "running",
					Health: "healthy",
				},
				Metrics: observe.Metrics{
					CPUPercent:    20,
					MemoryPercent: 96,
				},
			},
			wantOK: true,
			wantSignals: []SignalType{
				SignalHighMemory,
			},
		},
		{
			name: "unhealthy container",
			observation: observe.Observation{
				Container: observe.ContainerInfo{
					State:  "running",
					Health: "unhealthy",
				},
				Metrics: observe.Metrics{
					CPUPercent:    20,
					MemoryPercent: 30,
				},
			},
			wantOK: true,
			wantSignals: []SignalType{
				SignalUnhealthy,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncident, gotOK := detector.Detect(tt.observation)

			if gotOK != tt.wantOK {
				t.Fatalf(
					"Detect() ok = %v, want %v",
					gotOK,
					tt.wantOK,
				)
			}

			if !tt.wantOK {
				if gotIncident != nil {
					t.Fatalf("Detect() incident = %v, want nil", gotIncident)
				}

				return
			}

			if gotIncident == nil {
				t.Fatal("Detect() incident = nil, want non-nil")
			}

			if len(gotIncident.Signals) != len(tt.wantSignals) {
				t.Fatalf(
					"Detect() signals length = %d, want %d",
					len(gotIncident.Signals),
					len(tt.wantSignals),
				)
			}

			for i, wantSignal := range tt.wantSignals {
				if gotIncident.Signals[i].Type != wantSignal {
					t.Errorf(
						"signal[%d] = %q, want %q",
						i,
						gotIncident.Signals[i].Type,
						wantSignal,
					)
				}
			}
		})
	}
}
