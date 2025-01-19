package comparer

import (
	"context"
	"fmt"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/getter"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"log/slog"
	"reflect"
	"slices"
)

type CompareResult struct {
}

type NormalizedDictionary struct {
	Src *domain.DictionarySync
	Dst *domain.DictionarySync
}

type comparer struct {
	logger *slog.Logger
	gtSrc  getter.Getter
	gtDst  getter.Getter
}

func New(logger *slog.Logger, gtSrc getter.Getter, gtDst getter.Getter) (Comparer, error) {
	return &comparer{
		logger: logger,
		gtSrc:  gtSrc,
		gtDst:  gtDst,
	}, nil
}

// NormalizeIndexDictionary - выполнить их нормализацию -
// заполнить DictionarySync.ProcSync.TableIndex и DictionarySync.ProcSync.OriginalIndex
func (c *comparer) NormalizeIndexDictionary(dictionary *domain.DictionarySync) error {
	// выполняем поиск индекса таблицы, которую используем для копирования в новую БД
	var findTableIndex = false
	for index, table := range dictionary.Sh5ExecRep.ShTable {
		if table.Head == dictionary.ProcSync.Head {
			dictionary.TableIndex = index // TODO проверить что индексы все сохраняются при множественном выборе, когда несколько dictionary и несколько таблиц внутри dictionary
			findTableIndex = true
		}
	}
	if !findTableIndex {
		return fmt.Errorf("не найдена dictionary.Sh5ExecRep.ShTable.Head(%s) в ответе после выполнения процедуры ProcName(%s)", dictionary.ProcSync.Head, dictionary.ProcSync.Name)
	}

	for _, path := range dictionary.ProcSync.Original {
		indexOriginal := slices.Index(dictionary.Sh5ExecRep.ShTable[dictionary.TableIndex].Original, path)
		if indexOriginal != -1 {
			dictionary.OriginalsNormalized = append(dictionary.OriginalsNormalized, domain.OriginalNormalized{
				Path:         path,
				IndexInValue: indexOriginal,
			})
		} else {
			return fmt.Errorf("dictionary.ProcSync.Original.path(%v) не найден в dictionary.Sh5ExecRep.ShTable[%d].Original", path, dictionary.TableIndex)
		}
	}
	if len(dictionary.OriginalsNormalized) != len(dictionary.ProcSync.Original) {
		return fmt.Errorf("не совпадает количество найденных индексов procSync.OriginalIndex(%d) и количество значений procSync.Original(%d)", len(dictionary.OriginalsNormalized), len(dictionary.ProcSync.Original))
	}
	return nil
}

// NormalizeValuesDictionary - заполняет DictionarySync.ValuesNormalized
func (c *comparer) NormalizeValuesDictionary(dictionary *domain.DictionarySync) error {
	err := dictionary.CheckTableIndexAndValues()
	if err != nil {
		return err
	}

	valuesOriginal := dictionary.Sh5ExecRep.ShTable[dictionary.TableIndex].Values
	for _, valueOriginal := range valuesOriginal {
		guid, err := dictionary.GetValueByFieldPath(valueOriginal, sh5api.FIELD_4_GUID)
		if err != nil {
			return err
		}

		rid, err := dictionary.GetValueByFieldPath(valueOriginal, sh5api.FIELD_1_RID)
		if err != nil {
			return err
		}

		normalizedValue, err := dictionary.GetNormalizedValue(valueOriginal)
		if err != nil {
			return err
		}

		dictionary.ValuesNormalized[guid] = domain.ValuesNormalized{
			Value:         normalizedValue,
			Sh5ExecStatus: "",
			Guid:          guid,
			RID:           rid,
		}
	}

	return nil
}

