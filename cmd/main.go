package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/huanzichen00/remedion/internal/observe"
	"github.com/huanzichen00/remedion/internal/provider/docker"
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
}
