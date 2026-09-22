package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	if err := provider.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("docker daemon is reachable")

	metrics, err := provider.Stats(ctx, "remedion-deno")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", metrics)
}
