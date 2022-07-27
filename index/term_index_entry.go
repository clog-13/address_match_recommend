package index

import (
	"address_match_recommend/util"
)

type TermIndexEntry struct {
	Key      string
	Items    []TermIndexItem
	Children map[byte]TermIndexEntry
}

func (tie TermIndexEntry) buildIndex(text string, pos int, item TermIndexItem) {
	if len(text) == 0 || pos < 0 || pos >= len(text) {
		return
	}
	c := text[pos]
	if tie.Children == nil {
		tie.Children = make(map[byte]TermIndexEntry, 1)
		entry, ok := tie.Children[c]
		if !ok {
			entry = TermIndexEntry{
				Key:      util.Head(text, pos+1),
				Children: map[byte]TermIndexEntry{c: entry},
			}
		}
		if pos == len(text)-1 {
			entry.Items = append(entry.Items, item)
			return
		}
		entry.buildIndex(text, pos+1, item)
	}
}
