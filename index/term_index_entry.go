package index

import (
	"github.com/xiiv13/address_match_recommend/utils"
)

// TermIndexEntry 索引条目
type TermIndexEntry struct {
	Key      string                   // 条目的key
	Items    []*TermIndexItem         // 每个条目下的所有索引对象
	Children map[rune]*TermIndexEntry // 子条目
}

func NewTermIndexEntry(text string) *TermIndexEntry {
	return &TermIndexEntry{
		Key:      text,
		Items:    make([]*TermIndexItem, 0),
		Children: make(map[rune]*TermIndexEntry),
	}
}

// BuildIndex 初始化倒排索引
func (tie *TermIndexEntry) BuildIndex(text []rune, pos int, item *TermIndexItem) {
	if len(text) == 0 || pos < 0 || pos >= len(text) {
		return
	}

	c := text[pos]
	_, ok := tie.Children[c]
	if !ok {
		tie.Children[c] = NewTermIndexEntry(utils.Head(text, pos+1))
	}
	if pos == len(text)-1 {
		tie.Children[c].Items = append(tie.Children[c].Items, item)
		return
	}
	tie.Children[c].BuildIndex(text, pos+1, item)
}
