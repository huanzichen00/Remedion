package jev

import "github.com/huanzichen00/remedion/internal/incident"

type incidentState struct {
	ContainerName string   `json:"container_name"`
	State         string   `json:"state"`
	Health        string   `json:"health"`
	OOMKilled     bool     `json:"oom_killed"`
	CPUPercent    float64  `json:"cpu_percent"`
	MemoryPercent float64  `json:"memory_percent"`
	Signals       []string `json:"signals"`
	Logs          []string `json:"logs"`
}

func buildRequest(inc incident.Incident) SystemOneRequest {
	var signals []string

	for _, signal := range inc.Signals {
		signals = append(signals, string(signal.Type))
	}

	state := incidentState{
		ContainerName: inc.ContainerName,
		State:         inc.Observation.Container.State,
		Health:        inc.Observation.Container.Health,
		OOMKilled:     inc.Observation.Container.OOMKilled,
		CPUPercent:    inc.Observation.Metrics.CPUPercent,
		MemoryPercent: inc.Observation.Metrics.MemoryPercent,
		Signals:       signals,
		Logs:          inc.Observation.Logs,
	}

	questions := map[string]any{
		"cause": ChoiceQuestion{
			Type:         "choice",
			Instructions: "Identify the most likely primary cause of this container incident.",
			Criteria: map[string]string{
				"resource_exhaustion": "CPU or memory resource exhaustion.",
				"application_failure": "Failure originating from the application itself.",
				"unknown":             "There is not enough evidence to identify the primary cause.",
			},
		},

		"restart_helpful": NoulQuestion{
			Type:         "noul",
			Instructions: "Would restarting this container likely help restore normal service?",
		},

		"needs_human_review": NoulQuestion{
			Type:         "noul",
			Instructions: "Should a human review this incident before any remediation action is executed?",
		},
	}

	return SystemOneRequest{
		Model:     "jev-latest",
		State:     state,
		Questions: questions,
	}
}
