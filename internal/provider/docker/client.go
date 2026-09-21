package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type Provider struct {
	client *client.Client
}

func New() (*Provider, error) {
	client, err := client.New(client.FromEnv)
	if err != nil {
		return &Provider{}, err
	}
	return &Provider{
		client: client,
	}, nil
}

func (p *Provider) Ping(ctx context.Context) error {
	result, err := p.client.Ping(ctx, client.PingOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", result)
	return nil
}
