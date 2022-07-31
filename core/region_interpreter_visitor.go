package core

import (
	"github.com/xiiv13/address_match_recommend/index"
	. "github.com/xiiv13/address_match_recommend/models"
	"strings"
)

// RegionInterpreterVisitor 基于倒排索引搜索匹配 省市区行政区划的访问者
type RegionInterpreterVisitor struct {
	Persister AddressPersister

	ambiguousChars map[string]struct{}
	strict         bool

	CurrentLevel, DeepMostLevel            int
	currentPos, deepMostPos                int
	fullMatchCount, deepMostFullMatchCount int

	DeepMostDivision Address
	CurDivision      Address

	stack []*index.TermIndexItem
}

func NewRegionInterpreterVisitor(ap AddressPersister, strict bool) *RegionInterpreterVisitor {
	newRiv := &RegionInterpreterVisitor{
		Persister:   ap,
		currentPos:  -1,
		deepMostPos: -1,
		strict:      strict,
	}
	newRiv.ambiguousChars["市"] = struct{}{}
	newRiv.ambiguousChars["县"] = struct{}{}
	newRiv.ambiguousChars["区"] = struct{}{}
	newRiv.ambiguousChars["镇"] = struct{}{}
	newRiv.ambiguousChars["乡"] = struct{}{}

	return newRiv
}

func (riv *RegionInterpreterVisitor) StartRound() {
	riv.CurrentLevel++
}

func (riv *RegionInterpreterVisitor) Visit(entry *index.TermIndexEntry, text string, pos int) bool {
	// 找到最匹配的 被索引对象
	acceptableItem := riv.findAcceptableItem(entry, text, pos)
	if acceptableItem == nil { // 没有匹配对象，匹配不成功，返回
		return false
	}

	// acceptableItem可能为TermType.Ignore类型,
	// 其value并不是RegionEntity对象, 因此下面region可能为nil
	region := acceptableItem.Value
	riv.stack = append(riv.stack, acceptableItem) // 更新当前状态, 匹配项压栈
	if isFullMatch(entry, region) {               // 使用全名匹配的词条数
		riv.fullMatchCount++
	}
	riv.currentPos = riv.positioning(region, entry, text, pos) // 当前结束的位置
	riv.updateCurrentDivisionState(region, entry)              // 刷新当前已经匹配上的省市区

	return true
}

