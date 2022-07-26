package model

type DeliveryAddr struct {
	Id uint `gorm:"primaryKey;comment:'用户ID'" json:"Id"`

	Country  string `gorm:"type:varchar(32);comment:'国家'" json:"country"`
	Province string `gorm:"type:varchar(32);comment:'省'" json:"province"`
	City     string `gorm:"type:varchar(32);comment:'市'" json:"city"`
	Barrio   string `gorm:"type:varchar(32);comment:'行政区'" json:"barrio"`
	Local    string `gorm:"type:varchar(128);comment:'具体地址'" json:"local"`
}

func NewAddr(country, province, city, barrio, local string) *DeliveryAddr {
	return &DeliveryAddr{
		Country:  country,
		Province: province,
		City:     city,
		Barrio:   barrio,
		Local:    local,
	}
}
