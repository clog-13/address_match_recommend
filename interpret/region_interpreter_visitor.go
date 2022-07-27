package interpret

import (
	"address_match_recommend/index"
	"address_match_recommend/model"
	"address_match_recommend/models"
	"address_match_recommend/persist"
)

type RegionInterpreterVisitor struct {
	IsDebug        bool
	AmbiguousChars map[byte]struct{} // ‘市’
	Persister persist.AddressPersister

	CurrentLevel, DeepMostLevel            int
	CurrentPos, DeepMostPos                int //-1, -1
	FullMatchCount, DeepMostFullMatchCount int

	DeepMostDivision model.Division
	CurDivision      model.Division

	stack []index.TermIndexItem
}

func NewRegionInterpreterVisitor(ap persist.AddressPersister) RegionInterpreterVisitor{
	return RegionInterpreterVisitor{
		Persister: ap,
	}
}

func loadRegionChildren(parent model.RegionEntity) {
	//已经到最底层，结束
	if parent.IsNil() || parent.Types== models.StreetRegion || parent.Types== models.VillageRegion ||
		parent.Types== models.PlatformL4 ||parent.Types== models.TownRegion {
		return
	}
	//递归加载下一级
	children =
List<RegionEntity> children = this.regionDao.findByParent(parent.getId());
if(children!=null && children.size()>0){
parent.setChildren(children);
for(RegionEntity child : children) {
REGION_CACHE.put(child.getId(), child);
this.loadRegionChildren(child);
}
}
}