package core

import (
	"github.com/xiiv13/address_match_recommend/index"
	. "github.com/xiiv13/address_match_recommend/models"
	"github.com/xiiv13/address_match_recommend/utils"
	"regexp"
	"strings"
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

	// 特殊字符1
	specialChars1 = []byte("　 \r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_-\\/")
	// 特殊字符2
	specialChars2 = []byte(`{}【】〈〉<>[]「」“”（）()`)

	// 匹配没有路号的情况
	reBuildingNum0 = regexp.MustCompile(`((路|街|巷)\d+号([\dA-Z一二三四五六七八九十][\\#\\-一－—/\\\\]|楼)?)?([\dA-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|楼|\\#楼?))?([一二三四五六七八九十东西南北甲乙丙\d]+([\\#\\-一－—/\\\\]|单元|门|梯|层|座|组))?(\d+([\\#\\-一－—/\\\\]|室|房|层|楼|号|户)?)?(\d+号?)?`)

	// 标准匹配building的模式：xx栋xx单元xxx
	// 山东青岛市南区宁夏路118号4号楼6单元202。如果正则模式开始位置不使用(路[0-9]+号)?, 则第一个符合条件的匹配结果是【118号4】,
	// 按照逻辑会将匹配结果及之后的所有字符当做building，导致最终结果为：118号4号楼6单元202
	// 所以需要先匹配 (路\d+号)?
	reBuildingNum1 = regexp.MustCompile(`((路|街|巷)\d+号)?([\dA-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?))?([一二三四五六七八九十东西南北甲乙丙\d]+(单元|门|梯|层|座))?(\d+(室|房)?)?`)

	// 校验building的模式。building1M能够匹配到纯数字等不符合条件的文本，使用building1V排除掉
	reBuildingNumV = regexp.MustCompile(`(栋|幢|橦|号楼|号|\\#|\\#楼|单元|室|房|门)+`)

	// 匹配building的模式：12-2-302，12栋3单元302
	reBuildingNum2 = regexp.MustCompile(`[A-Za-z\d]+([\\#\\-一－/\\\\]+[A-Za-z\d]+)+`)

	// 匹配building的模式：10组21号，农村地址
	reBuildingNum3 = regexp.MustCompile(`[0-9]+(组|通道)[A-Z0-9\\-一]+号?`)

	// 简单括号匹配
	bracketPattern = regexp.MustCompile(`([\\(（\\{\\<〈\\[【「][^\\)）\\}\\>〉\\]】」]*[\\)）\\}\\>〉\\]】」])`)

	// 道路信息
	reROAD = regexp.MustCompile(`^(?P<road>([\u4e00-\u9fa5]{2,6}(路|街坊|街|道|大街|大道)))(?P<ex>[甲乙丙丁])?(?P<roadnum>[0-9０１２３４５６７８９一二三四五六七八九十]+(号院|号楼|号大院|号|號|巷|弄|院|区|条|\\#院|\\#))?`)

	// 道路中未匹配到的building信息
	reRoadBuilding = regexp.MustCompile(`[0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?)`)

	// 村信息
	//P_TOWN1 = regexp.MustCompile(`^((?P<z>[\u4e00-\u9fa5]{2}(镇|乡))(?P<c>[\u4e00-\u9fa5]{1,3}村)?)`)
	//P_TOWN2 = regexp.MustCompile(`^((?P<z>[\u4e00-\u9fa5]{1,3}镇)?(?P<x>[\u4e00-\u9fa5]{1,3}乡)?(?P<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?)`)
	//P_TOWN3 = regexp.MustCompile(`^(?P<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?`)

	invalidTown           = make(map[string]struct{})
	invalidTownFollowings = make(map[string]struct{})
)

