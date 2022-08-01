package core

import (
	"bufio"
	"encoding/base64"
	"fmt"
	. "github.com/xiiv13/address_match_recommend/models"
	"io"
	"os"
)

type AddressPersister struct {
	// REGION_TREE为中国国家区域对象，全国所有行政区域都以树状结构加载到REGION_TREE
	// 通过 Region#getChildren() 获取下一级列表
	RegionTree Region

	// 按区域ID缓存的全部区域对象。
	RegionCache map[uint]Region

	RegionLoaded              bool
	AddressIndexByHash        map[uint]struct{}
	AddressIndexByHashCreated bool
}

func (ap *AddressPersister) GetAddress(id int) Address {
	var addr Address
	DB.Where("id = ?", id).Find(&addr)
	return addr
}

func (ap *AddressPersister) RootRegion() Region {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	return ap.RegionTree
}

func (ap *AddressPersister) GetRegion(id uint) *Region {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	res := ap.RegionCache[id]
	return &res
}

func (ap *AddressPersister) CreateRegion(region Region) {
	// insert into `bas_region`(`id`,`parent_id`,`name`,`type`,`zip`,`alias`)
	// values(#{id},#{parentId},#{name},#{type},#{zip},#{alias})
	DB.Create(&region)
}

func (ap *AddressPersister) FindRegion(parentID uint, name string) Region {
	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip` from `bas_region`
	// where parent_id=#{pid} and `name`=#{name}
	// order by id
	var region Region
	DB.Where("parent_id = ? AND  name = ?", parentID, name).Order("id").Find(&region)
	return region
}

func (ap *AddressPersister) loadRegions() {
	if ap.RegionLoaded {
		return
	}

	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip` from `bas_region` where id=1
	DB.Where("id =", 1).First(ap.RegionTree)

	ap.RegionCache = make(map[uint]Region)
	ap.RegionCache[ap.RegionTree.ID] = ap.RegionTree
	ap.loadRegionChildren(&ap.RegionTree)
	ap.RegionLoaded = true
}

func (ap *AddressPersister) loadRegionChildren(parent *Region) {
	// 已经到最底层，结束
	if parent == nil || parent.Types == StreetRegion || parent.Types == VillageRegion ||
		parent.Types == PlatformL4 || parent.Types == TownRegion {
		return
	}

	// parent.ID
	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip`
	// from `bas_region`
	// where parent_id=#{pid}
	// order by id
	var children []*Region
	DB.Order("id").Where("parent_id =", parent.ID).Find(&children)

	// 递归加载下一级
	if children != nil && len(children) > 0 {
		parent.Children = children
		for _, child := range children {
			ap.RegionCache[child.ID] = *child
			ap.loadRegionChildren(child)
		}
	}
}

func (ap *AddressPersister) GetRootRegionChilden() []*Region {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	return ap.RegionTree.Children
}

//func (ap *AddressPersister) importAddress(addrs []Address) int {
//	batchSize, count, imported, duplicate := 2000, 0, 0, 0
//	batch := make([]Address, 0)
//	for _,v:=range addrs {
//		if {
//
//		}
//
//	}
//}

func ReadDat(filepath string) {
	//path, _ := os.Getwd()
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

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

		//line = strings.TrimSpace(line)
		fmt.Println(line)
		fmt.Println(decode(line))
	}
}

func decode(dat string) string {
	str, err := base64.StdEncoding.DecodeString(dat)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return string(str)
	//decoder := base64.NewDecoder(base64.StdEncoding)
	//str, _ := base64.StdEncoding.DecodeString(dat)
}
