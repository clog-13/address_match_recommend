package main

import (
	"bufio"
	"fmt"
	"github.com/xiiv13/address_match_recommend/core"
	"github.com/xiiv13/address_match_recommend/models"
	"io"
	"os"
	"strings"
)

func main() {
	filepath := "C:\\Users\\zx\\GolandProjects\\address_match_recommend\\resource\\test_addresses.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	addrs := make([]models.Address, 0)
	persister := models.NewAddressPersister()
	interpreter := core.NewAddressInterpreter(persister)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}

		importAddr := models.Address{}
		importAddr.RawText = strings.TrimSpace(line)
		importAddr.AddressText = strings.TrimSpace(line)
		interpreter.Interpret(&importAddr)
		addrs = append(addrs, importAddr)
	}

	// import addrs
	for _, v := range addrs {
		if v.Province != nil {
			v.ProvinceId = v.Province.ID
		}
		if v.City != nil {
			v.CityId = v.City.ID
		}
		if v.District != nil {
			v.DistrictId = v.District.ID
		}
		if v.Street != nil {
			v.StreetId = v.Street.ID
		}
		if v.Town != nil {
			v.TownId = v.Town.ID
		}
		if v.Village != nil {
			v.VillageId = v.Village.ID
		}
		models.DB.Create(&v)
	}
}
