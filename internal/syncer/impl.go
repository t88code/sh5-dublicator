package syncer

import (
	"context"
	"fmt"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/comparer"
	"github.com/t88code/sh5-dublicator/internal/repository/sqlite"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"log/slog"
)

type syncer struct {
	logger            *slog.Logger
	cmr               comparer.Comparer
	sh5ClientSrc      sh5api.ClientInterface
	sh5ClientDst      sh5api.ClientInterface
	gGroupsRepository *sqlite.GGroupsRepository
}

func New(
	logger *slog.Logger,
	cmr comparer.Comparer,
	sh5ClientSrc sh5api.ClientInterface,
	sh5ClientDst sh5api.ClientInterface,
	gGroupsRepository *sqlite.GGroupsRepository) (Syncer, error) {
	return &syncer{
		logger:            logger,
		cmr:               cmr,
		sh5ClientSrc:      sh5ClientSrc,
		sh5ClientDst:      sh5ClientDst,
		gGroupsRepository: gGroupsRepository,
	}, nil
}

// UpdateAttrInSrc - обновить в sh5-src Attr-Rid_PARENT = attr
func (s *syncer) UpdateAttrInSrc(ctx context.Context, ggroupSrc domain.GGroups, procSync *domain.ProcSync, attr uint32) error {

	s.logger.Debug("UpdateRidInSrc")

	values := make([][]interface{}, 0)
	values = append(values, []interface{}{ggroupSrc.Rid})
	values = append(values, []interface{}{ggroupSrc.Guid})
	values = append(values, []interface{}{ggroupSrc.Rid_PARENT})
	values = append(values, []interface{}{ggroupSrc.Name})
	values = append(values, []interface{}{attr})

	if len(values) != len(procSync.OriginalForUpdateAttrInSrc) {
		return fmt.Errorf("не совпадает количество значений values(%d) и procSync.OriginalForUpdateAttrInSrc(%d), проверьте настройки OriginalForUpdateAttrInSrc и функцию UpdateRidInSrc",
			len(values), len(procSync.OriginalForUpdateAttrInSrc))
	}

	_, err := s.sh5ClientSrc.Sh5ExecWithInput(
		ctx,
		sh5api.UpdGGroup,
		procSync.Head,
		procSync.OriginalForUpdateAttrInSrc,
		values,
		[]sh5api.Sh5ExecStatus{sh5api.Sh5ExecStatusModify},
	)
	if err != nil {
		return err
	}

	return nil
}

