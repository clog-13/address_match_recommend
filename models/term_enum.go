package models

type TermEnum byte

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
)
