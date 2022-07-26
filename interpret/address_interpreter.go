package interpret

import (
	"address_match_recommend/index"
	. "address_match_recommend/model"
	"address_match_recommend/persist"
	"regexp"
	"strings"
)

// AddressInterpreter 地址解析操作, 从地址文本中解析出省、市、区、街道、乡镇、道路等地址组成部分

var (
	termIndex index.TermIndexBuilder
	persister persist.AddressPersister

	specialChars1         = []byte(" \r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_-\\/")
	invalidTown           = make(map[string]struct{})
	invalidTownFollowings = make(map[string]struct{})

	BRACKET_PATTERN = "([\\(（\\{\\<〈\\[【「][^\\)）\\}\\>〉\\]】」]*[\\)）\\}\\>〉\\]】」])"

	// PBuildingNum1 匹配building的模式：xx栋xx单元xxx
	// 山东青岛市南区宁夏路118号4号楼6单元202。如果正则模式开始位置不使用(路[0-9]+号)?，则第一个符合条件的匹配结果是【118号4】，
	// 按照逻辑会将匹配结果及之后的所有字符当做building，导致最终结果为：118号4号楼6单元202
	PBuildingNum1          = "((路|街|巷)[0-9]+号)?([0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?)){0,1}([一二三四五六七八九十东西南北甲乙丙0-9]+(单元|门|梯|层|座))?([0-9]+(室|房)?)?"
	RegexpPBuildingNum1, _ = regexp.Compile(`((路|街|巷)[0-9]+号)?([0-9A-Z一二三四五六七八九十]+(栋|橦|幢|座|号楼|号|\\#楼?)){0,1}([一二三四五六七八九十东西南北甲乙丙0-9]+(单元|门|梯|层|座))?([0-9]+(室|房)?)?`)

	// P_BUILDING_NUM_V 校验building的模式。building1M能够匹配到纯数字等不符合条件的文本，使用building1V排除掉
	P_BUILDING_NUM_V = "(栋|幢|橦|号楼|号|\\#|\\#楼|单元|室|房|门)+"

	// P_BUILDING_NUM2 匹配building的模式：12-2-302，12栋3单元302
	P_BUILDING_NUM2 = "[A-Za-z0-9]+([\\#\\-一－/\\\\]+[A-Za-z0-9]+)+"

	// P_BUILDING_NUM3 匹配building的模式：10组21号，农村地址
	P_BUILDING_NUM3 = "[0-9]+组[0-9\\-一]+号?"

	P_TOWN1_Z = "[\u4e00-\u9fa5]{2,2}(镇|乡)"
	P_TOWN1_C = "[\u4e00-\u9fa5]{1,3}村"

	P_TOWN2 = "^((?<z>[\u4e00-\u9fa5]{1,3}镇)?(?<x>[\u4e00-\u9fa5]{1,3}乡)?(?<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?)"
	P_TOWN3 = "^(?<c>[\u4e00-\u9fa5]{1,3}村(?!(村|委|公路|(东|西|南|北)?(大街|大道|路|街))))?"
	P_ROAD  = "^(?<road>([\u4e00-\u9fa5]{2,4}(路|街坊|街|道|大街|大道)))(?<ex>[甲乙丙丁])?(?<roadnum>[0-9０１２３４５６７８９一二三四五六七八九十]+(号院|号楼|号大院|号|號|巷|弄|院|区|条|\\#院|\\#))?"
)

type AddressInterpreter struct {
}

func (ai AddressInterpreter) Interpret(addressText string) AddressEntity {
	visitor := NewRegionInterpreterVisitor(persister)
	return interpret(addressText, visitor)
}

func interpret(addressText string, visitor RegionInterpreterVisitor) AddressEntity {
	if len(addressText) == 0 || len(strings.TrimSpace(addressText)) <= 0 {
		return AddressEntity{}
	}
	addr := NewAddrEntity(addressText)
	extractBuildingNum(addr)

	removeSpecialChars(addr)
	brackets := extractBrackets(addr)

	extractRegion(addr, visitor)
	removeRedundancy(addr, visitor)
	extractRoad(addr)

	//addr.setText(addr.getText().replaceAll("[0-9A-Za-z\\#]+(单元|楼|室|层|米|户|\\#)", ""));
	//addr.setText(addr.getText().replaceAll("[一二三四五六七八九十]+(单元|楼|室|层|米|户)", ""));
	//if(brackets!=null && brackets.length()>0)
	//addr.setText(addr.getText()+brackets);

	return addr
}

func interprets(addrTextList []string, visitor RegionInterpreterVisitor) []AddressEntity {
	if addrTextList == nil {
		return nil
	}
	numSuccess, numFail := 0, 0
	addresses := make([]AddressEntity, 0)
	for _, addrText := range addrTextList {
		if len(addrText) == 0 {
			continue
		}
		address := interpretSimgle(addrText, visitor)
		if address.IsNil() || !address.City.IsNil() || !address.District.IsNil() {
			numFail++
			continue
		}
		numSuccess++
		addresses = append(addresses, address)
	}
	return addresses
}

func interpretSimgle(addressText string, visitor RegionInterpreterVisitor) AddressEntity {
}

// TODO point

func extractBuildingNum(addr AddressEntity) bool {
	if len(addr.Text) <= 0 {return false}

	//抽取building
	found := false
	var building string
	//xx[幢|幢|号楼|#]xx[单元]xxx
	//Matcher matcher = P_BUILDING_NUM1.matcher(addr.getText());
	//match, _ := regexp.MatchString(PBuildingNum1, addr.Text)
	RegexpPBuildingNum1.FindAllString(addr.Text, -1)


	while(matcher.find()){
		if(matcher.end()==matcher.start()) continue; //忽略null匹配结果
		building = StringUtil.substring(addr.getText(), matcher.start(), matcher.end()-1);
		//最小的匹配模式形如：7栋301，包括4个非空goup：[0:7栋301]、[1:7栋]、[2:栋]、[3:301]
		int nonEmptyGroups = 0;
		for(int i=0; i<matcher.groupCount(); i++){
			String groupStr = matcher.group(i);
			if(groupStr!=null) nonEmptyGroups++;
		}
		if(P_BUILDING_NUM_V.matcher(building).find() && nonEmptyGroups>3){
			//山东青岛市南区宁夏路118号4号楼6单元202。去掉【路xxx号】前缀
			building = StringUtil.substring(addr.getText(), matcher.start(), matcher.end()-1);
			int pos = matcher.start();
			if(building.startsWith("路") || building.startsWith("街") || building.startsWith("巷")){
				pos += building.indexOf("号")+1;
				building = StringUtil.substring(addr.getText(), pos, matcher.end()-1);
			}
			addr.setBuildingNum(building);
			addr.setText(StringUtil.head(addr.getText(), pos));
			found = true;
			break;
		}
	}
	if !found {
		//xx-xx-xx（xx栋xx单元xxx）
		matcher = P_BUILDING_NUM2.matcher(addr.getText());
		if(matcher.find()){
			addr.setBuildingNum(StringUtil.substring(addr.getText(), matcher.start(), matcher.end()-1));
			addr.setText(StringUtil.head(addr.getText(), matcher.start()));
			found = true;
		}
	}
	if !found {
		//xx组xx号
		matcher = P_BUILDING_NUM3.matcher(addr.getText());
		if(matcher.find()){
			addr.setBuildingNum(StringUtil.substring(addr.getText(), matcher.start(), matcher.end()-1));
			addr.setText(StringUtil.head(addr.getText(), matcher.start()));
			found = true;
		}
	}

	return found;
}
