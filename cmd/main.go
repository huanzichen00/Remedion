package main

import (
	"context"
	"remedion/internal/provider/docker"
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

	provider.Ping(ctx)
}
