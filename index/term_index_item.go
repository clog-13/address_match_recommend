package index

import (
	. "github.com/xiiv13/address_match_recommend/models"
)

// TermIndexItem 索引对象
type TermIndexItem struct {
	Types TermEnum
	Value *Region
}

func NewTermIndexItem(t TermEnum, v *Region) *TermIndexItem {
	return &TermIndexItem{
		Types: t,
		Value: v,
	}
}
