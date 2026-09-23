package decision

import (
	"context"

	"github.com/huanzichen00/remedion/internal/incident"
)

type Service struct {
	judge Judge
}

func NewService(judge Judge) *Service {
	return &Service{
		judge: judge,
	}
}

func (s *Service) Decide(ctx context.Context, inc incident.Incident) (Decision, error) {
	return s.judge.Judge(ctx, inc)
}
