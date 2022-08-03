package models

type Document struct {
	ID int
	// 文档所有词条, 按照文档顺序, 未去重
	Terms    []*Term
	TermsMap map[string]*Term

	Province *Term
	City     *Term
	District *Term

	Street       *Term
	Town         *Term // 乡镇相关的词条信息
	Village      *Term
	Road         *Term // 道路信息
	RoadNum      *Term
	RoadNumValue int
}

func NewDocument(id int) Document {
	return Document{ID: id}
}

// GetTerm 获取词语对象。
func (d *Document) GetTerm(term string) *Term {
	if len(d.Terms) == 0 || d.Terms == nil {
		return nil
	}
	if d.TermsMap == nil {
		if d.TermsMap == nil {
			d.TermsMap = make(map[string]*Term)
			for _, v := range d.Terms {
				d.TermsMap[v.Text] = v
			}
		}
	}
	return d.TermsMap[term]
}
