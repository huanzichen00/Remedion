package review

import (
	"time"

	"github.com/huanzichen00/remedion/internal/policy"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Review struct {
	ContainerID   string
	ContainerName string

	Recommendation policy.Recommendation

	Status Status

	CreatedAt  time.Time
	ReviewedAt *time.Time
}
