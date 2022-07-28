package core

import (
	"address_match_recommend/index"
	. "address_match_recommend/models"
	"regexp"
)

// AddressInterpreter 地址解析操作, 从地址文本中解析出省、市、区、街道、乡镇、道路等地址组成部分

var (
	ignoringRegionNames = []string{
		// JD, Tmall
		"其它区", "其他地区", "其它地区", "全境", "城区", "城区以内", "城区以外", "郊区", "县城内", "内环以内", "开发区", "经济开发区", "经济技术开发区",
		// ehaier (来自TMall或HP)
		"省直辖", "省直辖市县",
		// 其他
		"地区", "市区"}
	persister AddressPersister

	// 特殊字符1
	specialChars1 = []byte("　 \r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_-\\/")
	// 特殊字符2
	specialChars2 = []byte(`{}【】〈〉<>[]「」“”（）()`)
	// 匹配没有路号的情况
	P_BUILDING_NUM0 = regexp.MustCompile(`((路|街|巷)\d+号([\dA-Z一二三四五六七八九十][\\#\\-一－—/\\\\]|楼)?)?([0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|楼|\\#楼?)){0,1}([一二三四五六七八九十东西南北甲乙丙0-9]+([\\#\\-一－—/\\\\]|单元|门|梯|层|座|组))?([0-9]+([\\#\\-一－—/\\\\]|室|房|层|楼|号|户)?)?([0-9]+号?)?`)

	/**
	标准匹配building的模式：xx栋xx单元xxx
	注1：山东青岛市南区宁夏路118号4号楼6单元202。如果正则模式开始位置不使用(路[0-9]+号)?, 则第一个符合条件的匹配结果是【118号4】,
	按照逻辑会将匹配结果及之后的所有字符当做building，导致最终结果为：118号4号楼6单元202
	所以需要先匹配 (路[0-9]+号)?
	*/

	P_BUILDING_NUM1 = regexp.MustCompile(`((路|街|巷)\d+号)?([0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?)){0,1}([一二三四五六七八九十东西南北甲乙丙0-9]+(单元|门|梯|层|座))?([0-9]+(室|房)?)?`)

	// 校验building的模式。building1M能够匹配到纯数字等不符合条件的文本，使用building1V排除掉
	P_BUILDING_NUM_V = regexp.MustCompile(`(栋|幢|橦|号楼|号|\\#|\\#楼|单元|室|房|门)+`)

	// 匹配building的模式：12-2-302，12栋3单元302
	P_BUILDING_NUM2 = regexp.MustCompile(`[A-Za-z0-9]+([\\#\\-一－/\\\\]+[A-Za-z0-9]+)+`)

	// 匹配building的模式：10组21号，农村地址
	P_BUILDING_NUM3 = regexp.MustCompile(`[0-9]+(组|通道)[A-Z0-9\\-一]+号?`)

	// 简单括号匹配
	BRACKET_PATTERN = regexp.MustCompile(`(?P<bracket>([\\(（\\{\\<〈\\[【「][^\\)）\\}\\>〉\\]】」]*[\\)）\\}\\>〉\\]】」]))`)

	// 道路信息
	P_ROAD = regexp.MustCompile(`^(?P<road>([\u4e00-\u9fa5]{2,6}(路|街坊|街|道|大街|大道)))(?P<ex>[甲乙丙丁])?(?P<roadnum>[0-9０１２３４５６７８９一二三四五六七八九十]+(号院|号楼|号大院|号|號|巷|弄|院|区|条|\\#院|\\#))?`)
	// 道路中未匹配到的building信息
	P_ROAD_BUILDING = regexp.MustCompile(`[0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?)`)

	// 村信息
	P_TOWN1 = regexp.MustCompile(`^((?P<z>[\u4e00-\u9fa5]{2,2}(镇|乡))(?P<c>[\u4e00-\u9fa5]{1,3}村)?)`)
	P_TOWN2 = regexp.MustCompile(`^((?P<z>[\u4e00-\u9fa5]{1,3}镇)?(?P<x>[\u4e00-\u9fa5]{1,3}乡)?(?P<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?)`)
	P_TOWN3 = regexp.MustCompile(`^(?P<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?`)

	invalidTown           = make(map[string]struct{})
	invalidTownFollowings = make(map[string]struct{})
)

