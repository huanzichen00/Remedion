package executor

import "context"

type ContainerProvider interface {
	Restart(ctx context.Context, target string) error
}
