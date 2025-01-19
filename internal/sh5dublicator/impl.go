package sh5dublicator

import (
	"context"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
)

type duplicator struct {
	Duplicator

	src       sh5api.ClientInterface
	dst       sh5api.ClientInterface
	procNames []sh5api.ProcName
}

func New(src sh5api.ClientInterface, dst sh5api.ClientInterface, procNames []sh5api.ProcName) (Duplicator, error) {
	return &duplicator{
		dst:       dst,
		src:       src,
		procNames: procNames,
	}, nil
}

func (d *duplicator) CopyDictionary(ctx context.Context) (err error) {
	return nil
}
