package review

import (
	"fmt"
	"time"

	"github.com/huanzichen00/remedion/internal/incident"
	"github.com/huanzichen00/remedion/internal/policy"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(inc incident.Incident, rec policy.Recommendation) (Review, error) {
	if rec.Action == policy.ActionNone {
		return Review{}, fmt.Errorf("cannot create review for no-op recommendation")
	}

	return Review{
		ContainerID:    inc.ContainerID,
		ContainerName:  inc.ContainerName,
		Recommendation: rec,
		Status:         StatusPending,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *Service) Approve(r *Review) error {
	if r.Status != StatusPending {
		return fmt.Errorf("review is %q, expected pending", r.Status)
	}

	now := time.Now()

	r.Status = StatusApproved
	r.ReviewedAt = &now

	return nil
}

func (s *Service) Reject(r *Review) error {
	if r.Status != StatusPending {
		return fmt.Errorf("review is %q, expected pending", r.Status)
	}

	now := time.Now()

	r.Status = StatusRejected
	r.ReviewedAt = &now

	return nil
}
