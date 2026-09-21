package docker

import (
	"context"

	"github.com/huanzichen00/remedion/internal/observe"
	"github.com/moby/moby/client"
)

func (p *Provider) Inspect(ctx context.Context, target string) (observe.ContainerInfo, error) {
	inspectResult, err := p.client.ContainerInspect(ctx, target, client.ContainerInspectOptions{})
	if err != nil {
		return observe.ContainerInfo{}, err
	}

	container := inspectResult.Container

	info := observe.ContainerInfo{
		ID:           container.Name,
		Name:         container.ID,
		Image:        container.Image,
		State:        string(container.State.Status),
		Health:       string(container.State.Health.Status),
		RestartCount: container.RestartCount,
		OOMKilled:    container.State.OOMKilled,
	}

	return info, nil
}
