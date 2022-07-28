package test

import (
	. "address_match_recommend/models"
	"testing"
)

func TestInsertAddress(t *testing.T) {
	addr := &Address{
		AddressText: "_clog_addr",
		Road:        "_clog_addr",
		RoadNum:     "_clog_addr",
		BuildingNum: "_clog_addr",
	}
	addr.Div = Division{}
	addr.Div.Province = &Region{
		Name:         "_clog_d_p",
		Alias:        "_clog_d_p",
		Types:        13,
		Children:     make([]*Region, 0),
		OrderedNames: make([]string, 0),
	}
	addr.Div.Province.Children = append(addr.Div.Province.Children, &Region{
		Name:  "pc",
		Alias: "pc",
	})
	addr.Div.Province.Children = append(addr.Div.Province.Children, &Region{
		Name:  "pc1",
		Alias: "pc1",
	})
	addr.Div.Town = &Region{
		Name:         "_clog_d_t",
		Alias:        "_clog_d_t",
		Types:        13,
		Children:     make([]*Region, 0),
		OrderedNames: make([]string, 0),
	}
	DB.Create(addr)
}
