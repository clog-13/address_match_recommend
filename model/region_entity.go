package model

import (
	"address_match_recommend/enum"
	"strings"
)

const ()

// RegionEntity 行政区域实体
type RegionEntity struct {
	serialVersionUID int64 // -111163973997033386L

	id           int64
	parentId     int64
	name         string
	alias        string
	types        int // RegionType enum
	zip          string
	children     []RegionEntity
	orderedNames []string
}

func (r RegionEntity) IsTown() bool {
	switch r.types {
	case enum.Country:
		return true
	case enum.Street:
		if r.name == "" {
			return false
		}
		return len(r.name) <= 4 &&
			(string(r.name[len(r.name)-1]) == "镇" || string(r.name[len(r.name)-1]) == "乡")
	}
	return false
}

// OrderedNameAndAlias 获取所有名称和别名列表，按字符长度倒排序。
func (r RegionEntity) OrderedNameAndAlias() []string {
	if r.orderedNames == nil {
		return r.orderedNames
	}
	r.buildOrderedNameAndAlias()
	return r.orderedNames
}

func (r RegionEntity) buildOrderedNameAndAlias() {
	if r.orderedNames != nil {
		return
	}
	tokens := make([]string, 0)
	if r.alias != "" && len(strings.TrimSpace(r.alias)) > 0 {
		tokens = strings.Split(strings.TrimSpace(r.alias), ";")
	}
	if tokens == nil || len(tokens) <= 0 {
		r.orderedNames = make([]string, 1)
	} else {
		r.orderedNames = make([]string, len(tokens)+1)
	}
	r.orderedNames = append(r.orderedNames, r.name)
	if tokens != nil {
		for _, v := range tokens {
			if v == "" || len(strings.TrimSpace(v)) <= 0 {
				continue
			}
			r.orderedNames = append(r.orderedNames, strings.TrimSpace(v))
		}
	}

	exchanged := true
	endIndex := len(r.orderedNames) - 1
	for exchanged && endIndex > 0 {
		exchanged = false
		for i := 0; i < endIndex; i++ {
			if len(r.orderedNames[i]) < len(r.orderedNames[i+1]) {
				temp := r.orderedNames[i]
				r.orderedNames[i] = r.orderedNames[i+1]
				r.orderedNames[i+1] = temp
				exchanged = true
			}
		}
		endIndex--
	}
}
