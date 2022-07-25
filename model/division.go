package model

type Division struct {
	province RegionEntity
	city     RegionEntity
	district RegionEntity
	street   RegionEntity
	town     RegionEntity
	village  RegionEntity
}

// 获取最小一级有效行政区域对象
func (d Division) leastRegion() RegionEntity {
	// TODO
	//if(hasVillage()) return getVillage();
	//if(hasTown()) return getTown();
	//if(hasStreet()) return getStreet();
	//if(hasDistrict()) return getDistrict();
	//if(hasCity()) return getCity();
	//return getProvince();
	return d.province
}

// TODO
//public RegionEntity getTown() {
//if(this.town!=null) return this.town;
//if(this.street==null) return null;
//return this.street.isTown() ? this.street : null;
//}
//
//public void setTown(RegionEntity value) {
//if(value==null) {
//this.town=null;
//return;
//}
//switch(value.getType()){
//case Town:
//this.town=value;
//return;
//case Street:
//case PlatformL4:
//this.street = value;
//return;
//default:
//this.town=null;
//}
//}
