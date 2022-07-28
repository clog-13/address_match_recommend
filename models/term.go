package models

// Term 词条
type Term struct {
	Id uint `gorm:"primaryKey;comment:词条ID" json:"ID"`

	Text  string   `gorm:"type:text;comment:词条字段" json:"term_text"`
	Types TermEnum `gorm:"type:SMALLINT;comment:词条类型" json:"term_types"`
	Idf   float64  `gorm:"type:float;comment:IDF" json:"term_idf"`

	DocumentID uint
	RefID      *uint `json:"term_parent_id"`
	Ref        *Term `gorm:"comment:相关联的词条引用" json:"term_ref"`
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

func (t Term) TableName() string {
	return "term"
}
