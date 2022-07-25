package model

import "time"

type AddressEntity struct {
	serialVersionUID int64 // 111198101809627685L
	id               int
	text             string
	road             string
	roadNum          string
	buildingNum      string
	hash             int
	// 仅保存到持久化仓库，从持久化仓库读取时不加载该属性
	rawTest    string
	prop1      string
	prop2      string
	createTime time.Time
}
