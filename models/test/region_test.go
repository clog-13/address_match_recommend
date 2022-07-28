package test

import (
	. "address_match_recommend/models"
	"testing"
)

func TestInsert(t *testing.T) {
	re := &RegionEntity{
		ParentID: 1,
		Name:     "-1",
		Alias:    "-1",
		Types:    1,
	}

	re.Children = make([]*RegionEntity, 0)
	child := &RegionEntity{
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
