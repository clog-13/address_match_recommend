package models

type AddressEntity struct {
	Id int64 `gorm:"primaryKey;comment:'用户ID'" json:"Id"`

	AddressText string `gorm:"type:varchar(128);comment:'完整地址'" json:"address_text"`
	Road        string `gorm:"type:varchar(32);comment:'道路信息'" json:"road"`
	RoadNum     string `gorm:"type:varchar(32);comment:'道路号'" json:"road_num"`
	BuildingNum string `gorm:"type:varchar(32);comment:'建筑信息'" json:"building_num"`

	Div Division
}

func NewAddrEntity(text string) *AddressEntity {
	return &AddressEntity{
		AddressText: text,
	}
}
