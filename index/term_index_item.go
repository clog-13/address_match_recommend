package index

import "address_match_recommend/enum"

type TermIndexItem struct {
	Types enum.Enum
	Value any
}

func NewTermIndexItem(t enum.Enum, v any) TermIndexItem {
	return TermIndexItem{
		Types: t,
		Value: v,
	}
}
