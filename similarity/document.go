package similarity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Document struct {
	// 文档所有词条, 按照文档顺序, 未去重
	Terms []*Term

	// 乡镇相关的词条信息
	Town    *Term
	Village *Term

	// 道路信息
	Road         *Term
	RoadNum      *Term
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
