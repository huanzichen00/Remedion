package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/huanzichen00/remedion/internal/decision"
	"github.com/huanzichen00/remedion/internal/incident"
	"github.com/huanzichen00/remedion/internal/observe"
	"github.com/huanzichen00/remedion/internal/provider/jev"
)

func main() {
	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		log.Fatal("TYPESAFE_API_KEY is not set")
	}

	inc := incident.Incident{
		ContainerID:   "test-container",
		ContainerName: "jev-smoke-test",
		Signals: []incident.Signal{
			{Type: incident.SignalOOMKilled},
			{Type: incident.SignalHighMemory},
		},
		Observation: observe.Observation{
			Container: observe.ContainerInfo{
				ID:        "test-container",
				Name:      "jev-smoke-test",
				Image:     "example/app:latest",
				State:     "exited",
				Health:    "none",
				OOMKilled: true,
			},
			Metrics: observe.Metrics{
				CPUPercent:    21.4,
				MemoryUsage:   980 * 1024 * 1024,
				MemoryLimit:   1024 * 1024 * 1024,
				MemoryPercent: 95.7,
				PIDs:          12,
			},
			Logs: []string{
				"fatal error: runtime: out of memory",
				"process terminated unexpectedly",
			},
			ObservedAt: time.Now(),
		},
		DetectedAt: time.Now(),
	}

	client := jev.NewClient(apiKey)
	service := decision.NewService(client)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	result, err := service.Decide(ctx, inc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cause: %s\n", result.Cause)
	fmt.Printf("Cause confidence: %.2f\n", result.CauseConfidence)
	fmt.Printf(
		"Restart helpful probability: %.2f\n",
		result.RestartHelpfulProbability,
	)
	fmt.Printf(
		"Needs human review probability: %.2f\n",
		result.NeedsHumanReviewProbability,
	)
}
