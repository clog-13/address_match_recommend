package test

import (
	. "address_match_recommend/models"
	"testing"
)

func TestInsertAddress(t *testing.T) {
	addr := &AddressEntity{
		AddressText: "_clog_addr",
		Road:        "_clog_addr",
		RoadNum:     "_clog_addr",
		BuildingNum: "_clog_addr",
	}
	addr.Div = Division{}
	addr.Div.Province = &RegionEntity{
		Name:         "_clog_d_p",
		Alias:        "_clog_d_p",
		Types:        13,
		Children:     make([]*RegionEntity, 0),
		OrderedNames: make([]string, 0),
	}
	addr.Div.Province.Children = append(addr.Div.Province.Children, &RegionEntity{
		Name:  "pc",
		Alias: "pc",
	})
	addr.Div.Province.Children = append(addr.Div.Province.Children, &RegionEntity{
		Name:  "pc1",
		Alias: "pc1",
	})
	addr.Div.Town = &RegionEntity{
		Name:         "_clog_d_t",
		Alias:        "_clog_d_t",
		Types:        13,
		Children:     make([]*RegionEntity, 0),
		OrderedNames: make([]string, 0),
	}
	DB.Create(addr)
}
