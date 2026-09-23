package jev

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/huanzichen00/remedion/internal/decision"
	"github.com/huanzichen00/remedion/internal/incident"
	"github.com/huanzichen00/remedion/internal/observe"
)

func TestClientJudge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}

			if r.URL.Path != "/v1/systemone" {
				t.Fatalf(
					"path = %s, want /v1/systemone",
					r.URL.Path,
				)
			}

			if got := r.Header.Get("Authorization"); got != "Bearer test-api-key" {
				t.Fatalf(
					"Authorization = %q, want %q",
					got,
					"Bearer test-api-key",
				)
			}

			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf(
					"Content-Type = %q, want application/json",
					got,
				)
			}

			var req SystemOneRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode request: %v", err)
			}

			if req.Model != "jev-latest" {
				t.Fatalf(
					"model = %q, want jev-latest",
					req.Model,
				)
			}

			if len(req.Questions) != 3 {
				t.Fatalf(
					"questions = %d, want 3",
					len(req.Questions),
				)
			}

			w.Header().Set("Content-Type", "application/json")

			_, _ = w.Write([]byte(`{
				"model": "jev-latest",
				"answers": {
					"cause": {
						"type": "choice",
						"choice": "resource_exhaustion",
						"confidence": 0.92,
						"probabilities": {
							"resource_exhaustion": 0.92,
							"application_failure": 0.05,
							"unknown": 0.03
						}
					},
					"restart_helpful": {
						"type": "noul",
						"noul": 0.81
					},
					"needs_human_review": {
						"type": "noul",
						"noul": 0.97
					}
				},
				"usage": {
					"input_tokens": 100,
					"output_tokens": 0
				}
			}`))
		},
	))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	inc := incident.Incident{
		ContainerID:   "container-123",
		ContainerName: "remedion-demo",
		Signals: []incident.Signal{
			{Type: incident.SignalOOMKilled},
			{Type: incident.SignalHighMemory},
		},
		Observation: observe.Observation{
			Container: observe.ContainerInfo{
				ID:        "container-123",
				Name:      "remedion-demo",
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

	got, err := client.Judge(context.Background(), inc)
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}

	if got.Cause != decision.CauseResourceExhaustion {
		t.Errorf(
			"Cause = %q, want %q",
			got.Cause,
			decision.CauseResourceExhaustion,
		)
	}

	if got.CauseConfidence != 0.92 {
		t.Errorf(
			"CauseConfidence = %v, want 0.92",
			got.CauseConfidence,
		)
	}

	if got.RestartHelpfulProbability != 0.81 {
		t.Errorf(
			"RestartHelpfulProbability = %v, want 0.81",
			got.RestartHelpfulProbability,
		)
	}

	if got.NeedsHumanReviewProbability != 0.97 {
		t.Errorf(
			"NeedsHumanReviewProbability = %v, want 0.97",
			got.NeedsHumanReviewProbability,
		)
	}
}
