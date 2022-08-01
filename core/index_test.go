package core

import (
	"fmt"
	"github.com/xiiv13/address_match_recommend/utils"
	"testing"
)

func TestQueryIndex(t *testing.T) {
	persister := NewAddressPersister()
	builder := NewTermIndexBuilder(persister)
	visitor := NewRegionInterpreterVisitor(persister)

	text := "北京海淀区丹棱街18号创富大厦1106"
	//text := "山东青岛李沧区延川路116号绿城城园东区7号楼2单元802户"
	builder.DeepMostQuery(text, visitor)
	fmt.Println(text)
	fmt.Println(utils.Substring([]rune(text), 0, visitor.DeepMostPos))
	fmt.Println(visitor.DeepMostDivision)

	visitor.Reset()
	text = "青岛市南区"
	builder.DeepMostQuery(text, visitor)
	fmt.Println(text)
	fmt.Println(utils.Substring([]rune(text), 0, visitor.DeepMostPos))
	fmt.Println(visitor.DeepMostDivision)

	visitor.Reset()
	text = "新疆阿克苏地区阿拉尔市新苑祥和小区"
	builder.DeepMostQuery(text, visitor)
	fmt.Println(text)
	fmt.Println(utils.Substring([]rune(text), 0, visitor.DeepMostPos))
	fmt.Println(visitor.DeepMostDivision)

	visitor.Reset()
	text = "湖南湘潭市湘潭县易俗河镇中南建材市场"
	builder.DeepMostQuery(text, visitor)
	fmt.Println(text)
	fmt.Println(utils.Substring([]rune(text), 0, visitor.DeepMostPos))
	fmt.Println(visitor.DeepMostDivision)

	visitor.Reset()
	text = "广东从化区温泉镇新田村"
	builder.DeepMostQuery(text, visitor)
	fmt.Println(text)
	fmt.Println(utils.Substring([]rune(text), 0, visitor.DeepMostPos))
	fmt.Println(visitor.DeepMostDivision)
}
