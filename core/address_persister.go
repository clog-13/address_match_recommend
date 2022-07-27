package core

import . "address_match_recommend/models"

type AddressPersister struct {
	// REGION_TREE为中国国家区域对象，全国所有行政区域都以树状结构加载到REGION_TREE, 通过{@link RegionEntity#getChildren()}获取下一级列表
	RegionTree RegionEntity

	// 按区域ID缓存的全部区域对象。
	RegionCache map[int64]RegionEntity

	RegionLoaded              bool
	AddressIndexByHash        map[int64]struct{}
	AddressIndexByHashCreated bool
}

func (ap AddressPersister) GetRegion(id int64) RegionEntity {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	//if ap.RegionTree.IsNil() {
	//	panic("Region data not initialized")
	//}
	return ap.RegionCache[id]
}
func (ap AddressPersister) loadRegions() {
	if ap.RegionLoaded {
		return
	}

	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip` from `bas_region` where id=1
	ap.RegionTree = RegionEntity{}
	ap.RegionCache = make(map[int64]RegionEntity)
	ap.RegionCache[ap.RegionTree.Id] = ap.RegionTree
	ap.loadRegionChildren(ap.RegionTree)
	ap.RegionLoaded = true
}

func (ap AddressPersister) loadRegionChildren(parent RegionEntity) {
	// 已经到最底层，结束
	if parent.IsNil() || parent.Types == StreetRegion || parent.Types == VillageRegion ||
		parent.Types == PlatformL4 || parent.Types == TownRegion {
		return
	}
	// 递归加载下一级
	// parent.Id
	// select `id`,`parent_id`,`name`,`alias`,`type`,`zip`
	// from `bas_region`
	// where parent_id=#{pid}
	// order by id
	children := make([]RegionEntity, 0)
	if children != nil && len(children) > 0 {
		parent.Children = children
		for _, child := range children {
			ap.RegionCache[child.Id] = child
			ap.loadRegionChildren(child)
		}
	}
}

func (ap AddressPersister) RootRegion() RegionEntity {
	if !ap.RegionLoaded {
		ap.loadRegions()
	}
	return ap.RegionTree
}
