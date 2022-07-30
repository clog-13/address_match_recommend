package test

import (
	. "github.com/xiiv13/address_match_recommend/models"
	"testing"
)

func TestInsert(t *testing.T) {
	re := &Region{
		ParentID: 1,
		Name:     "-1",
		Alias:    "-1",
		Types:    1,
	}

	re.Children = make([]*Region, 0)
	child := &Region{
		ParentID: 2,
		Name:     "-2",
		Alias:    "-2",
		Types:    2,
	}
	re.Children = append(re.Children, child)

	re.OrderedNames = make([]string, 0)
	re.OrderedNames = append(re.OrderedNames, "-1")

	DB.Create(re)
}
