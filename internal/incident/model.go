package incident

import (
	"time"

	"github.com/huanzichen00/remedion/internal/observe"
)

type SignalType string

const (
	SignalNotRunning SignalType = "not_running"
	SignalUnhealthy  SignalType = "unhealthy"
	SignalOOMKilled  SignalType = "oom_killed"
	SignalHighCPU    SignalType = "high_cpu"
	SignalHighMemory SignalType = "high_memory"
)

type Signal struct {
	Type SignalType
}

type Incident struct {
	ContainerID   string
	ContainerName string

	Signals     []Signal
	Observation observe.Observation

	DetectedAt time.Time
}
