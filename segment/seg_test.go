package segment

import (
	"fmt"
	"testing"
)

const (
	text1 = "7号楼1单元102室"
	text2 = "九鼎2期B7号楼东数新都商贸购物中心CBD附近"
)

func TestSimpleSegmenter(t *testing.T) {
	simp := SimpleSegmenter{}
	fmt.Println(simp.Segment(text1))
	fmt.Println(simp.Segment(text2))
}

func TestGseSegment(t *testing.T) {
	gse := NewGseSegment()
	fmt.Println(gse.Segment(text1))
	fmt.Println(gse.Segment(text2))
}
