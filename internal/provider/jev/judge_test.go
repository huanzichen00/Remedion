package jev

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/huanzichen00/remedion/internal/incident"
	"github.com/huanzichen00/remedion/internal/observe"
)

func TestBuildRequest(t *testing.T) {
	inc := incident.Incident{
		ContainerName: "remedion-demo",
		Signals: []incident.Signal{
			{Type: incident.SignalOOMKilled},
			{Type: incident.SignalHighMemory},
		},
		Observation: observe.Observation{
			Container: observe.ContainerInfo{
				State:     "exited",
				Health:    "none",
				OOMKilled: true,
			},
			Metrics: observe.Metrics{
				CPUPercent:    12.3,
				MemoryPercent: 96.4,
			},
			Logs: []string{
				"fatal: out of memory",
			},
		},
	}

	req := buildRequest(inc)

	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}
