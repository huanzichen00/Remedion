package observe

import "context"

type ContainerProvider interface {
	Inspect(ctx context.Context, target string) (ContainerInfo, error)
	Stats(ctx context.Context, target string) (Metrics, error)
	Logs(ctx context.Context, target string, limit int) ([]string, error)
}
