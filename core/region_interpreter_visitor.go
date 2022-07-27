package core

import (
	"address_match_recommend/index"
	. "address_match_recommend/models"
	"strings"
)

type RegionInterpreterVisitor struct {
	IsDebug        bool
	AmbiguousChars map[string]struct{}
	Persister      AddressPersister
	strict         bool

	CurrentLevel, DeepMostLevel            int
	CurrentPos, DeepMostPos                int
	FullMatchCount, DeepMostFullMatchCount int

	DeepMostDivision Division
	CurDivision      Division

	stack []index.TermIndexItem
}

func NewRegionInterpreterVisitor(ap AddressPersister) RegionInterpreterVisitor {
	newRiv := RegionInterpreterVisitor{
		CurrentPos:  -1,
		DeepMostPos: -1,
		Persister:   ap,
	}
	newRiv.AmbiguousChars["市"] = struct{}{}
	newRiv.AmbiguousChars["县"] = struct{}{}
	newRiv.AmbiguousChars["区"] = struct{}{}
	newRiv.AmbiguousChars["镇"] = struct{}{}
	newRiv.AmbiguousChars["乡"] = struct{}{}

	return newRiv
}

func (riv RegionInterpreterVisitor) visit(entry index.TermIndexEntry, text string, pos int) bool {
	// 找到最匹配的 被索引对象
	acceptableItem := riv.findAcceptableItem(entry, text, pos)
	if acceptableItem.IsNil() { // 没有匹配对象，匹配不成功，返回
		return false
	}

	// acceptableItem可能为TermType.Ignore类型，此时其value并不是RegionEntity对象，因此下面region的值可能为null
	region := acceptableItem.Value.(RegionEntity)
	riv.stack = append(riv.stack, acceptableItem)
	if isFullMatch(entry, region) { // 使用全名匹配的词条数
		riv.FullMatchCount++
	}
	riv.CurrentPos = riv.positioning(region, entry, text, pos) // 当前结束的位置
	riv.updateCurrentDivisionState(region, entry)              // 刷新当前已经匹配上的省市区

	return true
}

