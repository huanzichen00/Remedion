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
	"github.com/huanzichen00/remedion/internal/provider/docker"
	"github.com/huanzichen00/remedion/internal/provider/jev"
)

func main() {
	provider, err := docker.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	service := observe.NewService(provider, 10)
	observation, err := service.Observe(ctx, "remedion-demo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", observation)

	detector := incident.NewDetector(90, 90)

	inc, ok := detector.Detect(observation)
	if !ok {
		fmt.Println("no incident detected")
		return
	}

	fmt.Printf("incident detected: %+v\n", inc)

	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		log.Fatal("TYPESAFE_API_KEY is not set")
	}

	jevClient := jev.NewClient(apiKey)
	decisionService := decision.NewService(jevClient)

	result, err := decisionService.Decide(ctx, *inc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("decision: %+v\n", result)
}
