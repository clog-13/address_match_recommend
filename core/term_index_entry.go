package core

import (
	"github.com/xiiv13/address_match_recommend/utils"
)

// TermIndexEntry 索引条目
type TermIndexEntry struct {
	Key      string                   // 条目的key
	Items    []*TermIndexItem         // 每个条目下的所有索引对象
	Children map[byte]*TermIndexEntry // 子条目
}

func NewTermIndexEntry() *TermIndexEntry {
	return new(TermIndexEntry)
}

// BuildIndex 初始化倒排索引
func (tie TermIndexEntry) BuildIndex(text string, pos int, item *TermIndexItem) {
	if len(text) == 0 || pos < 0 || pos >= len(text) {
		return
	}

	c := text[pos]
	entry, ok := tie.Children[c]
	if !ok {
		entry = NewTermIndexEntry()
		entry.Key = utils.Head(text, pos+1)

		tie.Children[c] = entry
	}
	if pos == len(text)-1 {
		entry.Items = append(entry.Items, item)
		return
	}
	entry.BuildIndex(text, pos+1, item)
}
