package model

type Document struct {
	Id           int
	Terms        []Term
	TermsMap     map[string]Term
	Town         Term
	Village      Term
	Road         Term
	RoadNum      Term
	RoadNumValue int
}

func NewDocument(id int) Document {
	return Document{
		Id: id,
	}
}

// GetTerm 获取词语对象。
func (d Document) GetTerm(term string) Term {
	if len(term) == 0 || len(d.Terms) == 0 {
		return Term{}
	}
	if d.TermsMap == nil {
		d.buildMapCache()
	}
	return d.TermsMap[term]
}

func (d Document) buildMapCache() {
	if d.TermsMap == nil {
		d.TermsMap = make(map[string]Term, len(d.Terms))
	}
	for _, v := range d.Terms {
		d.TermsMap[v.Text] = v
	}
}
