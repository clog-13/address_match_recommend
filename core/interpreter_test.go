package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xiiv13/address_match_recommend/models"
	"github.com/xiiv13/address_match_recommend/utils"
	"testing"
)

// TODO

var (
	addrs = []string{
		"()四{}川{aa}(bb)成（）都（cc）武[]侯[dd]区【】武【ee】侯<>大<ff>道〈〉铁〈gg〉佛「」段「hh」千盛百货对面200米金履三路288号绿地圣路易名邸[]",
		//"抚顺顺城区将军桥【将军水泥厂住宅4-1-102】 (将军桥附近)",
		//"辽宁沈阳于洪区沈阳市辽中县县城虹桥商厦西侧三单元外跨楼梯3-2-23-", // 冗余
		//"北京海淀区丹棱街18号创富大厦1106",
		//"江苏泰州兴化市昌荣镇【康琴网吧】 (昌荣镇附近)",
		//"中国山东临沂兰山区小埠东社区居委会【绿杨榭公寓31-1-101 】 (绿杨榭公寓附近)",
		//"山东济宁任城区金宇路【杨柳国际新城K8栋3单元1302】(杨柳国际新城·丽宫附近)",
	}

	nochar = []string{
		"山东青岛四方区（撤）山东省四方区人民路285号3号楼901室",
		"黑龙江哈尔滨南岗区文林街11号A栋901室",
		"山东青岛市北区错埠岭二路43号楼四单元901户",
		"山东青岛李沧区青岛市李沧区京口路64号1号楼2单元901",
		"山东青岛李沧区青岛市李沧区大崂路1024号三单元901",
		"辽宁大连甘井子区 大连甘井子区椒北路2号楼1单元4-901",
		"山东青岛李沧区振华路124号1单元901",
		"辽宁大连沙河口区沿河街22号1-704 华业玫瑰东方 5号楼1单元-901号",
	}
)

func TestExtraRoad(t *testing.T) {
	interpreter := NewAddressInterpreter(models.NewAddressPersister())
	arr := []string{
		"山东青岛即墨市龙山镇官庄村即墨市龙山街道办事处管庄村",               // 镇宁路
		"河北省石家庄市鹿泉市镇宁路贺庄回迁楼1号楼1单元602室",             // 镇宁路
		"北京北京海淀区北京市海淀区万寿路翠微西里13号楼1403室",            // 万寿路
		",海南海南省直辖市县定安县见龙大道财政局宿舍楼702",               // 见龙大道
		"河北石家庄长安区南村镇强镇街51号南村工商管理局",                 // 强镇街
		"吉林长春绿园区长春汽车产业开发区（省级）（特殊乡镇）长沈路1000号力旺格林春天", // 长沈路
	}
	for _, v := range arr {
		addr := &models.Address{AddressText: v}
		interpreter.Interpret(addr)
		fmt.Println(addr.RoadNum)
	}

}

func TestAddressInterpreter_Interpret(t *testing.T) {
	interpreter := NewAddressInterpreter(models.NewAddressPersister())
	addr := &models.Address{}
	for _, v := range nochar {
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

func TestExtractBuildingNum(t *testing.T) {
	ai := NewAddressInterpreter(models.NewAddressPersister())
	for _, v := range nochar {
		addr := &models.Address{AddressText: v}
		// 清洗下开头垃圾数据, 针对用户数据
		ai.prepare(addr)
		// 提取建筑物号
		ai.extractBuildingNum(addr)
		fmt.Println(addr.AddressText, addr.BuildingNum)
	}
}

func TestExtractBrackets(t *testing.T) {
	interpreter := NewAddressInterpreter(models.NewAddressPersister())
	addr := &models.Address{}
	addr.AddressText = "()四{}川{aa}(bb)成（）都（cc）武[]侯[dd]区【】武【ee】侯<>大<ff>道〈〉铁〈gg〉佛「」段「hh」千盛百货对面200米金履三路288号绿地圣路易名邸[]"
	brackets := interpreter.extractBrackets(addr)
	assert.Equal(t, brackets, "aabbccddeeffgghh")
	assert.Equal(t, addr.AddressText, "四川成都武侯区武侯大道铁佛段千盛百货对面200米金履三路288号绿地圣路易名邸")

	// FAIL TODO
	//addr.AddressText = "四川成都(武[]侯区武侯大道铁佛{aa}段千)盛百货对面200米金履三【bb】路288号绿地圣路易名邸"
	//brackets = interpreter.extractBrackets(addr)
	//assert.Equal(t, brackets, "aabb")
	//assert.Equal(t, addr.AddressText, "四川成都盛百货对面200米金履三路288号绿地圣路易名邸")
}

func TestSpecialChar(t *testing.T) {
	addr := &models.Address{
		AddressText: "(四)川成都武侯区武侯大【】〈〉<>[]「」道铁佛段千盛百货\\/ \r\n\t对面200米金履三路288号绿地610015圣路易名邸",
	}
	text := utils.Remove([]rune(addr.AddressText), specialChars1, "")
	text = utils.RemoveRepeatNum([]rune(text), 6)
	text = utils.Remove([]rune(text), specialChars2, "")
	assert.Equal(t, text, "四川成都武侯区武侯大道铁佛段千盛百货对面200米金履三路288号绿地圣路易名邸")
}
