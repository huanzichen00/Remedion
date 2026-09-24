package docker

import (
	"context"
	"fmt"

	"github.com/huanzichen00/remedion/internal/executor"
	"github.com/moby/moby/client"
)

var _ executor.ContainerProvider = (*Provider)(nil)

func (p *Provider) Restart(ctx context.Context, target string) error {
	_, err := p.client.ContainerRestart(ctx, target, client.ContainerRestartOptions{})
	if err != nil {
		return fmt.Errorf("restart container %q: %w", target, err)
	}

	return nil
}
