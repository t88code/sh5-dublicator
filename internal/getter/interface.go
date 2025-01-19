package getter

import (
	"context"
	"github.com/t88code/sh5-dublicator/domain"
)

type Getter interface {
	GetDictionary(ctx context.Context, procSync domain.ProcSync) (*domain.DictionarySync, error)
}