func init() {
	invalidTownFollowings["政府"] = struct{}{}
	invalidTownFollowings["大街"] = struct{}{}
	invalidTownFollowings["大道"] = struct{}{}
	invalidTownFollowings["社区"] = struct{}{}
	invalidTownFollowings["小区"] = struct{}{}
	invalidTownFollowings["小学"] = struct{}{}
	invalidTownFollowings["中学"] = struct{}{}
	invalidTownFollowings["医院"] = struct{}{}
	invalidTownFollowings["银行"] = struct{}{}
	invalidTownFollowings["中心"] = struct{}{}
	invalidTownFollowings["卫生"] = struct{}{}
	invalidTownFollowings["一小"] = struct{}{}
	invalidTownFollowings["一中"] = struct{}{}
	invalidTownFollowings["政局"] = struct{}{}
	invalidTownFollowings["企局"] = struct{}{}

	invalidTown["新村"] = struct{}{}
	invalidTown["外村"] = struct{}{}
	invalidTown["大村"] = struct{}{}
	invalidTown["后村"] = struct{}{}
	invalidTown["东村"] = struct{}{}
	invalidTown["南村"] = struct{}{}
	invalidTown["北村"] = struct{}{}
	invalidTown["西村"] = struct{}{}
	invalidTown["上村"] = struct{}{}
	invalidTown["下村"] = struct{}{}
	invalidTown["一村"] = struct{}{}
	invalidTown["二村"] = struct{}{}
	invalidTown["三村"] = struct{}{}
	invalidTown["四村"] = struct{}{}
	invalidTown["五村"] = struct{}{}
	invalidTown["六村"] = struct{}{}
	invalidTown["七村"] = struct{}{}
	invalidTown["八村"] = struct{}{}
	invalidTown["九村"] = struct{}{}
	invalidTown["十村"] = struct{}{}
	invalidTown["中村"] = struct{}{}
	invalidTown["街村"] = struct{}{}
	invalidTown["头村"] = struct{}{}
	invalidTown["店村"] = struct{}{}
	invalidTown["桥村"] = struct{}{}
	invalidTown["楼村"] = struct{}{}
	invalidTown["老村"] = struct{}{}
	invalidTown["户村"] = struct{}{}
	invalidTown["山村"] = struct{}{}
	invalidTown["才村"] = struct{}{}
	invalidTown["子村"] = struct{}{}
	invalidTown["旧村"] = struct{}{}
	invalidTown["文村"] = struct{}{}
	invalidTown["全村"] = struct{}{}
	invalidTown["和村"] = struct{}{}
	invalidTown["湖村"] = struct{}{}
	invalidTown["甲村"] = struct{}{}
	invalidTown["乙村"] = struct{}{}
	invalidTown["丙村"] = struct{}{}
	invalidTown["邻村"] = struct{}{}
	invalidTown["乡村"] = struct{}{}
	invalidTown["村二村"] = struct{}{}
	invalidTown["中关村"] = struct{}{}

	invalidTown["城乡"] = struct{}{}
	invalidTown["县乡"] = struct{}{}
	invalidTown["头乡"] = struct{}{}
	invalidTown["牌乡"] = struct{}{}
	invalidTown["茶乡"] = struct{}{}
	invalidTown["水乡"] = struct{}{}
	invalidTown["港乡"] = struct{}{}
	invalidTown["巷乡"] = struct{}{}
	invalidTown["七乡"] = struct{}{}
	invalidTown["站乡"] = struct{}{}
	invalidTown["西乡"] = struct{}{}
	invalidTown["宝乡"] = struct{}{}
	invalidTown["还乡"] = struct{}{}

	invalidTown["古镇"] = struct{}{}
	invalidTown["小镇"] = struct{}{}
	invalidTown["街镇"] = struct{}{}
	invalidTown["城镇"] = struct{}{}
	invalidTown["环镇"] = struct{}{}
	invalidTown["湾镇"] = struct{}{}
	invalidTown["岗镇"] = struct{}{}
	invalidTown["镇镇"] = struct{}{}
	invalidTown["场镇"] = struct{}{}
	invalidTown["新镇"] = struct{}{}
	invalidTown["乡镇"] = struct{}{}
	invalidTown["屯镇"] = struct{}{}
	invalidTown["大镇"] = struct{}{}
	invalidTown["南镇"] = struct{}{}
	invalidTown["店镇"] = struct{}{}
	invalidTown["铺镇"] = struct{}{}
	invalidTown["关镇"] = struct{}{}
	invalidTown["口镇"] = struct{}{}
	invalidTown["和镇"] = struct{}{}
	invalidTown["建镇"] = struct{}{}
	invalidTown["集镇"] = struct{}{}
	invalidTown["庙镇"] = struct{}{}
	invalidTown["河镇"] = struct{}{}
	invalidTown["村镇"] = struct{}{}
}

