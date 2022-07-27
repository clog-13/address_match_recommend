package index

import (
	"address_match_recommend/models"
)

type TermIndexItem struct {
	Types models.TermEnum
	Value any
}

func NewTermIndexItem(t models.TermEnum, v any) TermIndexItem {
	return TermIndexItem{
		Types: t,
		Value: v,
	}
}

func (tii TermIndexItem) IsNil() bool {

}
