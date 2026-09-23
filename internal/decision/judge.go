package decision

import (
	"context"

	"github.com/huanzichen00/remedion/internal/incident"
)

type Judge interface {
	Judge(ctx context.Context, incident incident.Incident) (Decision, error)
}
