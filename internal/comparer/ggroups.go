package comparer

import (
	"context"
	"fmt"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/helper"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
)

// SaveGGroupDictionary
//
// Добавляем запись справочника в базу SQL
// Изначально GGroup таблица в базе SQL должна быть пустая пустая, так как она создается при запуске сервиса
func (c *comparer) SaveGGroupDictionary(ctx context.Context, dictionarySync *domain.DictionarySync, tableName string) error {
	c.logger.Debug("заполняем таблицу", "table", tableName)
	table := dictionarySync.Sh5ExecRep.ShTable[dictionarySync.TableIndex]

	for indexValue, _ := range table.Values[0] {
		ggroup := domain.GGroups{}
		for _, field := range dictionarySync.OriginalsNormalized {
			switch field.Path {
			case sh5api.FIELD_1_RID:
				ggroup.Rid = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_4_GUID:
				ggroup.Guid = helper.GetStringFromInterfaceStringNullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_42:
				ggroup.FIELD_42 = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_3_NAME:
				ggroup.Name = helper.GetStringFromInterfaceStringNullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_6_RidDst:
				ggroup.Attr = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_209_1_RID_PARENT:
				ggroup.Rid_PARENT = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_209_4_GUID_PARENT:
				ggroup.Guid_PARENT = helper.GetStringFromInterfaceStringNullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_209_42:
				ggroup.FIELD_42_PARENT = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_209_3_NAME_PARENT:
				ggroup.Name_PARENT = helper.GetStringFromInterfaceStringNullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_209_6_RidDst_PARENT:
				ggroup.Attr_PARENT = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_239:
				ggroup.UserGroup = helper.GetUint64FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_106_1:
				ggroup.Rid_106 = helper.GetUint32FromInterfaceFloat64Nullable(table.Values[field.IndexInValue][indexValue])
			case sh5api.FIELD_106_3:
				ggroup.Name_106 = helper.GetStringFromInterfaceStringNullable(table.Values[field.IndexInValue][indexValue])
			}
		}
		c.logger.Debug(fmt.Sprintf("%+v", ggroup))
		err := c.gGroupsRepository.InsertAll(ctx, &ggroup, tableName)
		if err != nil {
			return err
		}

	}
	return nil
}

// CompareGGroupsTree
//
//	Перебираем domain.TABLE_ggroups_src - помечаем на insert, modify
//	Перебираем domain.TABLE_ggroups_dst - помечаем на delete
func (c *comparer) CompareGGroupsTree(ctx context.Context) error {
	c.logger.Debug("запускаем процесс сравнения таблиц",
		"src", domain.TABLE_ggroups_src, "dst", domain.TABLE_ggroups_dst)

	// перебираем все значения в базе SRC
	ggroupsSrc, err := c.gGroupsRepository.GetAll(ctx, domain.TABLE_ggroups_src)
	if err != nil {
		return err
	}
	for _, ggroupSrc := range ggroupsSrc {
		c.logger.Debug(fmt.Sprint(ggroupSrc))
		if ggroupSrc.Rid == 1 {
			continue
		}

		if ggroupSrc.Attr == 0 {
			c.logger.Warn("требуется добавить")
			ggroupSrc.Action = domain.ActionInsert
		} else {
			groupsDst, err := c.gGroupsRepository.GetByRid(
				ctx,
				domain.TABLE_ggroups_dst,
				ggroupSrc.Attr)
			if err != nil {
				return err
			}
			switch {
			case len(groupsDst) == 0:
				c.logger.Warn("требуется добавить в dst")
				ggroupSrc.Action = domain.ActionDeleteRidAndInsert
			case len(groupsDst) == 1:
				c.logger.Debug(fmt.Sprint(groupsDst))
				if ggroupSrc.Name != groupsDst[0].Name || ggroupSrc.FIELD_42 != groupsDst[0].FIELD_42 || ggroupSrc.Name_PARENT != groupsDst[0].Name_PARENT {
					c.logger.Warn("требуется обновить в dst")
					ggroupSrc.Action = domain.ActionModify
				} else {
					c.logger.Warn("обновление не требуется")
					continue
				}
			default:
				return fmt.Errorf("в dst найдены дубли rid(%s) в количестве(%d)", ggroupSrc.Attr, len(groupsDst))
			}
		}
		err = c.gGroupsRepository.UpdateAction(ctx, domain.TABLE_ggroups_src, &ggroupSrc)
		if err != nil {
			return err
		}
	}

	// перебираем все значения в базе DST
	ggroupsDst, err := c.gGroupsRepository.GetAll(ctx, domain.TABLE_ggroups_dst)
	if err != nil {
		return err
	}
	for _, ggroupDst := range ggroupsDst {
		c.logger.Debug(fmt.Sprint(ggroupDst))
		if ggroupDst.Rid == 1 {
			continue
		}

		groupsSrc, err := c.gGroupsRepository.GetByAttr(
			ctx,
			domain.TABLE_ggroups_src,
			ggroupDst.Rid)
		if err != nil {
			return err
		}

		if len(groupsSrc) == 0 {
			c.logger.Warn("требуется удалить из dst")
			ggroupDst.Action = domain.ActionDelete
		}

		err = c.gGroupsRepository.UpdateAction(ctx, domain.TABLE_ggroups_dst, &ggroupDst)
		if err != nil {
			return err
		}
	}

	return nil
}
