package models

import (
	"strings"
)

// RegionEntity 行政区域实体
type RegionEntity struct {
	Id int64

	ParentId     int64
	Name         string
	Alias        string
	Types        int // RegionType enum
	Zip          string
	Children     []*RegionEntity
	OrderedNames []string
}

func (r RegionEntity) IsTown() bool {
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
func (r RegionEntity) OrderedNameAndAlias() []string {
	if r.OrderedNames == nil {
		return r.OrderedNames
	}
	r.buildOrderedNameAndAlias()
	return r.OrderedNames
}

func (r RegionEntity) buildOrderedNameAndAlias() {
	if r.OrderedNames != nil {
		return
	}
	tokens := make([]string, 0)
	if r.Alias != "" && len(strings.TrimSpace(r.Alias)) > 0 {
		tokens = strings.Split(strings.TrimSpace(r.Alias), ";")
	}
	if tokens == nil || len(tokens) <= 0 {
		r.OrderedNames = make([]string, 1)
	} else {
		r.OrderedNames = make([]string, len(tokens)+1)
	}
	r.OrderedNames = append(r.OrderedNames, r.Name)
	if tokens != nil {
		for _, v := range tokens {
			if v == "" || len(strings.TrimSpace(v)) <= 0 {
				continue
			}
			r.OrderedNames = append(r.OrderedNames, strings.TrimSpace(v))
		}
	}

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
}

func (r *RegionEntity) Equal(t *RegionEntity) bool {
	return r.ParentId == t.ParentId && r.Name == t.Name && r.Alias == t.Alias &&
		r.Types == t.Types && r.Zip == t.Zip
}
