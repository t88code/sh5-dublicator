package sh5dublicator

import (
	"context"
)

type Duplicator interface {
	CopyDictionary(ctx context.Context) (err error)
}
