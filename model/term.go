package model

type Term struct {
	RoadNumIdf float64 // 7
	Types      byte
	Text       string
	Idf        float64
	Ref        *Term
}

// NewTerm TODO
func NewTerm(types byte, text string) Term {
	return Term{
		Types: types,
		Text:  text,
	}
}

func (t Term) IsNil() bool {

}
