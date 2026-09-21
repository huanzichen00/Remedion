package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type Provider struct {
	client *client.Client
}

func New() (*Provider, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Provider{
		client: cli,
	}, nil
}

func (p *Provider) Ping(ctx context.Context) error {
	_, err := p.client.Ping(ctx, client.PingOptions{})
	return err
}
