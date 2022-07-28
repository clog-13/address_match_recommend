package models

type Document struct {
	Id uint `gorm:"primaryKey;"`

	// 文档所有词条, 按照文档顺序, 未去重
	Terms    []*Term          `gorm:"-"`
	TermsMap map[string]*Term `gorm:"-"`

	TownId       uint
	Town         *Term `gorm:"foreignKey:term_id;references:town_id"` // 乡镇相关的词条信息
	VillageId    uint
	Village      *Term `gorm:"foreignKey:term_id;references:village_id"`
	RoadId       uint
	Road         *Term `gorm:"foreignKey:term_id;references:road_id"` // 道路信息
	RoadNumId    uint
	RoadNum      *Term `gorm:"foreignKey:term_id;references:road_num_id"`
	RoadNumValue int
}

func NewDocument() Document {
	return Document{}
}

// GetTerm 获取词语对象。
func (d Document) GetTerm(term string) *Term {
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

func (d Document) TableName() string {
	return "document"
}
