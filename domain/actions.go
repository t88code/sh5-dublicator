package domain

type Action string

var (
	ActionInsert             Action = "insert"
	ActionDelete             Action = "delete"
	ActionModify             Action = "modify"
	ActionDeleteRidAndInsert Action = "delete_rid_and_insert"
)
