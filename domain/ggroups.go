package domain

// GGroups - товарная группа - 209
type GGroups struct {
	ID              uint32 `db:"ID"`
	Rid             uint32 `db:"Rid"`
	Guid            string `db:"Guid"`
	FIELD_42        uint32 `db:"FIELD_42"`
	Name            string `db:"Name"`
	Attr            uint32 `db:"Attr"`
	Rid_PARENT      uint32 `db:"Rid_PARENT"`
	Guid_PARENT     string `db:"Guid_PARENT"`
	FIELD_42_PARENT uint32 `db:"FIELD_42_PARENT"`
	Name_PARENT     string `db:"Name_PARENT"`
	Attr_PARENT     uint32 `db:"Attr_PARENT"`
	UserGroup       uint64 `db:"UserGroup"`
	Rid_106         uint32 `db:"Rid_106"`
	Name_106        string `db:"Name_106"`
	Action          Action `db:"Action"`
	TableName       string `db:"TableName"`
}