func (riv RegionInterpreterVisitor) findAcceptableItem(
	entry index.TermIndexEntry, text string, pos int) index.TermIndexItem {
	mostPriority := -1
	var acceptableItem index.TermIndexItem

loop:
	for _, item := range entry.Items { // 每个 被索引对象循环，找出最匹配的
		// 仅处理省市区类型的 被索引对象，忽略其它类型的
		if !isAcceptableItemType(item.Types) {
			continue
		}

		// 省市区中的特殊名称
		if item.Types == IgnoreTerm {
			if acceptableItem.IsNil() {
				mostPriority = 4
				acceptableItem = item
			}
			continue
		}

		region := item.Value.(RegionEntity)
		// 从未匹配上任何一个省市区，则从全部被索引对象中找出一个级别最高的
		if !riv.CurDivision.Province.IsNil() {
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
			if mostPriority == -1 || region.Types < mostPriority {
				mostPriority = region.Types
				acceptableItem = item
			}
			continue
		}

		// 对于省市区全部匹配, 并且当前term属于非完全匹配的时候
		// 需要忽略掉当前term, 以免污染已经匹配的省市区
		if !isFullMatch(entry, region) && riv.hasThreeDivision() {
			switch region.Types {
			case ProvinceRegion:
				if region.Id != riv.CurDivision.Province.Id {
					continue loop
				}
			case CityRegion:
			case CityLevelDistrict:
				if region.Id != riv.CurDivision.City.Id {
					continue loop
				}
			case DistrictRegion:
				if region.Id != riv.CurDivision.District.Id {
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
		if region.ParentId == riv.CurDivision.LeastRegion().Id {
			acceptableItem = item
			break
		}

		// 2. 中间缺一级的情况
		if mostPriority == -1 || mostPriority > 2 {
			parent := persister.GetRegion(region.ParentId)
			// 2.1 缺地级市
			if riv.CurDivision.City.IsNil() && !riv.CurDivision.Province.IsNil() &&
				region.Types == DistrictRegion && riv.CurDivision.Province.Id == parent.ParentId {
				mostPriority = 2
				acceptableItem = item
				continue
			}
			// 2.2 缺区县
			if riv.CurDivision.District.IsNil() && !riv.CurDivision.City.IsNil() &&
				(region.Types == StreetRegion || region.Types == TownRegion ||
					region.Types == PlatformL4 || region.Types == VillageRegion) &&
				riv.CurDivision.City.Id == parent.ParentId {
				mostPriority = 2
				acceptableItem = item
				continue
			}
		}

		// 3. 地址中省市区重复出现的情况
		if mostPriority == -1 || mostPriority > 3 {
			if !riv.CurDivision.Province.IsNil() && riv.CurDivision.Province.Id == region.Id ||
				!riv.CurDivision.City.IsNil() && riv.CurDivision.City.Id == region.Id ||
				!riv.CurDivision.District.IsNil() && riv.CurDivision.District.Id == region.Id ||
				!riv.CurDivision.Street.IsNil() && riv.CurDivision.Street.Id == region.Id ||
				!riv.CurDivision.Town.IsNil() && riv.CurDivision.Town.Id == region.Id ||
				!riv.CurDivision.Village.IsNil() && riv.CurDivision.Village.Id == region.Id {
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
				!riv.CurDivision.Province.IsNil() && riv.CurDivision.Province.Id == region.ParentId {
				mostPriority = 4
				acceptableItem = item
				continue
			}

			// 4.2 地级市-区县从属关系错误，但区县对应的省份正确，则将使用区县的地级市覆盖已匹配的地级市
			// 主要是地级市的管辖范围有调整，或者由于外部系统地级市与区县对应关系有调整导致
			if region.Types == DistrictRegion && // 必须是普通区县
				!riv.CurDivision.City.IsNil() && !riv.CurDivision.Province.IsNil() &&
				isFullMatch(entry, region) && // 使用的全名匹配
				riv.CurDivision.City.Id != region.ParentId { // 区县的地级市
				city := persister.GetRegion(region.ParentId)
				if city.ParentId == riv.CurDivision.Province.Id && !riv.hasThreeDivision() {
					mostPriority = 4
					acceptableItem = item
					continue
				}
			}
		}

		// 5. 街道、乡镇，且不符合上述情况
		if region.Types == StreetRegion || region.Types == TownRegion ||
			region.Types == VillageRegion || region.Types == PlatformL4 {
			if !riv.CurDivision.District.IsNil() {
				parent := persister.GetRegion(region.ParentId) // parent为区县
				parent = persister.GetRegion(parent.ParentId)  // parent为地级市
				if !riv.CurDivision.City.IsNil() && riv.CurDivision.City.Id == parent.Id {
					mostPriority = 5
					acceptableItem = item
					continue
				}
			} else if region.ParentId == riv.CurDivision.District.Id {
				// 已经匹配上区县
				mostPriority = 5
				acceptableItem = item
				continue
			}
		}
	}
	return acceptableItem

}

func (riv RegionInterpreterVisitor) positioning(
	acceptedRegion RegionEntity, entry index.TermIndexEntry, text string, pos int) int {
	if acceptedRegion.IsNil() {
		return pos
	}
	//需要调整指针的情况
	//1. 山东泰安肥城市桃园镇桃园镇山东省泰安市肥城县桃园镇东伏村
	//   错误匹配方式：提取省市区时，将【肥城县】中的字符【肥城】匹配成【肥城市】，剩下一个【县】
	if (acceptedRegion.Types == CityRegion || acceptedRegion.Types == DistrictRegion ||
		acceptedRegion.Types == StreetRegion) &&
		!isFullMatch(entry, acceptedRegion) && pos+1 <= len(text)-1 {
		c := string(text[pos+1])
		_, ok := riv.AmbiguousChars[c]
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

func isFullMatch(entry index.TermIndexEntry, region RegionEntity) bool {
	if region.IsNil() {
		return false
	}
	if len(entry.Key) == len(region.Name) {
		return true
	}
	if region.Types == StreetRegion && strings.HasSuffix(region.Name, "街道") &&
		len(region.Name) == len(entry.Key)+1 {
		return true
	}
	return false
}

func (riv RegionInterpreterVisitor) hasThreeDivision() bool {
	return riv.CurDivision.City.ParentId == riv.CurDivision.Province.Id &&
		riv.CurDivision.District.ParentId == riv.CurDivision.City.Id

}

// 更新当前已匹配区域对象的状态
func (riv RegionInterpreterVisitor) updateCurrentDivisionState(
	region RegionEntity, entry index.TermIndexEntry) {
	if region.IsNil() {
		return
	}
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
		riv.CurDivision.City = RegionEntity{}
	case CityRegion:
	case ProvinceLevelCity2:
		riv.CurDivision.City = region
		if !riv.CurDivision.Province.IsNil() {
			riv.CurDivision.Province = persister.GetRegion(region.ParentId)
		}
	case CityLevelDistrict:
		riv.CurDivision.City = region
		riv.CurDivision.District = region
		if !riv.CurDivision.Province.IsNil() {
			riv.CurDivision.Province = persister.GetRegion(region.ParentId)
		}
	case DistrictRegion:
		riv.CurDivision.District = region
		// 成功匹配了区县，则强制更新地级市
		riv.CurDivision.City = persister.GetRegion(riv.CurDivision.District.ParentId)
		if riv.CurDivision.Province.IsNil() {
			riv.CurDivision.Province = persister.GetRegion(riv.CurDivision.City.ParentId)
		}
	case StreetRegion:
	case PlatformL4:
		if riv.CurDivision.Street.IsNil() {
			riv.CurDivision.Street = region
		}
		if riv.CurDivision.District.IsNil() {
			riv.CurDivision.District = persister.GetRegion(region.ParentId)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District)
		}
	case TownRegion:
		if riv.CurDivision.Town.IsNil() {
			riv.CurDivision.Town = region
		}
		if riv.CurDivision.District.IsNil() {
			riv.CurDivision.District = persister.GetRegion(region.ParentId)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District)
		}
	case VillageRegion:
		if riv.CurDivision.Village.IsNil() {
			riv.CurDivision.Village = region
		}
		if riv.CurDivision.District.IsNil() {
			riv.CurDivision.District = persister.GetRegion(region.ParentId)
		}
		if needUpdateCityAndProvince {
			riv.updateCityAndProvince(riv.CurDivision.District)
		}
	}
}

// TODO

func (riv RegionInterpreterVisitor) updateCityAndProvince(distinct RegionEntity) {
	if distinct.IsNil() {
		return
	}
	if riv.CurDivision.City.IsNil() {
		riv.CurDivision.City = persister.GetRegion(distinct.ParentId)
		if riv.CurDivision.Province.IsNil() {
			riv.CurDivision.Province = persister.GetRegion(riv.CurDivision.City.ParentId)
		}
	}
}
