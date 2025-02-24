package domain

import (
	"fmt"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
)

type DictionarySync struct {
	ProcSync            *ProcSync            // Все настройки для синхронизации справочника
	Sh5ExecRep          *sh5api.Sh5ExecRep   // Ответ, который вернул API на ProcSync.Name
	TableIndex          int                  // Индекс таблицы, в которой хранится справочник ProcSync.Head. Если 0, то значит ProcSync.Head не найден и вероятно некорректное поведение.
	OriginalsNormalized []OriginalNormalized // Originals оригинальный c индексами
}

type OriginalNormalized struct {
	Path         sh5api.FieldPath
	IndexInValue int
}

// CheckTableIndexAndValues - (не используется)проверка результатов нормализации
func (d *DictionarySync) CheckTableIndexAndValues() error {
	if d.TableIndex == -1 {
		return fmt.Errorf("Dictionary.TableIndex не найден = -1")
	}

	if len(d.Sh5ExecRep.ShTable) <= d.TableIndex {
		return fmt.Errorf("Dictionary.TableIndex(%d) не может быть больше или равен len(d.Sh5ExecRep.ShTable)(%d)", d.TableIndex, len(d.Sh5ExecRep.ShTable))
	}

	if len(d.Sh5ExecRep.ShTable[d.TableIndex].Values) != len(d.OriginalsNormalized) {
		return fmt.Errorf("длина len(d.Sh5ExecRep.ShTable[d.TableIndex].Values)(%d) != len(d.OriginalsNormalized)(%d)", len(d.Sh5ExecRep.ShTable[d.TableIndex].Values), len(d.OriginalsNormalized))

	}
	return nil
}
