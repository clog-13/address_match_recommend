package similarity

type Document struct {
	// 文档所有词条, 按照文档顺序, 未去重
	Terms    []*Term
	TermsMap map[string]*Term

	// 乡镇相关的词条信息
	Town    *Term
	Village *Term

	// 道路信息
	Road         *Term
	RoadNum      *Term
	RoadNumValue int
}

func NewDocument() Document {
	return Document{}
}

// TODO nil

// GetTerm 获取词语对象。
func (d Document) GetTerm(term string) *Term {
	if len(d.Terms) == 0 || d.Terms == nil {
		return nil
	}
	if d.TermsMap == nil {
		if d.TermsMap == nil { // buildCache TODO go mutex
			d.TermsMap = make(map[string]*Term)
			for _, v := range d.Terms {
				d.TermsMap[v.Text] = v
			}
		}

	}
	return d.TermsMap[term]
}
