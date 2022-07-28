package similarity

import (
	. "address_match_recommend/models"
)

// Term 词条
type Term struct {
	Text  string
	Types TermEnum
	Idf   float64
	Ref   *Term
}

func NewTerm(types TermEnum, text string) *Term {
	return &Term{
		Types: types,
		Text:  text,
	}
}

func (t *Term) GetIdf() float64 {
	switch t.Types {
	case ProvinceTerm:
	case CityTerm:
	case DistrictTerm:
		t.Idf = 0.0
	case StreetTerm:
		t.Idf = 1.0
	}
	return t.Idf
}

func (t *Term) Equals(a *Term) bool {
	return t.Text == a.Text && t.Types == a.Types && t.Idf == a.Idf
}