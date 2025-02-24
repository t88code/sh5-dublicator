package comparer

import (
	"context"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
)

type Comparer interface {
	CompareDictionary(ctx context.Context, procsSync []*domain.ProcSync) (map[sh5api.Head]*NormalizedDictionary, error)
}