func (riv *RegionInterpreterVisitor) findAcceptableItem(
	entry *index.TermIndexEntry, text string, pos int) *index.TermIndexItem {
	mostPriority := -1
	var acceptableItem *index.TermIndexItem

loop:
	for _, item := range entry.Items { // 每个 被索引对象循环，找出最匹配的
		// 仅处理省市区类型的 被索引对象，忽略其它类型的
		if !isAcceptableItemType(item.Types) {
			continue
		}

		// 省市区中的特殊名称
		if item.Types == IgnoreTerm {
			if acceptableItem == nil {
				mostPriority = 4
				acceptableItem = item
			}
			continue
		}

		region := item.Value
		// 从未匹配上任何一个省市区，则从全部被索引对象中找出一个级别最高的
		if riv.CurDivision.Province != nil {
			// 在为匹配上任务省市区情况下, 由于 `xx路` 的xx是某县区/市区/省的别名, 如江苏路, 绍兴路等等, 导致错误的匹配。
			// 如 延安路118号, 错误匹配上了延安县
			if !isFullMatch(entry, region) && pos+1 <= len(text)-1 { // 使用别名匹配，并且后面还有一个字符
				if region.Types == ProvinceRegion || region.Types == CityRegion || region.Types == CityLevelDistrict ||
					region.Types == DistrictRegion || region.Types == StreetRegion || region.Types == PlatformL4 ||
					region.Types == TownRegion { // 县区或街道

					switch string(text[pos+1]) { // 如果是某某路, 街等
					case "路":
					case "街":
					case "巷":
					case "道":
						continue loop
					}
				}
			}
			if mostPriority == -1 || int(region.Types) < mostPriority {
				mostPriority = int(region.Types)
				acceptableItem = item
			}
			continue
		}

		// 对于省市区全部匹配, 并且当前term属于非完全匹配的时候
		// 需要忽略掉当前term, 以免污染已经匹配的省市区
		if !isFullMatch(entry, region) && riv.hasThreeDivision() {
			switch region.Types {
			case ProvinceRegion:
				if region.ID != riv.CurDivision.Province.ID {
					continue loop
				}
			case CityRegion:
			case CityLevelDistrict:
				if region.ID != riv.CurDivision.City.ID {
					continue loop
				}
			case DistrictRegion:
				if region.ID != riv.CurDivision.District.ID {
					continue loop
				}
			}
		}

		// 已经匹配上部分省市区，按下面规则判断最匹配项
		// 高优先级的排除情况
		if !isFullMatch(entry, region) && pos+1 <= len(text)-1 { // 使用别名匹配，并且后面还有一个字符
			// 1. 湖南益阳沅江市万子湖乡万子湖村
			//   错误匹配方式：提取省市区时，将【万子湖村】中的字符【万子湖】匹配成【万子湖乡】，剩下一个【村】。
			// 2. 广东广州白云区均和街新市镇
			//   白云区下面有均和街道，街道、乡镇使用别名匹配时，后续字符不能是某些行政区域和道路关键字符
			if region.Types == ProvinceRegion || region.Types == CityRegion ||
				region.Types == CityLevelDistrict || region.Types == DistrictRegion ||
				region.Types == StreetRegion || region.Types == TownRegion {
				switch string(text[pos+1]) {
				case "区":
				case "县":
				case "乡":
				case "镇":
				case "村":
				case "街":
				case "路":
					continue loop
				case "大":
					if pos+2 <= len(text)-1 {
						c := string(text[pos+2])
						if c == "街" || c == "道" {
							continue loop
						}
					}
				}
			}
		}

		// 1. 匹配度最高的情况，正好是下一级行政区域
		if region.ParentID == riv.CurDivision.LeastRegion().ID {
			acceptableItem = item
			break
		}

		// 2. 中间缺一级的情况
		if mostPriority == -1 || mostPriority > 2 {
			parent := riv.Persister.GetRegion(region.ParentID)
			// 2.1 缺地级市
			if riv.CurDivision.City == nil && riv.CurDivision.Province != nil &&
				region.Types == DistrictRegion && riv.CurDivision.Province.ID == parent.ParentID {
				mostPriority = 2
				acceptableItem = item
				continue
			}
			// 2.2 缺区县
			if riv.CurDivision.District == nil && riv.CurDivision.City != nil &&
				(region.Types == StreetRegion || region.Types == TownRegion ||
					region.Types == PlatformL4 || region.Types == VillageRegion) &&
				riv.CurDivision.City.ID == parent.ParentID {
				mostPriority = 2
				acceptableItem = item
				continue
			}
		}

		// 3. 地址中省市区重复出现的情况
		if mostPriority == -1 || mostPriority > 3 {
			if riv.CurDivision.Province != nil && riv.CurDivision.Province.ID == region.ID ||
				riv.CurDivision.City != nil && riv.CurDivision.City.ID == region.ID ||
				riv.CurDivision.District != nil && riv.CurDivision.District.ID == region.ID ||
				riv.CurDivision.Street != nil && riv.CurDivision.Street.ID == region.ID ||
				riv.CurDivision.Town != nil && riv.CurDivision.Town.ID == region.ID ||
				riv.CurDivision.Village != nil && riv.CurDivision.Village.ID == region.ID {
				mostPriority = 3
				acceptableItem = item
				continue
			}
		}

		// 4. 容错
		if mostPriority == -1 || mostPriority > 4 {
			// 4.1 新疆阿克苏地区阿拉尔市
			// 到目前为止，新疆下面仍然有地级市【阿克苏地区】
			//【阿拉尔市】是县级市，以前属于地级市【阿克苏地区】，目前已变成新疆的省直辖县级行政区划
			// 即，老的行政区划关系为：新疆->阿克苏地区->阿拉尔市
			// 新的行政区划关系为：
			// 新疆->阿克苏地区
			// 新疆->阿拉尔市
			// 错误匹配方式：新疆 阿克苏地区 阿拉尔市，会导致在【阿克苏地区】下面无法匹配到【阿拉尔市】
			// 正确匹配结果：新疆 阿拉尔市
			if region.Types == CityLevelDistrict &&
				riv.CurDivision.Province != nil && riv.CurDivision.Province.ID == region.ParentID {
				mostPriority = 4
				acceptableItem = item
				continue
			}

			// 4.2 地级市-区县从属关系错误，但区县对应的省份正确，则将使用区县的地级市覆盖已匹配的地级市
			// 主要是地级市的管辖范围有调整，或者由于外部系统地级市与区县对应关系有调整导致
			if region.Types == DistrictRegion && // 必须是普通区县
				riv.CurDivision.City != nil && riv.CurDivision.Province != nil &&
				isFullMatch(entry, region) && // 使用的全名匹配
				riv.CurDivision.City.ID != region.ParentID { // 区县的地级市
				city := riv.Persister.GetRegion(region.ParentID)
				if city.ParentID == riv.CurDivision.Province.ID && !riv.hasThreeDivision() {
					mostPriority = 4
					acceptableItem = item
					continue
				}
			}
		}

		// 5. 街道、乡镇，且不符合上述情况
		if region.Types == StreetRegion || region.Types == TownRegion ||
			region.Types == VillageRegion || region.Types == PlatformL4 {
			if riv.CurDivision.District != nil {
				parent := riv.Persister.GetRegion(region.ParentID) // parent为区县
				parent = riv.Persister.GetRegion(parent.ParentID)  // parent为地级市
				if riv.CurDivision.City != nil && riv.CurDivision.City.ID == parent.ID {
					mostPriority = 5
					acceptableItem = item
					continue
				}
			} else if region.ParentID == riv.CurDivision.District.ID {
				// 已经匹配上区县
				mostPriority = 5
				acceptableItem = item
				continue
			}
		}
	}
	return acceptableItem
}

