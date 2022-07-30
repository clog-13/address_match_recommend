package test

import (
	"fmt"
	"testing"
)

func TestSimilarity(t *testing.T) {
	// 一般匹配
	//text1 := "山东省沂水县四十里堡镇东艾家庄村205号"
	//text2 := "山东省沂水县四十里堡镇东艾家庄村206号"

	fmt.Println()

	// 带有building匹配
	//buiText1 := "湖南衡阳常宁市湖南省衡阳市常宁市泉峰街道泉峰街道消防大队南园小区A栋1单元601"
	//buiText2 := "湖南衡阳常宁市湖南省衡阳市常宁市泉峰街道泉峰街道消防大队南园小区A栋2单元601"
	//
	//// 特殊
	//speText1 := "山东青岛李沧区延川路116号绿城城园东区7号楼2单元802户"
	//speText2 := "山东青岛李沧区延川路绿城城园东区7-2-802"
	//
	//// 标准化
	//val addr1 = Geocoding.normalizing(text1)
	//val addr2 = Geocoding.normalizing(text2)
	//
	//println("相似度结果分析 >>>>>>>>> " + Geocoding.similarityWithResult(addr1, addr2))
}