type AddressInterpreter struct {
	indexBuilder index.TermIndexBuilder
}

func NewAddressInterpreter(persister AddressPersister) *AddressInterpreter {
	return &AddressInterpreter{
		indexBuilder: index.NewTermIndexBuilder(persister, ignoringRegionNames),
	}
}

// Interpret 将地址进行标准化处理, 解析成 Address
func (ai AddressInterpreter) Interpret(entity *Address) {

	visitor := NewRegionInterpreterVisitor(persister)
	ai.interpret(entity, visitor)
}

func (ai AddressInterpreter) interpret(entity *Address, visitor index.TermIndexVisitor) {
	// 清洗下开头垃圾数据, 针对用户数据
	ai.prepare(entity)
	// 提取建筑物号
	ai.extractBuildingNum(entity)
	// 去除特殊字符
	ai.removeSpecialChars(entity)
	// 提取包括的数据
	brackets := ai.extractBrackets(entity)
	// 去除包括的特殊字符
	brackets = utils.Remove(brackets, specialChars2, "")
	entity.AddressText = utils.Remove(entity.AddressText, specialChars2, "")

	// 提取行政规划标准地址
	ai.extractRegion(entity, visitor)
	// 规整省市区街道等匹配的结果
	ai.removeRedundancy(entity, visitor)

	// 提取道路信息
	ai.extractRoad(entity)

	// TODO

	// ai.extractTownVillage(entity)

	entity.AddressText = regexp.MustCompile(`[0-9A-Za-z\\#]+(单元|楼|室|层|米|户|\\#)`).ReplaceAllString(entity.AddressText, "")
	entity.AddressText = regexp.MustCompile(`[一二三四五六七八九十]+(单元|楼|室|层|米|户)`).ReplaceAllString(entity.AddressText, "")
	if len(brackets) != 0 {
		entity.AddressText = entity.AddressText + brackets
		if len(entity.RoadText) != 0 { // 如果没有道路信息, 可能存在于 Brackets 中
			ai.extractRoad(entity)
		}
	}
}

func interprets(addrTextList []string, visitor RegionInterpreterVisitor) []Address {
	return nil
	//	if addrTextList == nil {
	//		return nil
	//	}
	//	numSuccess, numFail := 0, 0
	//	addresses := make([]Address, 0)
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
	//func interpretSimgle(addressText string, visitor RegionInterpreterVisitor) Address {

}

// 清洗下开头垃圾数据
func (ai AddressInterpreter) prepare(entity *Address) {
	if len(entity.AddressText) == 0 {
		return
	}

	// 去除开头的数字, 字母, 空格, 符号
	prefix := regexp.MustCompile("[ \\da-zA-Z\r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_\\-\\/]")
	strings.TrimPrefix(entity.AddressText, string(prefix.Find([]byte(entity.AddressText))))

	// 将地址中的 ー－—- 等替换为-
	replace := regexp.MustCompile(`[ー_－—/]|(--)`)
	replace.ReplaceAll([]byte(entity.AddressText), []byte("-"))
}