func init() {

}

type AddressInterpreter struct {
	indexBuilder index.TermIndexBuilder
}

func NewAddressInterpreter(persister AddressPersister, visitor TermIndexVisitor) {
	newAi := AddressInterpreter{}
	newAi.indexBuilder = index.NewTermIndexBuilder(persister, ignoringRegionNames)
}

// Interpret 将地址进行标准化处理, 解析成 AddressEntity
func (ai AddressInterpreter) Interpret(entity *AddressEntity) {
	visitor := NewRegionInterpreterVisitor(persister)
	ai.interpret(entity, visitor)
}

func (ai AddressInterpreter) interpret(entity *AddressEntity, visitor RegionInterpreterVisitor) {
	// 清洗下开头垃圾数据, 针对用户数据
	ai.prepare(entity)

	// extractBuildingNum, 提取建筑物号
	ai.extractBuildingNum(entity)

	//// 去除特殊字符
	//removeSpecialChars(entity)
	//// 提取包括的数据
	//var brackets = extractBrackets(entity)
	//// 去除包括的特殊字符
	//brackets = brackets.remove(specialChars2)
	//removeBrackets(entity)
	//// 提取行政规划标准地址
	//extractRegion(entity, visitor)
	//// 规整省市区街道等匹配的结果
	//removeRedundancy(entity, visitor)
	//// 提取道路信息
	//extractRoad(entity)

	/**
	  entity.text = entity.text!!.replace("[0-9A-Za-z\\#]+(单元|楼|室|层|米|户|\\#)", "")
	  entity.text = entity.text!!.replace("[一二三四五六七八九十]+(单元|楼|室|层|米|户)", "")
	  if (brackets.isNotEmpty()) {
	      entity.text = entity.text + brackets
	      // 如果没有道路信息, 可能存在于 Brackets 中
	      if (entity.road.isNullOrBlank()) extractRoad(entity)
	  }
	*/
}

// TODO

// 清洗下开头垃圾数据
func (ai AddressInterpreter) prepare(entity *AddressEntity) {
	if len(entity.AddressText) == 0 {
		return
	}
	// 去除开头的数字, 字母, 空格等
	//p, _ := regexp.Compile("[ \\da-zA-Z\r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_\\-\\/]")

	/**
	  entity.text = entity.text?.trimStart {
	      p.matcher("$it").find()
	  }

	  // 将地址中的 ー－—- 等替换为-
	  entity.text = entity.text?.replace(Regex("[ー_－—/]|(--)"), "-")
	*/
}

// 提取建筑物号
func (ai AddressInterpreter) extractBuildingNum(entity *AddressEntity) bool {
	if len(entity.AddressText) == 0 {
		return false
	}
	found := false // 是否找到的标志
	//var building string // 最后匹配的文本

	return found
}

//func interprets(addrTextList []string, visitor RegionInterpreterVisitor) []AddressEntity {
//	if addrTextList == nil {
//		return nil
//	}
//	numSuccess, numFail := 0, 0
//	addresses := make([]AddressEntity, 0)
//	for _, addrText := range addrTextList {
//		if len(addrText) == 0 {
//			continue
//		}
//		address := interpretSimgle(addrText, visitor)
//		if address.IsNil() || !address.City.IsNil() || !address.District.IsNil() {
//			numFail++
//			continue
//		}
//		numSuccess++
//		addresses = append(addresses, address)
//	}
//	return addresses
//}
//
//func interpretSimgle(addressText string, visitor RegionInterpreterVisitor) AddressEntity {
//}
