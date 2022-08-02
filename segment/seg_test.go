package segment

import (
	"fmt"
	"testing"
)

var (
	addrs = []string{
		"北京海淀区丹棱街18号创富大厦1106",
		"7号楼1单元102室",
		"九鼎2期B7号楼东数新都商贸购物中心CBD附近",
		"山东青岛李沧区虎山路街道北崂路993号东山峰景6号楼1单元602室",
		"辽宁省沈阳市沈河区东陵街道海上五月花三期302楼2",
		"安徽省合肥市瑶海区长江东路8号琥珀名城和园10栋2203",
		"河南省南阳市邓州市花洲街道新华东路刘庄村兴德旅社",
		"河北省唐山市路北区唐山高新技术产业开发区龙泽路于龙福南道交叉口南行50米维也纳音乐城",
	}
)

func TestSimpleSegmenter(t *testing.T) {
	simp := NewSimpleSegmenter()
	for _, v := range addrs {
		fmt.Println(simp.Segment(v))
	}
}

func TestGseSegment(t *testing.T) {
	g := NewGseSegment()
	for _, v := range addrs {
		fmt.Println(g.Stop(g.Cut(v)))
	}
}

func TestHMMSegment(t *testing.T) {
	g := NewGseSegment()
	for _, v := range addrs {
		fmt.Println(g.Cut(v, true))
	}
}
