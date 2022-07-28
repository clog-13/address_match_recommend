package index

import (
	. "address_match_recommend/models"
)

// TermIndexItem 索引对象
type TermIndexItem struct {
	Types TermEnum
	Value *RegionEntity
}

func NewTermIndexItem(t TermEnum, v *RegionEntity) *TermIndexItem {
	return &TermIndexItem{
		Types: t,
		Value: v,
	}
}