func (riv *RegionInterpreterVisitor) positioning(
	acceptedRegion *Region, entry *index.TermIndexEntry, text string, pos int) int {
	//需要调整指针的情况
	//1. 山东泰安肥城市桃园镇桃园镇山东省泰安市肥城县桃园镇东伏村
	//   错误匹配方式：提取省市区时，将【肥城县】中的字符【肥城】匹配成【肥城市】，剩下一个【县】
	if (acceptedRegion.Types == CityRegion || acceptedRegion.Types == DistrictRegion ||
		acceptedRegion.Types == StreetRegion) &&
		!isFullMatch(entry, acceptedRegion) && pos+1 <= len(text)-1 {
		c := string(text[pos+1])
		_, ok := riv.ambiguousChars[c]
		if ok { // 后面跟着特殊字符
			if acceptedRegion.Children != nil {
				for _, child := range acceptedRegion.Children {
					if string(child.Name[0]) == c {
						return pos
					}
				}
			}
			return pos + 1
		}
	}
	return pos
}

func isAcceptableItemType(types TermEnum) bool {
	switch types {
	case ProvinceTerm:
	case CityTerm:
	case DistrictTerm:
	case StreetTerm:
	case TownTerm:
	case VillageTerm:
	case IgnoreTerm:
		return true
	}
	return false
}

func isFullMatch(entry *index.TermIndexEntry, region *Region) bool {
	if len(entry.Key) == len(region.Name) {
		return true
	}
	if region.Types == StreetRegion && strings.HasSuffix(region.Name, "街道") &&
		len(region.Name) == len(entry.Key)+1 {
		return true
	}
	return false
}

