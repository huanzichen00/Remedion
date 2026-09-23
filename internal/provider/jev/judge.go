package jev

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huanzichen00/remedion/internal/decision"
	"github.com/huanzichen00/remedion/internal/incident"
)

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

// 编译期检查 *jev.Client 是否实现 decision.Judge
var _ decision.Judge = (*Client)(nil)

func (c *Client) Judge(ctx context.Context, inc incident.Incident) (decision.Decision, error) {
	var result decision.Decision

	req := buildRequest(inc)

	resp, err := c.systemOne(ctx, req)
	if err != nil {
		return result, err
	}

	rawCause, ok := resp.Answers["cause"]
	if !ok {
		return result, fmt.Errorf("Jev response missing cause answer")
	}

	var cause ChoiceAnswer

	if err := json.Unmarshal(rawCause, &cause); err != nil {
		return result, fmt.Errorf("decode Jev cause answer: %w", err)
	}

	switch cause.Choice {
	case "resource_exhaustion":
		result.Cause = decision.CauseResourceExhaustion

	case "application_failure":
		result.Cause = decision.CauseApplicationFailure

	case "unknown":
		result.Cause = decision.CauseUnknown

	default:
		return result, fmt.Errorf("unexpected Jev cause %q", cause.Choice)
	}

	result.CauseConfidence = cause.Confidence

	rawRestart, ok := resp.Answers["restart_helpful"]
	if !ok {
		return result, fmt.Errorf(
			"Jev response missing restart_helpful answer",
		)
	}

	var restart NoulAnswer
	if err := json.Unmarshal(rawRestart, &restart); err != nil {
		return result, fmt.Errorf(
			"decode Jev restart_helpful answer: %w",
			err,
		)
	}

	result.RestartHelpfulProbability = restart.Noul

	rawReview, ok := resp.Answers["needs_human_review"]
	if !ok {
		return result, fmt.Errorf(
			"Jev response missing needs_human_review answer",
		)
	}

	var review NoulAnswer
	if err := json.Unmarshal(rawReview, &review); err != nil {
		return result, fmt.Errorf(
			"decode Jev needs_human_review answer: %w",
			err,
		)
	}

	result.NeedsHumanReviewProbability = review.Noul

	return result, nil
}
