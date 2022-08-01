package core

import (
	. "github.com/xiiv13/address_match_recommend/models"
)

// TermIndexItem 索引对象
type TermIndexItem struct {
	Types int
	Value *Region
}

func NewTermIndexItem(t int, v *Region) *TermIndexItem {
	return &TermIndexItem{
		Types: t,
		Value: v,
	}
}
