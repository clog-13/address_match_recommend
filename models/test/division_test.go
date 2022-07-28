package test

import (
	. "address_match_recommend/models"
	"testing"
)

func TestInsertDivision(t *testing.T) {
	d := &Division{}
	d.Street = &Region{
		Name:  "xiiv_street",
		Alias: "xiiv_street",
		Types: 10,
	}
	DB.Create(d)
}
