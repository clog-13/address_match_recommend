package models

type AddressEntity struct {
	Id int64

	Text        string
	Road        string
	RoadNum     string
	BuildingNum string
	Hash        int
	division    Division
}

func NewAddrEntity(text string) AddressEntity {
	return AddressEntity{
		Text: text,
	}
}

func (a AddressEntity) IsNil() bool {

}
