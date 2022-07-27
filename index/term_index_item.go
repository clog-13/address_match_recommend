package index

import (
	"address_match_recommend/models"
)

type TermIndexItem struct {
	Types models.Enum
	Value any
}

func NewTermIndexItem(t models.Enum, v any) TermIndexItem {
	return TermIndexItem{
		Types: t,
		Value: v,
	}
}
