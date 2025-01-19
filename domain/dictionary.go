package domain

import (
	"fmt"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"slices"
)

type DictionarySync struct {
	ProcSync            *ProcSync
	Sh5ExecRep          *sh5api.Sh5ExecRep
	TableIndex          int // Индекс таблицы, по которой происходит синхронизация. Если 0, то значит ProcSync.Head не найден и вероятно некорректное поведение.
	OriginalsNormalized []OriginalNormalized
	ValuesNormalized    map[string]ValuesNormalized // GUID -> value
}

type OriginalNormalized struct {
	Path         sh5api.FieldPath
	IndexInValue int
}

type ValuesNormalized struct {
	Value         []interface{}
	Sh5ExecStatus sh5api.Sh5ExecStatus
	Guid          string // 4
	RID           string // 1
}

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

func (d *DictionarySync) GetValueByFieldPath(valueOriginal []interface{}, fieldPath sh5api.FieldPath) (string, error) {
	indexFieldPath := slices.Index(d.ProcSync.Original, fieldPath)
	if indexFieldPath == -1 {
		return "", fmt.Errorf("не найден FieldPath(%v) в Dictionary.ProcSync.Original(%v)", fieldPath, d.ProcSync.Original)
	}

	if valueFieldPath, ok := valueOriginal[indexFieldPath].(string); ok {
		if valueFieldPath == "" {
			return "", fmt.Errorf("valueFieldPath(%v) не найден для value(%v), d.TableIndex(%d), indexFieldPath(%d)", valueFieldPath, valueOriginal, d.TableIndex, indexFieldPath)
		}
		return valueFieldPath, nil
	} else {
		return "", fmt.Errorf("приведение типа к (string) невозможно для value(%v), valueOriginal[%d](%v), d.TableIndex(%d), indexFieldPath(%d)", valueOriginal, indexFieldPath, valueOriginal[indexFieldPath], d.TableIndex, indexFieldPath)
	}
}

func (d *DictionarySync) GetNormalizedValue(valueOriginal []interface{}) ([]interface{}, error) {
	value := make([]interface{}, 0)
	for _, originalNormalized := range d.OriginalsNormalized {
		value = append(value, valueOriginal[originalNormalized.IndexInValue])
	}
	return value, nil
}
