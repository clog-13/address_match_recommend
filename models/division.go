package models

// Division 行政区域集合
type Division struct {
	Id int64 `gorm:"primaryKey;comment:行政区规范ID" json:"ID"`

	//ProvinceId uint
	// `gorm:"foreignKey:division_id"`
	Province *RegionEntity
	//CityId     uint
	City *RegionEntity
	//DistrictId uint
	District *RegionEntity
	//StreetId   uint
	Street *RegionEntity
	//TownId     uint
	Town *RegionEntity
	//VillageId  uint
	Village *RegionEntity
}

// LeastRegion 获取最小一级有效行政区域对象
func (d Division) LeastRegion() *RegionEntity {
	if d.Village != nil {
		return d.Village
	}
	if d.Town != nil {
		return d.Town
	}
	if d.Street != nil {
		return d.Street
	}
	if d.District != nil {
		return d.District
	}
	if d.City != nil {
		return d.City
	}
	return d.Province
}

func (d Division) GetTown() *RegionEntity {
	if d.Town != nil {
		return d.Town
	}
	if d.Street != nil && d.Street.IsTown() {
		return d.Street
	}
	return nil
}

func (d Division) SetTown(value *RegionEntity) {
	if value == nil {
		d.Town = nil
		return
	}
	switch value.Types {
	case TownRegion:
		d.Town = value
	case StreetRegion:
	case PlatformL4:
		d.Street = value
	default:
		d.Town = nil
	}
}

func (d Division) TableName() string {
	return "division"
}
