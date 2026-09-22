package observe

import (
	"context"
	"time"
)

type Service struct {
	provider ContainerProvider
	logLimit int
}

func NewService(provider ContainerProvider, logLimit int) *Service {
	return &Service{
		provider: provider,
		logLimit: logLimit,
	}
}

func (s *Service) Observe(ctx context.Context, target string) (Observation, error) {
	info, err := s.provider.Inspect(ctx, target)
	if err != nil {
		return Observation{}, err
	}

	stats, err := s.provider.Stats(ctx, target)
	if err != nil {
		return Observation{}, err
	}

	logs, err := s.provider.Logs(ctx, target, s.logLimit)
	if err != nil {
		return Observation{}, err
	}

	return Observation{
		Container:  info,
		Metrics:    stats,
		Logs:       logs,
		ObservedAt: time.Now(),
	}, nil
}
