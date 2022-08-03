package main

import (
	"fmt"
	"github.com/xiiv13/address_match_recommend/core"
	"github.com/xiiv13/address_match_recommend/models"
)

func main() {
	persister := models.NewAddressPersister()
	arr := []string{
		//"北京海淀区丹棱街18号创富大厦1106",
		"江苏连云港赣榆区江苏省赣榆县青口镇黄海路56号电影公司宿舍楼北楼1-901",
		"江津连云港赣榆区江苏省赣榆县青口镇黄海路56号科技有限公司宿舍楼北楼",
		"山东省济南市章丘区山东省章丘区明水开发区环路海尔公司",
		"湖北武汉汉阳区汉阳经济技术开发区车城东路901号",
		"湖西武汉汉阳区东荆河路海尔配套园武汉钣金有限公司",
		"上海武汉汉阳区东荆河路海尔配套园武汉钣金有限公司",
		"山东潍坊潍城区潍坊市潍城区望留西安村浮烟山风景区东侧",
	}
	for _, a := range arr {
		fmt.Println("------------------------------------------")
		fmt.Println(a)
		querys, ok := core.FindsimilarAddress(a, 5, true)
		if ok {
			fmt.Println("地址存在")
		} else {
			for i, v := range querys.SimiDocs {
				text := persister.LoadAddr(v.Doc.Id).RawText
				fmt.Printf("%d.%s\n", i+1, text)
				fmt.Printf("  %g\n", v.Similarity)
			}
		}
	}
}
