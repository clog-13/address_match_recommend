package core

import . "github.com/xiiv13/address_match_recommend/models"

type AddressPersister struct {
	// REGION_TREE为中国国家区域对象，全国所有行政区域都以树状结构加载到REGION_TREE, 通过{@link Region#getChildren()}获取下一级列表
	RegionTree *Region

	// 按区域ID缓存的全部区域对象。
	RegionCache map[uint]*Region

	RegionLoaded              bool
	AddressIndexByHash        map[uint]struct{}
	AddressIndexByHashCreated bool
}

func (ap AddressPersister) GetRegion(id uint) *Region {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	return ap.RegionCache[id]
}

func (ap AddressPersister) loadRegions() {
	if ap.RegionLoaded {
		return
	}

	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip` from `bas_region` where id=1
	DB.Where("id =", 1).First(ap.RegionTree)

	ap.RegionCache = make(map[uint]*Region)
	ap.RegionCache[ap.RegionTree.ID] = ap.RegionTree
	ap.loadRegionChildren(ap.RegionTree)
	ap.RegionLoaded = true
}

func (ap AddressPersister) loadRegionChildren(parent *Region) {
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
			ap.RegionCache[child.ID] = child
			ap.loadRegionChildren(child)
		}
	}
}

func (ap AddressPersister) GetRootRegionChilden() []*Region {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	return ap.RegionTree.Children
}
