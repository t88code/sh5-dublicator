package domain

import "github.com/t88code/sh5-dublicator/pkg/sh5api"

// ProcSync - основная настройка для справочника, который будет синхронизироваться
type ProcSync struct {
	Name                       sh5api.ProcName    // Название процедуры, в результате которой можно: получить таблицу, выполнить изменение/создание/удаление.
	Head                       sh5api.Head        // Идентификатор Head по которому будет выполняться поиск таблицы(table) в ответе на запрос /api/sh5exec; замечено, что в ответе может быть несколько таблиц, например 209#1 и 209.
	OriginalForSaveToSql       []sh5api.FieldPath // Поля, которые необходимо сохранить в SQL
	OriginalForInsertToDst     []sh5api.FieldPath // Поля, которые будут скопированы в новую БД
	OriginalForUpdateAttrInSrc []sh5api.FieldPath // Поля, которые будут обновлены в исходной БД, используется для обновления Attr=RID_DST
	OriginalForUpdateInDst     []sh5api.FieldPath // Поля, которые будут обновлены в новой БД
	OriginalForDeleteInDst     []sh5api.FieldPath // Поля, которые будут использоваться для удаления в новой БД
	IsActionDoIt               bool               // Флаг включает выполнение действий в SH5. Если выключен, то будет только заполняться БД SQL.
}
