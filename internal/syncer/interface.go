package syncer

import (
	"context"
	"github.com/t88code/sh5-dublicator/domain"
)

type Syncer interface {
	SyncDictionary(ctx context.Context, procsSync []*domain.ProcSync) (err error)
}
