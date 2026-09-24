package executor

import (
	"context"
	"fmt"

	"github.com/huanzichen00/remedion/internal/policy"
	"github.com/huanzichen00/remedion/internal/review"
)

type Service struct {
	provider ContainerProvider
}

func NewService(provider ContainerProvider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) Execute(ctx context.Context, r review.Review) error {
	if r.Status != review.StatusApproved {
		return fmt.Errorf("cannot execute review with status %q", r.Status)
	}

	switch r.Recommendation.Action {
	case policy.ActionRestartContainer:
		return s.provider.Restart(ctx, r.ContainerID)

	case policy.ActionNone:
		return fmt.Errorf("cannot execute no-op recommendation")

	default:
		return fmt.Errorf("unsupported action %q", r.Recommendation.Action)
	}
}
