package observe

import "time"

type Observation struct {
	Container  ContainerInfo
	Metrics    Metrics
	Logs       []string
	ObservedAt time.Time
}

type ContainerInfo struct {
	ID           string
	Name         string
	Image        string
	State        string
	Health       string
	RestartCount int
	OOMKilled    bool
}

type Metrics struct {
	CPUPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	MemoryPercent float64
	PIDs          uint64
}
