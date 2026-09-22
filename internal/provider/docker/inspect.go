package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/huanzichen00/remedion/internal/observe"
	"github.com/moby/moby/client"
)

func (p *Provider) Inspect(ctx context.Context, target string) (observe.ContainerInfo, error) {
	inspectResult, err := p.client.ContainerInspect(ctx, target, client.ContainerInspectOptions{})
	if err != nil {
		return observe.ContainerInfo{}, fmt.Errorf("inspect container %q: %w", target, err)
	}

	container := inspectResult.Container

	state := ""
	health := "none"
	oomKilled := false

	if container.State != nil && container.State.Health != nil {
		state = string(container.State.Status)
		health = string(container.State.Health.Status)
		oomKilled = container.State.OOMKilled
	}

	info := observe.ContainerInfo{
		ID:           container.ID,
		Name:         strings.TrimPrefix(container.Name, "/"),
		Image:        container.Config.Image,
		State:        state,
		Health:       health,
		RestartCount: container.RestartCount,
		OOMKilled:    oomKilled,
	}

	return info, nil
}
