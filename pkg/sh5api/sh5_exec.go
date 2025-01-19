package sh5api

// Sh5ExecReq - структура запроса /api/sh5exec
type Sh5ExecReq struct {
	UserName string         `json:"UserName"`
	Password string         `json:"Password"`
	ProcName ProcName       `json:"procName"`
	Input    []Sh5ExecInput `json:"Input"`
}

// ProcName - название процедуры в запросе /api/sh5exec
type ProcName string

var (
	GGroups     ProcName = "GGroups"     // Получить справочник Товарная группа
	GGroupsTree ProcName = "GGroupsTree" // Получить дерево Товарная группа
)

// Sh5ExecInput - структура таблицы в запросе /api/sh5exec
type Sh5ExecInput struct {
	Head     Head            `json:"head"`
	Original []FieldPath     `json:"original"`
	Fields   []FieldAlt      `json:"fields"`
	Values   [][]interface{} `json:"values"`
	Status   []Sh5ExecStatus `json:"status"`
}

// Sh5ExecStatus - статус в запросе /api/sh5exec
type Sh5ExecStatus string

var (
	Sh5ExecStatusInsert Sh5ExecStatus = "Insert"
	Sh5ExecStatusModify Sh5ExecStatus = "Modify"
	Sh5ExecStatusDelete Sh5ExecStatus = "Delete"
)

// Sh5ExecRep - структура ответа на запрос /api/sh5exec
type Sh5ExecRep struct {
	ErrorCode  int       `json:"errorCode"`
	ErrMessage string    `json:"errMessage"`
	Version    string    `json:"Version"`
	UserName   string    `json:"UserName"`
	ActionName string    `json:"actionName"`
	ActionType string    `json:"actionType"`
	ShTable    []ShTable `json:"shTable"`
}

// ShTable - структура таблицы в ответе на запрос /api/sh5struct или /api/sh5exec
type ShTable struct {
	Head     Head            `json:"head"`
	RecCount int             `json:"recCount"`
	Original []FieldPath     `json:"original"`
	Fields   []FieldAlt      `json:"fields"`
	Values   [][]interface{} `json:"values"`
}

// Head - идентификатор таблицы
type Head string

var (
	HeadCodeGGROUP Head = "209" // Товарная группа
)

// Field - поле таблицы
type Field struct {
	Path FieldPath
	Type FieldType
	Size int
	Alt  FieldAlt
}

// FieldPath - оригинальное название поля
type FieldPath string

var (
	FIELD_1_RID  FieldPath = "1"         // RID
	FIELD_4_GUID FieldPath = "4"         // GUID
	FIELD_42     FieldPath = "42"        //
	FIELD_3      FieldPath = "3"         // Наименование
	FIELD_209_1  FieldPath = "209#1\\1"  // RID группы предка
	FIELD_209_4  FieldPath = "209#1\\4"  // GUID группы предка
	FIELD_209_42 FieldPath = "209#1\\42" //
	FIELD_209_3  FieldPath = "209#1\\3"  // Наименование
	FIELD_239    FieldPath = "239"       //
	FIELD_106_1  FieldPath = "106\\1"    // null
	FIELD_106_3  FieldPath = "106\\3"    // null
	FIELD_6      FieldPath = "6"         // null
)

// FieldType - тип поля
type FieldType string

var (
	FIELD_TYPE_UINT32 FieldType = "tUint32"
	FIELD_TYPE_GUID   FieldType = "tGuid"
	FIELD_TYPE_STRZ   FieldType = "tStrZ"
	FIELD_TYPE_STRP   FieldType = "tStrP"
)

// Альтернативное название поля
type FieldAlt string

var (
	FIELD_ALT_RID      FieldAlt = "Rid"      // Rid
	FIELD_ALT_GUID     FieldAlt = "Guid"     // Guid
	FIELD_ALT_NAME     FieldAlt = "Name"     // Name
	FIELD_ALT_MAXCOUNT FieldAlt = "MaxCount" // MaxCount
)
