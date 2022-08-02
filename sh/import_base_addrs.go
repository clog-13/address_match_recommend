package sh

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
	filepath := "../resource/test_addresses.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	addrs := make([]models.Address, 0)
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				return
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}

		persister := models.NewAddressPersister()
		interpreter := core.NewAddressInterpreter(persister)
		importAddr := models.Address{}
		importAddr.RoadText = strings.TrimSpace(line)
		interpreter.Interpret(&importAddr)
		addrs = append(addrs, importAddr)
	}

	// import addrs
	for _, v := range addrs {
		v.ProvinceId = v.Province.ID
		v.CityId = v.City.ID
		v.DistrictId = v.District.ID
		v.StreetId = v.Street.ID
		v.TownId = v.Town.ID
		v.VillageId = v.Village.ID

	}
}
