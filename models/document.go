package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Document struct {
	Id uint `gorm:"primaryKey;comment:文档ID" json:"ID"`

	// 文档所有词条, 按照文档顺序, 未去重
	TermsId uint    `gorm:"uniqueIndex"`
	Terms   []*Term `gorm:"foreignKey:document_id;references:terms_id;"`

	// 乡镇相关的词条信息
	TownId    uint  `gorm:"uniqueIndex"`
	Town      *Term `gorm:"foreignKey:document_id;references:town_id;"`
	VillageId uint  `gorm:"uniqueIndex"`
	Village   *Term `gorm:"foreignKey:document_id;references:village_id;"`

	// 道路信息
	RoadId       uint  `gorm:"uniqueIndex"`
	Road         *Term `gorm:"foreignKey:document_id;references:road_id;"`
	RoadNumId    uint  `gorm:"uniqueIndex"`
	RoadNum      *Term `gorm:"foreignKey:document_id;references:road_num_id;"`
	RoadNumValue int

	TermsMap map[string]*Term `gorm:"-"` // 数据库不存储
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

// Args 参数
type Args map[string]*Term

// Scan Scanner
func (args Args) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	return json.Unmarshal(b, &args)
}

// Value Valuer
func (args Args) Value() (driver.Value, error) {
	if args == nil {
		return nil, nil
	}

	return json.Marshal(args)
}

func (d Document) TableName() string {
	return "document"
}
