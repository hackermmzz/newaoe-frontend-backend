package model

type MapInfo struct {
	Key   string `xorm:"pk varchar(255) notnull 'key'"`
	Value string `xorm:"text  'value'"`
}

func (p MapInfo) TableName() string {
	return "Map"
}
