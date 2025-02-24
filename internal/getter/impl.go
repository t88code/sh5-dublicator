package getter

import (
	"context"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
)

type getter struct {
	sh5Client sh5api.ClientInterface
}

func New(sh5Client sh5api.ClientInterface) (Getter, error) {
	return &getter{
		sh5Client: sh5Client,
	}, nil
}

// GetDictionary
//
// - получить справочник по названию процедуры
// - необходимо учитывать, что ProcName должен возвращать справочник, а не выполнять создание/модификацию/удаление
func (g *getter) GetDictionary(ctx context.Context, procSync domain.ProcSync) (*domain.DictionarySync, error) {
	sh5ExecRep, err := g.sh5Client.Sh5ExecOnlyProcName(ctx, procSync.Name)
	if err != nil {
		return nil, err
	}

	return &domain.DictionarySync{
		ProcSync:            &procSync,
		Sh5ExecRep:          sh5ExecRep,
		TableIndex:          -1,
		OriginalsNormalized: nil,
	}, nil
}
