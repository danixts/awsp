package port

import (
	"context"

	"github.com/danixts/awsp/internal/domain"
)

type LogService interface {
	ListLogGroups(prefix string) ([]domain.LogGroup, error)
	TailLogs(ctx context.Context, logGroup, since, region string) (<-chan string, error)
	FetchLogs(logGroup, since, region string) (string, error)
}
