package models

const (
	UndefinedTerm = '0'
	CountryTerm   = 'C'
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

// Term 词条
type Term struct {
	Id uint

	TermId uint
	Text   string
	Types  int
	Idf    float64

	Ref *Term
}

func NewTerm(types int, text string) *Term {
	return &Term{
		Types: types,
		Text:  text,
	}
}

func (t *Term) GetIdf() float64 {
	switch {
	case t.Types == ProvinceTerm:
	case t.Types == CityTerm:
	case t.Types == DistrictTerm:
		t.Idf = 0.0
	case t.Types == StreetTerm:
		t.Idf = 1.0
	}
	return t.Idf
}

func (t *Term) Equals(a *Term) bool {
	return t.Text == a.Text && t.Types == a.Types && t.Idf == a.Idf
}

func (t Term) TableName() string {
	return "term"
}
