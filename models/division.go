package models

type Division struct {
	Id int64

	Province RegionEntity
	City     RegionEntity
	District RegionEntity
	Street   RegionEntity
	Town     RegionEntity
	Village  RegionEntity
}

// LeastRegion 获取最小一级有效行政区域对象
func (d Division) LeastRegion() RegionEntity {
	if !d.Village.IsNil() {
		return d.Village
	}
	if !d.Town.IsNil() {
		return d.Town
	}
	if !d.Street.IsNil() {
		return d.Street
	}
	if !d.District.IsNil() {
		return d.District
	}
	if !d.City.IsNil() {
		return d.City
	}
	return d.Province
}

func (d Division) GetTown() RegionEntity {
	if !d.Town.IsNil() {
		return d.Town
	}
	if d.Street.IsNil() {
		return RegionEntity{}
	}
	if d.Street.IsTown() {
		return d.Street
	}
	return RegionEntity{}
}

func (d Division) SetTown(value RegionEntity) {
	if value.IsNil() {
		d.Town = value
		return
	}
	switch value.Types {
	case TownRegion:
		d.Town = value
	case StreetRegion:
	case PlatformL4:
		d.Street = value
	default:
		d.Town = RegionEntity{}
	}
}
