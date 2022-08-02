package core

import (
	"fmt"
	"github.com/xiiv13/address_match_recommend/models"
	"testing"
)

var (
	addrs = []string{
		"抚顺顺城区将军桥【将军水泥厂住宅4-1-102】 (将军桥附近)",
		"辽宁沈阳于洪区沈阳市辽中县县城虹桥商厦西侧三单元外跨楼梯3-2-23-", // 冗余
		"北京海淀区丹棱街18号创富大厦1106",
		"江苏泰州兴化市昌荣镇【康琴网吧】 (昌荣镇附近)",
		"中国山东临沂兰山区小埠东社区居委会【绿杨榭公寓31-1-101 】 (绿杨榭公寓附近)",
		"山东济宁任城区金宇路【杨柳国际新城K8栋3单元1302】(杨柳国际新城·丽宫附近)",
	}
)

func TestAddressInterpreter_Interpret(t *testing.T) {
	interpreter := NewAddressInterpreter(models.NewAddressPersister())
	addr := &models.Address{}
	for _, v := range addrs {
		addr.AddressText = v
		interpreter.Interpret(addr)
		fmt.Println(addr, addr.Province.Name, addr.City.Name, addr.District.Name)
	}
}

func TestRemoveSpecialChars(t *testing.T) {
	ai := NewAddressInterpreter(models.NewAddressPersister())
	addr := &models.Address{}
	for _, v := range addrs {
		addr.AddressText = v
		ai.prepare(addr)
		// 提取建筑物号
		ai.extractBuildingNum(addr)
		// 去除特殊字符
		ai.removeSpecialChars(addr)
		fmt.Println(addr)
	}
}
