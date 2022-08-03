package models

type Address struct {
	Id int64 `gorm:"primaryKey;comment:地址ID" json:"ID"`

	RawText     string `gorm:"type:text;" json:"raw_text"`
	AddressText string `gorm:"type:text;" json:"address_text"`
	RoadText    string `gorm:"type:text;" json:"road"`
	RoadNum     string `gorm:"type:text;" json:"road_num"`
	BuildingNum string `gorm:"type:text;" json:"building_num"`

	ProvinceId, CityId, DistrictId, StreetId, VillageId, TownId int

	Province *Region `gorm:"-"`
	City     *Region `gorm:"-"`
	District *Region `gorm:"-"`
	Street   *Region `gorm:"-"`
	Town     *Region `gorm:"-"`
	Village  *Region `gorm:"-"`
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
