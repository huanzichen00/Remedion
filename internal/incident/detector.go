package incident

import (
	"time"

	"github.com/huanzichen00/remedion/internal/observe"
)

type Detector struct {
	cpuThreshold    float64
	memoryThreshold float64
}

func NewDetector(cpuThreshold float64, memoryThreshold float64) *Detector {
	return &Detector{
		cpuThreshold:    cpuThreshold,
		memoryThreshold: memoryThreshold,
	}
}

func (d *Detector) Detect(obs observe.Observation) (*Incident, bool) {
	var signals []Signal

	if obs.Container.State != "running" {
		signals = append(signals, Signal{Type: SignalNotRunning})
	}
	if obs.Container.Health == "unhealthy" {
		signals = append(signals, Signal{Type: SignalUnhealthy})
	}

	if obs.Container.OOMKilled {
		signals = append(signals, Signal{Type: SignalOOMKilled})
	}

	if obs.Metrics.CPUPercent >= d.cpuThreshold {
		signals = append(signals, Signal{Type: SignalHighCPU})
	}

	if obs.Metrics.MemoryPercent >= d.memoryThreshold {
		signals = append(signals, Signal{Type: SignalHighMemory})
	}

	if len(signals) == 0 {
		return nil, false
	}

	return &Incident{
		ContainerID:   obs.Container.ID,
		ContainerName: obs.Container.Name,

		Signals:     signals,
		Observation: obs,
		DetectedAt:  time.Now(),
	}, true
}
