package model

import "time"

type AddressEntity struct {
	Division
	SerialVersionUID int64 // 111198101809627685L
	Id               int
	Text             string
	Road             string
	RoadNum          string
	BuildingNum      string
	Hash             int
	// 仅保存到持久化仓库，从持久化仓库读取时不加载该属性
	RawTest    string
	Prop1      string
	Prop2      string
	CreateTime time.Time
}

func FormatAddressEntity(addr string) AddressEntity {
}

func (a AddressEntity) IsNil() bool {

}