// UpdateInDst - обновить в sh5-dst Rid_PARENT
func (s *syncer) UpdateInDst(ctx context.Context, ggroupDst domain.GGroups, procSync *domain.ProcSync) error {

	s.logger.Debug("UpdateInDst")

	values := make([][]interface{}, 0)
	values = append(values, []interface{}{ggroupDst.Rid})
	values = append(values, []interface{}{ggroupDst.Guid})
	values = append(values, []interface{}{ggroupDst.Rid_PARENT})
	values = append(values, []interface{}{ggroupDst.Name})
	values = append(values, []interface{}{ggroupDst.Attr_PARENT})

	if len(values) != len(procSync.OriginalForUpdateInDst) {
		return fmt.Errorf("не совпадает количество значений values(%d) и procSync.OriginalForUpdateInDst(%d), проверьте настройки OriginalForUpdateInDst и функцию UpdateInDst",
			len(values), len(procSync.OriginalForUpdateInDst))
	}

	_, err := s.sh5ClientDst.Sh5ExecWithInput(
		ctx,
		sh5api.UpdGGroup,
		procSync.Head,
		procSync.OriginalForUpdateInDst,
		values,
		[]sh5api.Sh5ExecStatus{sh5api.Sh5ExecStatusModify},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteInDst - удалить в sh5-dst
func (s *syncer) DeleteInDst(ctx context.Context, ggroupDst domain.GGroups, procSync *domain.ProcSync) error {
	s.logger.Debug("DeleteInDst")

	values := make([][]interface{}, 0)
	values = append(values, []interface{}{ggroupDst.Rid})
	values = append(values, []interface{}{ggroupDst.Guid})

	if len(values) != len(procSync.OriginalForDeleteInDst) {
		return fmt.Errorf("не совпадает количество значений values(%d) и procSync.OriginalForDeleteInDst(%d), проверьте настройки OriginalForDeleteInDst и функцию DeleteInDst",
			len(values), len(procSync.OriginalForDeleteInDst))
	}

	_, err := s.sh5ClientDst.Sh5ExecWithInput(
		ctx,
		sh5api.DelGGroup,
		procSync.Head,
		procSync.OriginalForDeleteInDst,
		values,
		[]sh5api.Sh5ExecStatus{sh5api.Sh5ExecStatusDelete},
	)
	if err != nil {
		return err
	}

	return nil
}

// InsertWithoutParentToDst - добавить запись справочника в базе sh5-dst с parent=1 и вернуть id
func (s *syncer) InsertWithoutParentToDst(ctx context.Context, ggroupSrc domain.GGroups, procSync *domain.ProcSync) (uint32, error) {

	s.logger.Debug("InsertWithoutParentToDst")

	values := make([][]interface{}, 0)
	values = append(values, []interface{}{ggroupSrc.Guid})
	values = append(values, []interface{}{ggroupSrc.Name})
	values = append(values, []interface{}{1})

	if len(values) != len(procSync.OriginalForInsertToDst) {
		return 0, fmt.Errorf("не совпадает количество значений values(%d) и procSync.OriginalForInsertToDst(%d), проверьте настройки OriginalForInsertToDst и функцию UpdateRidInSrc",
			len(values), len(procSync.OriginalForInsertToDst))
	}

	sh5ExecRep, err := s.sh5ClientDst.Sh5ExecWithInput(
		ctx,
		sh5api.InsGGroup,
		procSync.Head,
		procSync.OriginalForInsertToDst,
		values,
		[]sh5api.Sh5ExecStatus{sh5api.Sh5ExecStatusInsert},
	)
	if err != nil {
		return 0, err
	}

	id := uint32(sh5ExecRep.ShTable[0].Values[6][0].(float64))

	return id, nil
}

// SyncDictionary
//
// Обновляем ParentAttr
//   - CompareDictionary: считать все справочники и пометить их на insert, modify, delete в базе SQL
//   - HeadCodeGGROUP:
//     -- Проходимся по всем delete_rid_and_insert - обновляем ATTR-RID-PARENT в SRC
//     -- Проходимся по всем insert - обновляем ATTR-RID-PARENT в SRC
//     -- Проходимся по всем delete_rid_and_insert - обновляем ATTR-RID-PARENT в DST
//     -- Проходимся по всем insert - обновляем ATTR-RID-PARENT в DST
//     -- Проходимся по всем modify - обновляем все записи в DST
func (s *syncer) SyncDictionary(ctx context.Context, procsSync []*domain.ProcSync) (err error) {

	// считать все справочники и пометить их на insert, modify, delete в базе SQL
	compareDictionaries, err := s.cmr.CompareDictionary(ctx, procsSync)
	if err != nil {
		return err
	}

	s.logger.Debug("Приступаем к выполнению Action")
	for _, procSync := range procsSync {
		if procSync.IsActionDoIt {
			s.logger.Debug("Включен флаг procSync.IsActionDoIt, выполняем копирование",
				"procSync.Name", procSync.Name,
				"procSync.Head", procSync.Head,
			)
			compareDictionary, ok := compareDictionaries[procSync.Head]
			if !ok {
				return fmt.Errorf("procSync.Head(%v), procSync.Name(%v) не найден в compareDictionaries", procSync.Head, procSync.Name)
			}
			_ = compareDictionary

			switch procSync.Head {
			case sh5api.HeadCodeGGROUP:
				// Получить все записи из SRC-SQL-DB
				ggroupsSrc, err := s.gGroupsRepository.GetAll(ctx, domain.TABLE_ggroups_src)
				if err != nil {
					return err
				}

				// Проходимся по всем delete_rid_and_insert
				// - обновляем запись в SH5-SRC ATTR-RID-PARENT=0
				// - создаем запись в SH5-DST
				// - обновляем запись в SH5-SRC ATTR-RID-PARENT=ID
				for _, ggroupSrc := range ggroupsSrc {
					if ggroupSrc.Action == domain.ActionDeleteRidAndInsert {
						s.logger.Debug(string(ggroupSrc.Action),
							"procSync.Head", procSync.Head,
							"ggroupSrc.Rid", ggroupSrc.Rid,
							"ggroupSrc.Guid", ggroupSrc.Guid,
							"ggroupSrc.FIELD_42", ggroupSrc.FIELD_42,
							"ggroupSrc.Name", ggroupSrc.Name,
							"ggroupSrc.Attr", ggroupSrc.Attr,
							"ggroupSrc.Rid_PARENT", ggroupSrc.Rid_PARENT,
							"ggroupSrc.Guid_PARENT", ggroupSrc.Guid_PARENT,
							"ggroupSrc.FIELD_42_PARENT", ggroupSrc.FIELD_42_PARENT,
							"ggroupSrc.Name_PARENT", ggroupSrc.Name_PARENT,
							"ggroupSrc.Attr_PARENT", ggroupSrc.Attr_PARENT)

						err = s.UpdateAttrInSrc(ctx, ggroupSrc, procSync, 0)
						if err != nil {
							return err
						}
						id, err := s.InsertWithoutParentToDst(ctx, ggroupSrc, procSync)
						if err != nil {
							return err
						}
						err = s.UpdateAttrInSrc(ctx, ggroupSrc, procSync, id)
						if err != nil {
							return err
						}
					}
				}

				// Проходимся по всем insert
				// - создаем запись в SH5-DST
				// - обновляем запись в SH5-SRC ATTR-RID-PARENT=ID
				for _, ggroupSrc := range ggroupsSrc {
					if ggroupSrc.Action == domain.ActionInsert {
						s.logger.Debug(string(ggroupSrc.Action),
							"procSync.Head", procSync.Head,
							"ggroupSrc.Rid", ggroupSrc.Rid,
							"ggroupSrc.Guid", ggroupSrc.Guid,
							"ggroupSrc.FIELD_42", ggroupSrc.FIELD_42,
							"ggroupSrc.Name", ggroupSrc.Name,
							"ggroupSrc.Attr", ggroupSrc.Attr,
							"ggroupSrc.Rid_PARENT", ggroupSrc.Rid_PARENT,
							"ggroupSrc.Guid_PARENT", ggroupSrc.Guid_PARENT,
							"ggroupSrc.FIELD_42_PARENT", ggroupSrc.FIELD_42_PARENT,
							"ggroupSrc.Name_PARENT", ggroupSrc.Name_PARENT,
							"ggroupSrc.Attr_PARENT", ggroupSrc.Attr_PARENT)

						id, err := s.InsertWithoutParentToDst(ctx, ggroupSrc, procSync)
						if err != nil {
							return err
						}
						err = s.UpdateAttrInSrc(ctx, ggroupSrc, procSync, id)
						if err != nil {
							return err
						}
					}
				}

				// Проходимся по всем delete_rid_and_insert и insert и modify
				// - обновляем запись в SH5-DST=SH5-SRC with RID-PARENT-ATTR
				for _, ggroupSrc := range ggroupsSrc {
					if ggroupSrc.Action == domain.ActionDeleteRidAndInsert || ggroupSrc.Action == domain.ActionInsert || ggroupSrc.Action == domain.ActionModify {
						s.logger.Debug(string(ggroupSrc.Action),
							"procSync.Head", procSync.Head,
							"ggroupSrc.Rid", ggroupSrc.Rid,
							"ggroupSrc.Guid", ggroupSrc.Guid,
							"ggroupSrc.FIELD_42", ggroupSrc.FIELD_42,
							"ggroupSrc.Name", ggroupSrc.Name,
							"ggroupSrc.Attr", ggroupSrc.Attr,
							"ggroupSrc.Rid_PARENT", ggroupSrc.Rid_PARENT,
							"ggroupSrc.Guid_PARENT", ggroupSrc.Guid_PARENT,
							"ggroupSrc.FIELD_42_PARENT", ggroupSrc.FIELD_42_PARENT,
							"ggroupSrc.Name_PARENT", ggroupSrc.Name_PARENT,
							"ggroupSrc.Attr_PARENT", ggroupSrc.Attr_PARENT)

						if ggroupSrc.Attr == 0 {
							continue
						}

						ggroupSrcParent, err := s.gGroupsRepository.GetByRid(ctx, domain.TABLE_ggroups_src, ggroupSrc.Rid_PARENT)
						if err != nil {
							return err
						}

						switch {
						case len(ggroupSrcParent) == 0:
							return fmt.Errorf("не найден ggroupSrcParent для ggroupSrc(%v)", ggroupSrc)
						case len(ggroupSrcParent) == 1:
							ggroupDstNew := domain.GGroups{
								Rid:         ggroupSrc.Attr,
								Guid:        ggroupSrc.Guid,
								Rid_PARENT:  ggroupSrcParent[0].Attr,
								Name:        ggroupSrc.Name,
								Attr_PARENT: 0,
							}

							err = s.UpdateInDst(ctx, ggroupDstNew, procSync)
							if err != nil {
								return err
							}
						default:
							return fmt.Errorf("найдено несколько ggroupSrcParent для ggroupSrc(%v)", ggroupSrc)
						}

					}
				}

				// Получить все записи из DST-SQL-DB
				ggroupsDst, err := s.gGroupsRepository.GetAll(ctx, domain.TABLE_ggroups_dst)
				if err != nil {
					return err
				}

				// Проходимся по всем delete
				// - удаляем запись в SH5-DST
				for _, ggroupDst := range ggroupsDst {
					if ggroupDst.Action == domain.ActionDelete {
						s.logger.Debug(string(ggroupDst.Action),
							"procSync.Head", procSync.Head,
							"ggroupDst.Rid", ggroupDst.Rid,
							"ggroupDst.Guid", ggroupDst.Guid,
							"ggroupDst.FIELD_42", ggroupDst.FIELD_42,
							"ggroupDst.Name", ggroupDst.Name,
							"ggroupDst.Attr", ggroupDst.Attr,
							"ggroupDst.Rid_PARENT", ggroupDst.Rid_PARENT,
							"ggroupDst.Guid_PARENT", ggroupDst.Guid_PARENT,
							"ggroupDst.FIELD_42_PARENT", ggroupDst.FIELD_42_PARENT,
							"ggroupDst.Name_PARENT", ggroupDst.Name_PARENT,
							"ggroupDst.Attr_PARENT", ggroupDst.Attr_PARENT)

						err = s.DeleteInDst(ctx, ggroupDst, procSync)
						if err != nil {
							return err
						}
					}
				}
			}
		} else {
			s.logger.Debug("Выключен флаг procSync.IsActionDoIt, копирование пропускаем",
				"procSync.Name", procSync.Name,
				"procSync.Head", procSync.Head,
			)
		}
	}
	return nil
}
