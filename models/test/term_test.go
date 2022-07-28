package test

import (
	. "address_match_recommend/models"
	"testing"
)

func TestTerm(t *testing.T) {
	DB.Create(&Term{
		Text:  "pure",
		Types: 10,
		Idf:   10,
		Ref: &Term{
			Text:  "pure_ref",
			Types: 10,
			Idf:   10,
			Ref:   nil,
		},
	})
}