func (riv *RegionInterpreterVisitor) hasThreeDivision() bool {
	return riv.CurDivision.City.ParentID == riv.CurDivision.Province.ID &&
		riv.CurDivision.District.ParentID == riv.CurDivision.City.ID

}

// 更新当前已匹配区域对象的状态
func (riv *RegionInterpreterVisitor) updateCurrentDivisionState(
	region *Region, entry *index.TermIndexEntry) {
	// region为重复项，无需更新状态
	if region.Equal(riv.CurDivision.Province) || region.Equal(riv.CurDivision.City) ||
		region.Equal(riv.CurDivision.District) || region.Equal(riv.CurDivision.Street) ||
		region.Equal(riv.CurDivision.Town) || region.Equal(riv.CurDivision.Village) {
		return
	}

	// 非严格模式 || 只有一个父项
	needUpdateCityAndProvince := !riv.strict || len(entry.Items) == 1
	switch region.Types {
	case ProvinceRegion:
	case ProvinceLevelCity1:
		riv.CurDivision.Province = region
		riv.CurDivision.City = nil
	case CityRegion:
	case ProvinceLevelCity2:
		riv.CurDivision.City = region
		if riv.CurDivision.Province != nil {
			riv.CurDivision.Province = riv.Persister.GetRegion(region.ParentID)
		}
	case CityLevelDistrict:
		riv.CurDivision.City = region
		riv.CurDivision.District = region
		if riv.CurDivision.Province != nil {
			riv.CurDivision.Province = riv.Persister.GetRegion(region.ParentID)
		}
	case DistrictRegion:
		riv.CurDivision.District = region
		// 成功匹配了区县，则强制更新地级市
		riv.CurDivision.City = riv.Persister.GetRegion(riv.CurDivision.District.ParentID)
		if riv.CurDivision.Province == nil {
			riv.CurDivision.Province = riv.Persister.GetRegion(riv.CurDivision.City.ParentID)
		}
	case StreetRegion:
	case PlatformL4:
		if riv.CurDivision.Street == nil {
			riv.CurDivision.Street = region
		}
		if riv.CurDivision.District == nil {
			riv.CurDivision.District = riv.Persister.GetRegion(region.ParentID)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District.ParentID)
		}
	case TownRegion:
		if riv.CurDivision.Town == nil {
			riv.CurDivision.Town = region
		}
		if riv.CurDivision.District == nil {
			riv.CurDivision.District = riv.Persister.GetRegion(region.ParentID)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District.ParentID)
		}
	case VillageRegion:
		if riv.CurDivision.Village == nil {
			riv.CurDivision.Village = region
		}
		if riv.CurDivision.District == nil {
			riv.CurDivision.District = riv.Persister.GetRegion(region.ParentID)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District.ParentID)
		}
	}
}

// TODO

func (riv *RegionInterpreterVisitor) updateCityAndProvince(parentId uint) {
	if riv.CurDivision.City == nil {
		riv.CurDivision.City = riv.Persister.GetRegion(parentId)
		if riv.CurDivision.Province == nil {
			riv.CurDivision.Province = riv.Persister.GetRegion(riv.CurDivision.City.ParentID)
		}
	}
}

func (riv *RegionInterpreterVisitor) PositionAfterAcceptItem() int {
	return riv.currentPos
}

