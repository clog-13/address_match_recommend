package models

import (
	"github.com/lib/pq"
	"strings"
)

// Region 行政区域实体
type Region struct {
	ID         uint `gorm:"primaryKey;comment:行政区域ID" json:"ID"`
	DivisionID uint

	ParentID uint       `gorm:"type:uint;comment:完整地址" json:"region_parent_id"`
	Name     string     `gorm:"type:string;comment:区域名称" json:"region_name"`
	Alias    string     `gorm:"type:string;comment:区域别名" json:"region_alias"`
	Types    RegionEnum `gorm:"type:SMALLINT;comment:区域类型" json:"region_types"`

	Children     []*Region      `gorm:"-"`
	OrderedNames pq.StringArray `gorm:"-"`
	//_varchar OrderedNames pq.StringArray `gorm:"type:varchar(255)[]" json:"region_ordered_names"`
}

func (r *Region) IsTown() bool {
	switch r.Types {
	case CountryRegion:
		return true
	case StreetRegion:
		if r.Name == "" {
			return false
		}
		return len(r.Name) <= 4 &&
			(string(r.Name[len(r.Name)-1]) == "镇" || string(r.Name[len(r.Name)-1]) == "乡")
	}
	return false
}

// OrderedNameAndAlias 获取所有名称和别名列表，按字符长度倒排序。
func (r *Region) OrderedNameAndAlias() []string {
	if r.OrderedNames != null {
		return r.OrderedNames
	}

	r.OrderedNames = make([]string, 0)
	r.OrderedNames = append(r.OrderedNames, r.Name)
	tokens := make([]string, 0)
	if r.Alias != "" && len(strings.TrimSpace(r.Alias)) > 0 {
		tokens = strings.Split(strings.TrimSpace(r.Alias), ";")
	}
	for _, v := range tokens {
		if len(v) > 0 && len(strings.TrimSpace(v)) > 0 {
			r.OrderedNames = append(r.OrderedNames, strings.TrimSpace(v))
		}
	}

	// 按长度倒序
	exchanged := true
	endIndex := len(r.OrderedNames) - 1
	for exchanged && endIndex > 0 {
		exchanged = false
		for i := 0; i < endIndex; i++ {
			if len(r.OrderedNames[i]) < len(r.OrderedNames[i+1]) {
				temp := r.OrderedNames[i]
				r.OrderedNames[i] = r.OrderedNames[i+1]
				r.OrderedNames[i+1] = temp
				exchanged = true
			}
		}
		endIndex--
	}

	return r.OrderedNames
}

func (r *Region) Equal(t *Region) bool {
	return r.ParentID == t.ParentID && r.Name == t.Name &&
		r.Alias == t.Alias && r.Types == t.Types
}

func (r *Region) TableName() string {
	return "region"
}
