package core

import (
	"fmt"
	"testing"
)

func TestFindsimilarAddress(t *testing.T) {
	querys, ok := FindsimilarAddress("北京海淀区丹棱街18号创富大厦1106", 5, true)
	// querys,ok := FindsimilarAddress("江苏连云港赣榆区江苏省赣榆县青口镇黄海路56号电影公司宿舍楼北楼1-901", 5, true)
	// querys,ok := FindsimilarAddress("江津连云港赣榆区江苏省赣榆县青口镇黄海路56号科技有限公司宿舍楼北楼", 5, true)
	//querys,ok := FindsimilarAddress("山东省济南市章丘区山东省章丘区明水开发区环路海尔公司", 5, true)
	//querys, ok := FindsimilarAddress("湖北武汉汉阳区汉阳经济技术开发区车城东路901号", 5, true)
	//querys, ok := FindsimilarAddress("湖西武汉汉阳区东荆河路海尔配套园武汉钣金有限公司", 5, true)
	//querys, ok := FindsimilarAddress("上海武汉汉阳区东荆河路海尔配套园武汉钣金有限公司", 5, true)
	//querys, ok := FindsimilarAddress("山东潍坊潍城区潍坊市潍城区望留西安村浮烟山风景区东侧", 5, true)
	if ok {
		fmt.Println("地址存在")
	} else {
		for _, v := range querys.SimiDocs {
			text := persister.LoadAddr(v.Doc.ID).RawText
			fmt.Println(text)
			fmt.Println(v.Similarity)
		}
	}

}

func TestProvince(t *testing.T) {
	//querys := FindsimilarAddress("湖北武汉汉阳区汉阳经济技术开发区车城东路901号", 5, true)
	//querys := FindsimilarAddress("湖西武汉汉阳区东荆河路海尔配套园武汉钣金有限公司", 5, true)
	//querys := FindsimilarAddress("湖南武汉汉阳区东荆河路海尔配套园武汉钣金有限公司", 5, true)
	querys, ok := FindsimilarAddress("北京武汉汉阳区东荆河路海尔配套园武汉钣金有限公司", 5, true)
	if ok {
		fmt.Println("地址存在")
	} else {
		for _, v := range querys.SimiDocs {
			text := persister.LoadAddr(v.Doc.ID).RawText
			fmt.Println(text)
			fmt.Println(v.Similarity)
		}
	}

}

func TestInsert(t *testing.T) {
	//t1 := "北京武汉汉阳区东荆河路海尔配套园武汉钣金有限公司"
	//ImportAddr(t1)
	//querys := FindsimilarAddress(t1, 5, true)
	//for _, v := range querys.SimiDocs {
	//	text := models.NewAddressPersister().LoadAddr(v.Doc.Id).RawText
	//	fmt.Println(text)
	//	fmt.Println(v.Similarity)
	//}
	t2 := "湖北武汉汉阳区天府软件园北区"
	ImportAddr(t2)
	querys, ok := FindsimilarAddress(t2, 5, true)
	if ok {
		fmt.Println("地址存在")
	} else {
		for _, v := range querys.SimiDocs {
			text := persister.LoadAddr(v.Doc.ID).RawText
			fmt.Println(text)
			fmt.Println(v.Similarity)
		}
	}
}
