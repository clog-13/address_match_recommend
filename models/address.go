package models

type Address struct {
	Id int64 `gorm:"primaryKey;comment:地址ID" json:"ID"`

	AddressText string `gorm:"type:text;comment:完整地址" json:"address_text"`
	Road        string `gorm:"type:text;comment:道路信息" json:"road"`
	RoadNum     string `gorm:"type:text;comment:道路号" json:"road_num"`
	BuildingNum string `gorm:"type:text;comment:建筑信息" json:"building_num"`

	Div Division `gorm:"embedded;embeddedPrefix:division_;comment:区域"`
}

func NewAddrEntity(text string) *Address {
	return &Address{
		AddressText: text,
	}
}

func (d *Address) TableName() string {
	return "address"
}
