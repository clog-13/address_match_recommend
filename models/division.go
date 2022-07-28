package models

// Division 行政区域集合
type Division struct {
	Id uint `gorm:"primaryKey;comment:行政区规范ID" json:"ID"`

	//ProvinceId uint
	// `gorm:"foreignKey:division_id"`
	Province *Region
	//CityId     uint
	City *Region
	//DistrictId uint
	District *Region
	//StreetId   uint
	Street *Region
	//TownId     uint
	Town *Region
	//VillageId  uint
	Village *Region
}

// LeastRegion 获取最小一级有效行政区域对象
func (d Division) LeastRegion() *Region {
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

func (d Division) GetTown() *Region {
	if d.Town != nil {
		return d.Town
	}
	if d.Street != nil && d.Street.IsTown() {
		return d.Street
	}
	return nil
}

func (d Division) SetTown(value *Region) {
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
