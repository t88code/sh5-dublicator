package comparer

import (
	"context"
	"fmt"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/getter"
	"github.com/t88code/sh5-dublicator/internal/repository/sqlite"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"log/slog"
	"slices"
)

type NormalizedDictionary struct {
	Src *domain.DictionarySync
	Dst *domain.DictionarySync
}

type comparer struct {
	logger            *slog.Logger
	gtSrc             getter.Getter
	gtDst             getter.Getter
	gGroupsRepository *sqlite.GGroupsRepository
}

func New(logger *slog.Logger, gtSrc getter.Getter, gtDst getter.Getter, gGroupsRepository *sqlite.GGroupsRepository) (Comparer, error) {
	return &comparer{
		logger:            logger,
		gtSrc:             gtSrc,
		gtDst:             gtDst,
		gGroupsRepository: gGroupsRepository,
	}, nil
}

// CompareDictionary
//
//	 Реализовано для словарей:
//	 sh5api.GGroupsTree:
//		-GetNormalizeDictionary - получить справочники(Dictionary) Src/Dst и выполнить их нормализацию
//		-SaveGGroupDictionary - сохранить GGroup(торговую группа) в базе SQL
//		-CompareGGroupsTree - помечаем на insert, modify, delete
func (c *comparer) CompareDictionary(ctx context.Context, procsSync []*domain.ProcSync) (map[sh5api.Head]*NormalizedDictionary, error) {

	normalizedDictionaries := make(map[sh5api.Head]*NormalizedDictionary)

	for _, procSync := range procsSync {
		c.logger.Debug(fmt.Sprintf("procSync: %v", procSync.Name))
		normalizedDictionary, err := c.GetNormalizeDictionary(ctx, procSync)
		if err != nil {
			return nil, err
		}
		normalizedDictionaries[procSync.Head] = normalizedDictionary

		switch procSync.Name {
		case sh5api.GGroupsTree:
			err = c.SaveGGroupDictionary(ctx, normalizedDictionary.Src, domain.TABLE_ggroups_src)
			if err != nil {
				return nil, err
			}

			err = c.SaveGGroupDictionary(ctx, normalizedDictionary.Dst, domain.TABLE_ggroups_dst)
			if err != nil {
				return nil, err
			}

			err = c.CompareGGroupsTree(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return normalizedDictionaries, nil
}

// GetNormalizeDictionary - получить справочники(Dictionary) Src/Dst и выполнить их нормализацию:
//
//	-заполнить DictionarySync.ProcSync.TableIndex и DictionarySync.ProcSync.OriginalIndex
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

	*procSync = *dictionarySrc.ProcSync

	normalizedDictionary := NormalizedDictionary{
		Src: dictionarySrc,
		Dst: dictionaryDst,
	}

	return &normalizedDictionary, nil
}

// NormalizeIndexDictionary - выполнить их нормализацию
//
//	найти DictionarySync.ProcSync.TableIndex - индекс таблицы, в котором хранятся значения Values,
//	необходимы для сравнения в будущем
//
//	найти DictionarySync.OriginalsNormalized - соответствия filed.path и IndexInValue(indexOriginal) в Values
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

	for _, path := range dictionary.ProcSync.OriginalForSaveToSql {
		indexOriginal := slices.Index(dictionary.Sh5ExecRep.ShTable[dictionary.TableIndex].Original, path)
		if indexOriginal != -1 {
			dictionary.OriginalsNormalized = append(dictionary.OriginalsNormalized, domain.OriginalNormalized{
				Path:         path,
				IndexInValue: indexOriginal,
			})
		} else {
			//return fmt.Errorf("dictionary.ProcSync.OriginalForSaveToSql.path(%v) не найден в dictionary.Sh5ExecRep.ShTable[%d].OriginalForSaveToSql", path, dictionary.TableIndex)
			c.logger.Error(fmt.Sprintf("dictionary.ProcSync.OriginalForSaveToSql.path(%v) не найден в dictionary.Sh5ExecRep.ShTable[%d].OriginalForSaveToSql", path, dictionary.TableIndex))
		}
	}
	if len(dictionary.OriginalsNormalized) != len(dictionary.ProcSync.OriginalForSaveToSql) {
		//return fmt.Errorf("не совпадает количество найденных индексов procSync.OriginalIndex(%d) и количество значений procSync.OriginalForSaveToSql(%d)", len(dictionary.OriginalsNormalized), len(dictionary.ProcSync.OriginalForSaveToSql))
		c.logger.Error(fmt.Sprintf("не совпадает количество найденных индексов procSync.OriginalIndex(%d) и количество значений procSync.OriginalForSaveToSql(%d)", len(dictionary.OriginalsNormalized), len(dictionary.ProcSync.OriginalForSaveToSql)))
	}
	return nil
}
