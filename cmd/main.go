package main

import (
	"context"
	"fmt"
	"github.com/huanzichen00/remedion/internal/provider/docker"
	"log"
	"time"
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
}
