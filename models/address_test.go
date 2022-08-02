package models

import (
	"fmt"
	"testing"
)

func TestInsertAddress(t *testing.T) {
	addr := &Address{
		AddressText: "<>clog<>addr",
		RoadText:    "<>clog<>addr",
		RoadNum:     "<>clog<>addr",
		BuildingNum: "<>clog<>addr",
		ProvinceId:  29,
		TownId:      22,
	}

	addr.Province = &Region{
		Name:  "<>clog<>d<>p",
		Alias: "<>clog<>d<>p",
	}
	addr.Province.Children = append(addr.Province.Children, &Region{
		Name:  "pc",
		Alias: "pc",
	})
	addr.Province.Children = append(addr.Province.Children, &Region{
		Name:  "pc1",
		Alias: "pc1",
	})

	addr.Town = &Region{
		Name:  "<>clog<>d<>t",
		Alias: "<>clog<>d<>t",
	}
	DB.Create(addr)
}

func TestQueryAddress(t *testing.T) {
	var addrs []Address
	//DB.Preload(clause.Associations).Find(&addrs)
	DB.Preload("Province").Preload("Town").Find(&addrs)
	fmt.Println(addrs[0])
	fmt.Println(addrs[0].Province)
	fmt.Println(addrs[0].Town)
}