// 提取建筑物号
func (ai AddressInterpreter) extractBuildingNum(entity *Address) {
	if len(entity.AddressText) == 0 {
		return
	}
	found := false // 是否找到的标志
	matches := reBuildingNum0.FindAllStringSubmatch(entity.AddressText, -1)
	matchesIdx := reBuildingNum0.FindAllStringSubmatchIndex(entity.AddressText, -1)
	for i, match := range matches {
		if len(match[0]) == 0 {
			continue
		}

		var notEmptyCnt int
		for _, v := range match {
			if len(v) > 0 {
				notEmptyCnt++
			}
		}

		build := match[0]
		if notEmptyCnt > 3 && reBuildingNumV.MatchString(build) {
			pos := matchesIdx[i][0]
			if strings.HasPrefix(build, "路") || strings.HasPrefix(build, "街") ||
				strings.HasPrefix(build, "巷") {
				if strings.Contains(build, "号楼") {
					pos += strings.Index(build, "路") + 1
				} else {
					pos += strings.Index(build, "号") + 1
				}
				build = entity.AddressText[pos:matchesIdx[i][1]]
			}
			entity.BuildingNum = build
			entity.AddressText = entity.AddressText[:pos] + entity.AddressText[matchesIdx[i][1]:]
			found = true
			break
		}
	}

	if !found {
		matches := reBuildingNum1.FindAllStringSubmatch(entity.AddressText, -1)
		matchesIdx := reBuildingNum1.FindAllStringSubmatchIndex(entity.AddressText, -1)
		for i, match := range matches {
			if len(match[0]) == 0 {
				continue
			}

			var notEmptyCnt int
			for _, v := range match {
				if len(v) > 0 {
					notEmptyCnt++
				}
			}

			build := match[0]
			if notEmptyCnt > 3 && reBuildingNumV.MatchString(build) {
				pos := matchesIdx[i][0]
				if strings.HasPrefix(build, "路") || strings.HasPrefix(build, "街") ||
					strings.HasPrefix(build, "巷") {
					pos += strings.Index(build, "号") + 1
					build = entity.AddressText[pos:matchesIdx[i][1]]
				}
				entity.BuildingNum = build
				entity.AddressText = entity.AddressText[:pos] + entity.AddressText[matchesIdx[i][1]:]
				found = true
				break
			}
		}
	}

	if !found { // xx-xx-xx（xx栋xx单元xxx）
		match := reBuildingNum2.FindString(entity.AddressText)
		if len(match) > 0 {
			entity.BuildingNum = match
			pos := reBuildingNum2.FindStringIndex(entity.AddressText)
			entity.AddressText = entity.AddressText[:pos[0]] + entity.AddressText[pos[1]:]
			found = true
		}
	}

	if !found { // xx组xx号, xx通道xx号
		match := reBuildingNum3.FindString(entity.AddressText)
		if len(match) > 0 {
			entity.BuildingNum = match
			pos := reBuildingNum2.FindStringIndex(entity.AddressText)
			entity.AddressText = entity.AddressText[:pos[0]] + entity.AddressText[pos[1]:]
			found = true
		}
	}
}

// 去除特殊字符
func (ai AddressInterpreter) removeSpecialChars(entity *Address) {
	if len(entity.AddressText) == 0 {
		return
	}

	// 性能优化：使用String.replaceAll()和Matcher.replaceAll()方法性能相差不大，都比较耗时
	// 这种简单替换场景，自定义方法的性能比String.replaceAll()和Matcher.replaceAll()快10多倍接近20倍
	// 删除特殊字符
	text := utils.Remove(entity.AddressText, specialChars1, "")

	// 删除连续出现5个以上的数字 TODO: 可能会出现, 这个暂做这个处理
	text = utils.RemoveRepeatNum(text, 6)
	entity.AddressText = text

	// 去除building
	build := entity.BuildingNum
	if len(build) == 0 {
		return
	}
	build = utils.Remove(entity.AddressText, specialChars1, "-一－_#")
	build = utils.RemoveRepeatNum(text, 6)
	entity.BuildingNum = build
}

// 提取包括的数据
func (ai AddressInterpreter) extractBrackets(entity *Address) string {
	if len(entity.AddressText) == 0 {
		return ""
	}

	// 匹配出带有 `Brackets` 的文字, 拼接到 text 中
	matches := bracketPattern.FindAllString(entity.AddressText, -1)
	var found bool
	var sb strings.Builder
	for _, match := range matches {
		if len(match) <= 2 { // 如果没有文字
			continue
		}
		sb.WriteString(match[1 : len(match)-1])
		found = true
	}
	if found {
		entity.AddressText = bracketPattern.ReplaceAllString(entity.AddressText, "")
		return sb.String()
	}
	return ""
}

// 提取行政规划标准地址
func (ai AddressInterpreter) extractRegion(entity *Address, visitor index.TermIndexVisitor) {
	if len(entity.AddressText) == 0 {
		return
	}

	// 开始匹配
	visitor.Reset()
	ai.indexBuilder.DeepMostQuery(entity.AddressText, visitor)
	entity.Province = visitor.GetDevision().Province
	entity.City = visitor.GetDevision().City
	entity.District = visitor.GetDevision().District
	entity.Street = visitor.GetDevision().Street
	entity.Town = visitor.GetDevision().Town
	entity.Village = visitor.GetDevision().Village
	entity.AddressText = entity.AddressText[visitor.EndPosition()+1:]
}

