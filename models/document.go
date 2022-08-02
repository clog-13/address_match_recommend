package models

type Document struct {
	Id uint

	// 文档所有词条, 按照文档顺序, 未去重
	Terms    []*Term
	TermsMap map[string]*Term

	TownId       uint
	Town         *Term // 乡镇相关的词条信息
	VillageId    uint
	Village      *Term
	RoadId       uint
	Road         *Term // 道路信息
	RoadNumId    uint
	RoadNum      *Term
	RoadNumValue int
}

func NewDocument(id uint) Document {
	return Document{Id: id}
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

func (d *Document) TableName() string {
	return "document"
}