// GetNormalizeDictionary - получить справочники(Dictionary) Src/Dst и выполнить их нормализацию:
//
//	-заполнить DictionarySync.ProcSync.TableIndex и DictionarySync.ProcSync.OriginalIndex
//	-заполнить ValuesNormalized и присвоить им тип операции: insert, modify, delete
func (c *comparer) GetNormalizeDictionary(ctx context.Context, procSync *domain.ProcSync) (*NormalizedDictionary, error) {
	// получить справочники
	dictionarySrc, err := c.gtSrc.GetDictionary(ctx, *procSync)
	if err != nil {
		return nil, err
	}

	dictionaryDst, err := c.gtDst.GetDictionary(ctx, *procSync)
	if err != nil {
		return nil, err
	}

	// получить индексы DictionarySync.ProcSync.TableIndex и DictionarySync.ProcSync.OriginalIndex
	err = c.NormalizeIndexDictionary(dictionarySrc)
	if err != nil {
		return nil, err
	}

	err = c.NormalizeIndexDictionary(dictionaryDst)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(dictionarySrc.OriginalsNormalized, dictionaryDst.OriginalsNormalized) &&
		dictionarySrc.TableIndex == dictionaryDst.TableIndex {
		*procSync = *dictionarySrc.ProcSync
	} else {
		return nil, fmt.Errorf("не совадают значения dictionarySrc.ProcSync.OriginalIndex(%v)!=dictionaryDst.ProcSync.OriginalIndex(%v) и dictionarySrc.ProcSync.TableIndex(%d)!=dictionaryDst.ProcSync.TableIndex(%d)",
			dictionarySrc.OriginalsNormalized, dictionaryDst.OriginalsNormalized, dictionarySrc.TableIndex, dictionaryDst.TableIndex)
	}

	// получить values и их типы операций
	err = c.NormalizeValuesDictionary(dictionarySrc)
	if err != nil {
		return nil, err
	}

	err = c.NormalizeValuesDictionary(dictionaryDst)
	if err != nil {
		return nil, err
	}

	normalizedDictionary := NormalizedDictionary{
		Src: dictionarySrc,
		Dst: dictionaryDst,
	}

	return &normalizedDictionary, nil
}

func (c *comparer) CompareDictionary(ctx context.Context, procsSync []*domain.ProcSync) (map[sh5api.ProcName]CompareResult, error) {

	compareResults := make(map[sh5api.ProcName]CompareResult)

	for _, procSync := range procsSync {

		normalizedDictionary, err := c.GetNormalizeDictionary(ctx, procSync)
		if err != nil {
			return nil, err
		}

		// TODO сравнить dictionarySrc и dictionaryDst и проставить статусы

		//////////////////////////

		fmt.Printf("%+v\n", normalizedDictionary.Src)
		fmt.Printf("%+v\n", normalizedDictionary.Dst)

		fmt.Printf("%+v\n", normalizedDictionary.Src.Sh5ExecRep)
		fmt.Printf("%+v\n", normalizedDictionary.Dst.Sh5ExecRep)

		fmt.Printf("%+v\n", normalizedDictionary.Src.ValuesNormalized)
		fmt.Printf("%+v\n", normalizedDictionary.Dst.ValuesNormalized)

		fmt.Println(normalizedDictionary.Src.ProcSync)

		break
	}

	//for _, procName := range procNames {
	//
	//	reqRepSrc, ok := reqRepMapSrc[procName]
	//	if !ok {
	//		return nil, fmt.Errorf("справочник (%s) не найден в базе источнике", procName)
	//	}
	//
	//	reqRepDst, ok := reqRepMapDst[procName]
	//	if !ok {
	//		return nil, fmt.Errorf("справочник (%s) не найден в базе получателе", procName)
	//	}
	//
	//	dataJsonBytesSrc, err := json.MarshalIndent(reqRepSrc.Rep, "", "    ")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	dataJsonBytesDst, err := json.MarshalIndent(reqRepDst.Rep, "", "    ")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	if !reflect.DeepEqual(dataJsonBytesSrc, dataJsonBytesDst) {
	//		compareResults[procName] = CompareResult{
	//			ReqRepSrc: reqRepSrc,
	//			ReqRepDst: reqRepDst,
	//		}
	//	}
	//}

	return compareResults, nil
}