// 规整省市区街道等匹配的结果
func (ai AddressInterpreter) removeRedundancy(entity *Address, visitor index.TermIndexVisitor) {
	if len(entity.AddressText) == 0 || entity.Province == nil || entity.City == nil {
		return
	}
	// 采用后序数组方式匹配省市区
	endIndex := len(entity.AddressText) - 2
	var i int
	for i < endIndex {
		visitor.Reset()
		ai.indexBuilder.DeepPosMostQuery(entity.AddressText, i, visitor)
		// 没有匹配上，或者匹配上的行政区域个数少于2个认当做无效匹配
		if visitor.MatchCount() < 2 || visitor.FullMatchCount() < 1 {
			i++
			continue
		}
		// 匹配上的省份、地级市不正确
		if entity.Province != visitor.GetDevision().Province || entity.City != visitor.GetDevision().City {
			i++
			continue
		}

		// TODO forkey_id

		devision := visitor.GetDevision() // 正确匹配
		// 修复 区 信息
		if entity.District == nil && devision.District != nil && devision.District.ParentID == entity.City.ID {
			entity.District = devision.District
		}
		// 修复 街道 信息
		if entity.District != nil && entity.Street == nil && devision.Street != nil &&
			devision.Street.ParentID == entity.District.ID {
			entity.Street = devision.Street
		}
		// 修复 乡镇 信息
		if entity.District != nil && entity.Town == nil && devision.Town != nil &&
			devision.Town.ParentID == entity.District.ID {
			entity.Town = devision.Town
		} else if entity.District != nil && entity.Town != nil && entity.Town.Equal(entity.Street) &&
			devision.Town != nil && !devision.Town.Equal(devision.Street) &&
			devision.Town.ParentID == entity.District.ID {
			entity.Town = devision.Town
		}
		if entity.District != nil && entity.Village == nil && devision.Village != nil &&
			devision.Village.ParentID == entity.District.ID {
			entity.Village = devision.Village
		}

		// 正确匹配上，删除
		entity.AddressText = entity.AddressText[visitor.EndPosition()+1:]
		endIndex = len(entity.AddressText)
		i = 0
	}
}

func (ai AddressInterpreter) extractRoad(entity *Address) {
	// 如果已经提取过了
	if len(entity.AddressText) == 0 || len(entity.RoadText) > 0 {
		return
	}
	matches := reROAD.FindStringSubmatch(entity.AddressText)
	if len(matches) > 0 {
		road := matches[reROAD.SubexpIndex("road")]
		ex := matches[reROAD.SubexpIndex("ex")]
		roadNum := ex + matches[reROAD.SubexpIndex("roadnum")]

		leftText := entity.AddressText[len(road)+len(roadNum):]
		if strings.HasPrefix(leftText, "小区") {
			return
		}
		entity.RoadText = fixRoad(road)
		// 仅包含【甲乙丙丁】单个汉字，不能作为门牌号
		if len(roadNum) == 1 {
			entity.AddressText = roadNum + leftText
		} else {
			entity.RoadNum = roadNum
			entity.AddressText = leftText
		}
		// 修复road中存在building的问题
		if len(entity.BuildingNum) == 0 {
			fixRoadBuilding(entity)
		}
	}
}

func fixRoad(road string) string {
	if len(road) == 0 || road[:len(road)/2] != road[len(road)/2:] {
		return road
	}
	return road[:len(road)/2]
}

func fixRoadBuilding(entity *Address) {
	if len(entity.AddressText) == 0 {
		return
	}
	// 最开始匹配, 先这样处理
	matchIdx := reRoadBuilding.FindStringIndex(entity.AddressText)
	if reRoadBuilding.MatchString(entity.AddressText) &&
		matchIdx[0] == 0 {
		entity.BuildingNum = entity.AddressText[matchIdx[0]:matchIdx[1]]
		entity.AddressText = entity.AddressText[:matchIdx[0]] + entity.AddressText[matchIdx[1]:]
	}
}
