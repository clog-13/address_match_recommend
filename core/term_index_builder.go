package core

import (
	. "github.com/xiiv13/address_match_recommend/models"
	"github.com/xiiv13/address_match_recommend/utils"
	"strings"
)

var (
	ignoringRegionNames = []string{
		// JD, Tmall
		"其它区", "其他地区", "其它地区", "全境", "城区", "城区以内", "城区以外", "郊区", "县城内", "内环以内", "开发区", "经济开发区", "经济技术开发区",
		// ehaier (来自TMall或HP)
		"省直辖", "省直辖市县",
		// 其他
		"地区", "市区"}
)

// TermIndexBuilder 行政区划建立倒排索引
type TermIndexBuilder struct {
	indexRoot *TermIndexEntry
	persister *AddressPersister
}

func NewTermIndexBuilder(persister *AddressPersister) *TermIndexBuilder {
	newTib := &TermIndexBuilder{
		indexRoot: NewTermIndexEntry(""),
		persister: persister,
	}
	newTib.indexRegions(persister.RootRegion().Children)
	newTib.indexIgnoring(ignoringRegionNames)
	return newTib
}

// 为行政区划建立倒排索引
func (tib *TermIndexBuilder) indexRegions(regions []*Region) {
	if regions == nil || len(regions) == 0 {
		return
	}
	for _, region := range regions {
		tii := NewTermIndexItem(convertRegionType(region), region)
		for _, name := range region.OrderedNameAndAlias() {
			tib.indexRoot.BuildIndex([]rune(name), 0, tii)
		}

		// 1. 为xx街道，建立xx镇、xx乡的别名索引项
		// 2. 为xx镇，建立xx乡的别名索引项
		// 3. 为xx乡，建立xx镇的别名索引项
		autoAlias := len(region.Name) <= 5 && len(region.Alias) == 0 &&
			(region.IsTown() || strings.HasSuffix(region.Name, "街道"))
		reName := []rune(region.Name)
		if autoAlias && len(reName) == 5 {
			switch string(reName[2]) {
			case "路":
			case "街":
			case "门":
			case "镇":
			case "村":
			case "区":
				autoAlias = false
			}
		}
		if autoAlias {
			var shortName string
			if region.IsTown() {
				shortName = utils.Head(reName, len(reName)-1)
			} else {
				shortName = utils.Head(reName, len(reName)-2)
			}
			if len(shortName) >= 2 {
				tib.indexRoot.BuildIndex([]rune(shortName), 0, tii)
			}
			if strings.HasSuffix(region.Name, "街道") || strings.HasSuffix(region.Name, "镇") {
				tib.indexRoot.BuildIndex([]rune(shortName+"乡"), 0, tii)
			}
			if strings.HasSuffix(region.Name, "街道") || strings.HasSuffix(region.Name, "乡") {
				tib.indexRoot.BuildIndex([]rune(shortName+"镇"), 0, tii)
			}
		}

		if region.Children != nil { // 递归
			tib.indexRegions(region.Children)
		}
	}
}

// 为忽略列表建立倒排索引
func (tib *TermIndexBuilder) indexIgnoring(ignoreList []string) {
	if len(ignoreList) == 0 {
		return
	}
	for _, str := range ignoreList {
		tib.indexRoot.BuildIndex([]rune(str), 0, NewTermIndexItem(IgnoreTerm, nil))
	}
}

// DeepMostQuery 深度优先匹配词条
func (tib *TermIndexBuilder) DeepMostQuery(text string, visitor *RegionInterpreterVisitor) {
	if len(text) == 0 {
		return
	}
	var pos int // 判断是否有中国开头
	if strings.HasPrefix(text, "中国") || strings.HasPrefix(text, "天朝") {
		pos += 2
	}
	tib.DeepMostPosQuery([]rune(text), pos, visitor)
}

func (tib *TermIndexBuilder) DeepMostPosQuery(text []rune, pos int, visitor *RegionInterpreterVisitor) {
	if len(text) == 0 {
		return
	}
	visitor.StartRound() // 开始匹配
	tib.deepFirstQueryRound(text, pos, tib.indexRoot.Children, visitor)
	visitor.EndRound()
}

// 获取索引对象
func (tib *TermIndexBuilder) deepFirstQueryRound(
	text []rune, pos int, entries map[rune]*TermIndexEntry, visitor *RegionInterpreterVisitor) {
	if pos > len(text)-1 {
		return
	}
	entry, ok := entries[text[pos]]
	if !ok {
		return
	}
	if entry.Children != nil && pos+1 <= len(text)-1 {
		tib.deepFirstQueryRound(text, pos+1, entry.Children, visitor)
	}
	if entry.Items != nil && len(entry.Items) > 0 {
		if visitor.Visit(entry, text, pos) {
			p := visitor.PositionAfterAcceptItem() // 给访问者一个调整当前指针的机会
			if p+1 <= len(text)-1 {
				tib.DeepMostPosQuery(text, p+1, visitor)
			}
			visitor.EndVisit(entry, p)
		}
	}
}

func convertRegionType(region *Region) int {
	switch region.Types {
	case CountryRegion:
		return CountryTerm
	case ProvinceRegion:
	case ProvinceLevelCity1:
		return ProvinceTerm
	case CityRegion:
	case ProvinceLevelCity2:
		return CityTerm
	case DistrictRegion:
	case CityLevelDistrict:
		return DistrictTerm
	case PlatformL4:
		return StreetTerm
	case TownRegion:
		return TownTerm
	case VillageRegion:
		return VillageTerm
	case StreetRegion:
		if region.IsTown() {
			return TownTerm
		} else {
			return StreetTerm
		}
	}
	return UndefinedTerm
}

func (tib *TermIndexBuilder) FullMatch(text []rune, pos int, entries map[rune]*TermIndexEntry) []*TermIndexItem {
	if entries == nil || len(text) == 0 || len(tib.indexRoot.Children) == 0 {
		return nil
	}
	entry, ok := entries[text[pos]]
	if !ok {
		return nil
	}
	if pos == len(text)-1 {
		return entry.Items
	}
	return tib.FullMatch(text, pos+1, entry.Children)
}
