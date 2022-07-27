package models

type Enum byte

const (
	UndefinedTerm = '0'
	ProvinceTerm  = '1'
	CityTerm      = '2'
	DistrictTerm  = '3'
	StreetTerm    = '4'
	TownTerm      = 'T'
	VillageTerm   = 'V'
	RoadTerm      = 'R'
	RoadNumTerm   = 'N'
	TextTerm      = 'X'
	IgnoreTerm    = 'I'

	UndefinedRegion    = 0
	CountryRegion      = 10
	ProvinceRegion     = 100
	ProvinceLevelCity1 = 150
	ProvinceLevelCity2 = 151
	CityRegion         = 200
	CityLevelDistrict  = 250
	DistrictRegion     = 300
	StreetRegion       = 450
	PlatformL4         = 460
	TownRegion         = 400
	VillageRegion      = 410
)
