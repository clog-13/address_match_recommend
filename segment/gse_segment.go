package segment

import (
	"github.com/go-ego/gse"
)

type GseSegment struct {
	seg gse.Segmenter
}

func NewGseSegment() *GseSegment {
	g := new(GseSegment)
	g.seg.LoadDict() // 加载默认词典

	// 载入自定义词典
	err := g.seg.LoadDict("../resource/dic/region.dic")
	err = g.seg.LoadDict("../resource/dic/comminity.dic")
	err = g.seg.LoadStop("../resource/dic/stop.dic")
	if err != nil {
		panic(err)
	}

	return g
}

// Segment 简单分词器, 直接按单个字符切分, 连续出现的数字、英文字母会作为一个词条
func (g GseSegment) Segment(text string) []string {
	return g.seg.Stop(g.seg.Cut(text))
}