func (riv *RegionInterpreterVisitor) EndVisit(entry *index.TermIndexEntry, pos int) {
	riv.checkDeepMost()

	indexTerm := riv.stack[len(riv.stack)-1] // 当前访问的索引对象出栈
	riv.stack = riv.stack[:len(riv.stack)-1]
	riv.currentPos = pos - len(entry.Key) // 恢复当前位置指针

	if isFullMatch(entry, indexTerm.Value) {
		riv.fullMatchCount++ // 更新全名匹配的数量
	}
	if indexTerm.Types == IgnoreTerm { // 如果是忽略项，无需更新当前已匹配的省市区状态
		return
	}

	// 扫描一遍stack，找出街道street、乡镇town、村庄village，以及省市区中级别最低的一个least
	var street, town, village, least *Region
loop:
	for _, v := range riv.stack {
		if v.Types == IgnoreTerm {
			continue
		}
		r := v.Value
		switch r.Types {
		case StreetRegion:
		case PlatformL4:
			street = r
			continue loop
		case TownRegion:
			town = r
			continue loop
		case VillageRegion:
			village = r
			continue loop
		}
		if least == nil {
			least = r
			continue loop
		}
		if r.Types > least.Types {
			least = r
		}
	}

	if street == nil {
		riv.CurDivision.Street = nil // 剩余匹配项中没有街道了
	}
	if town == nil {
		riv.CurDivision.Town = nil // 剩余匹配项中没有乡镇了
	}
	if village == nil {
		riv.CurDivision.Village = nil // 剩余匹配项中没有村庄了
	}

	// 只有街道、乡镇、村庄都没有时，才开始清空省市区
	if riv.CurDivision.Street != nil || riv.CurDivision.Town != nil || riv.CurDivision.Village != nil {
		return
	}

	if least != nil {
		switch least.Types {
		case ProvinceRegion:
		case ProvinceLevelCity1:
			riv.CurDivision.City = nil
			riv.CurDivision.District = nil
			return
		case CityRegion:
		case ProvinceLevelCity2:
			riv.CurDivision.District = nil
			return
		default:
			return
		}
	}

	// least为null，说明stack中什么都不剩了
	riv.CurDivision.Province = nil
	riv.CurDivision.City = nil
	riv.CurDivision.District = nil
}

func (riv *RegionInterpreterVisitor) EndRound() {
	riv.checkDeepMost()
	riv.CurrentLevel--
}

func (riv *RegionInterpreterVisitor) checkDeepMost() {
	if len(riv.stack) > riv.DeepMostLevel {
		riv.DeepMostLevel = len(riv.stack)
		riv.deepMostPos = riv.currentPos
		riv.deepMostFullMatchCount = riv.fullMatchCount
		riv.DeepMostDivision.Province = riv.CurDivision.Province
		riv.DeepMostDivision.City = riv.CurDivision.City
		riv.DeepMostDivision.District = riv.CurDivision.District
		riv.DeepMostDivision.Street = riv.CurDivision.Street
		riv.DeepMostDivision.Town = riv.CurDivision.Town
		riv.DeepMostDivision.Village = riv.CurDivision.Village
	}
}

func (riv *RegionInterpreterVisitor) HasResult() bool {
	return riv.deepMostPos > 0 && riv.DeepMostDivision.District != nil
}

func (riv *RegionInterpreterVisitor) GetDevision() Address {
	return riv.DeepMostDivision
}

func (riv *RegionInterpreterVisitor) MatchCount() int {
	return riv.DeepMostLevel
}

func (riv *RegionInterpreterVisitor) FullMatchCount() int {
	return riv.deepMostFullMatchCount
}

func (riv *RegionInterpreterVisitor) EndPosition() int {
	return riv.deepMostPos
}

func (riv *RegionInterpreterVisitor) Reset() {
	riv.CurrentLevel = 0
	riv.DeepMostLevel = 0
	riv.currentPos = -1
	riv.deepMostPos = -1
	riv.fullMatchCount = 0
	riv.deepMostFullMatchCount = 0
	riv.DeepMostDivision.Province = nil
	riv.DeepMostDivision.City = nil
	riv.DeepMostDivision.District = nil
	riv.DeepMostDivision.Street = nil
	riv.DeepMostDivision.Town = nil
	riv.DeepMostDivision.Village = nil
	riv.CurDivision.Province = nil
	riv.CurDivision.City = nil
	riv.CurDivision.District = nil
	riv.CurDivision.Street = nil
	riv.CurDivision.Town = nil
	riv.CurDivision.Village = nil
}
