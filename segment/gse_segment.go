package segment

import (
	"github.com/go-ego/gse"
)

type GseSegment struct {
	seg gse.Segmenter
}

func NewGseSegment() gse.Segmenter {
	//g := new(GseSegment)
	g, err := gse.New("zh_s,../resource/dic/region.dic,../resource/dic/community.dic") // 加载默认词典

	// 载入自定义词典
	//err := g.seg.LoadDict("../resource/dic/region.dic")
	//err = g.seg.LoadDict("../resource/dic/community.dic")

	err = g.LoadStop("../resource/dic/stop.dic")
	if err != nil {
		panic(err)
	}

	return g
}
