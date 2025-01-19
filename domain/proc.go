package domain

import "github.com/t88code/sh5-dublicator/pkg/sh5api"

type ProcSync struct {
	Name     sh5api.ProcName    // Название процедуры, в результате которой можно: получить таблицу, выполнить изменение/создание/удаление.
	Head     sh5api.Head        // Идентификатор Head по которому будет выполняться поиск таблицы(table) в ответе на запрос /api/sh5exec; замечено, что в ответе может быть несколько таблиц, например 209#1 и 209.
	Original []sh5api.FieldPath // Поля, которые необходимо проверять и копировать в новую БД
}
