package models

type Address struct {
	Id int64 `gorm:"primaryKey;comment:地址ID" json:"ID"`

	AddressText string `gorm:"type:text;comment:完整地址" json:"address_text"`
	RoadText    string `gorm:"type:text;comment:道路信息" json:"road"`
	RoadNum     string `gorm:"type:text;comment:道路号" json:"road_num"`
	BuildingNum string `gorm:"type:text;comment:建筑信息" json:"building_num"`

	ProvinceId, CityId, DistrictId, StreetId, VillageId, TownId uint

	Province *Region `gorm:"-"`
	City     *Region `gorm:"-"`
	District *Region `gorm:"-"`
	Street   *Region `gorm:"-"`
	Town     *Region `gorm:"-"`
	Village  *Region `gorm:"-"`

	//Province *Region `gorm:"foreignKey:division_id;references:province_id"`
	//City     *Region `gorm:"foreignKey:division_id;references:city_id"`
	//District *Region `gorm:"foreignKey:division_id;references:district_id"`
	//Street   *Region `gorm:"foreignKey:division_id;references:street_id"`
	//Town     *Region `gorm:"foreignKey:division_id;references:town_id"`
	//Village  *Region `gorm:"foreignKey:division_id;references:village_id"`
}

// LeastRegion 获取最小一级有效行政区域对象
func (d *Address) LeastRegion() *Region {
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

func (d *Address) GetTown() *Region {
	if d.Town != nil {
		return d.Town
	}
	if d.Street != nil && d.Street.IsTown() {
		return d.Street
	}
	return nil
}

func (d *Address) SetTown(value *Region) {
	if value == nil {
		d.Town = nil
		return
	}
	switch {
	case value.Types == TownRegion:
		d.Town = value
	case value.Types == StreetRegion:
	case value.Types == PlatformL4:
		d.Street = value
	default:
		d.Town = nil
	}
}

func (d *Address) TableName() string {
	return "address"
}
