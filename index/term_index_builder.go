package index

import (
	"github.com/xiiv13/address_match_recommend/core"
	. "github.com/xiiv13/address_match_recommend/models"
	"github.com/xiiv13/address_match_recommend/utils"
	"strings"
)

// TermIndexBuilder 行政区划建立倒排索引
type TermIndexBuilder struct {
	indexRoot TermIndexEntry
}

func NewTermIndexBuilder(persister core.AddressPersister, ingoringRegionNames []string) TermIndexBuilder {
	newTib := TermIndexBuilder{}
	newTib.indexRegions(persister.GetRootRegionChilden())
	newTib.indexIgnoring(ingoringRegionNames)
	return newTib
}

// 为行政区划建立倒排索引
func (tib *TermIndexBuilder) indexRegions(regions *[]*Region) {
	if len(*regions) == 0 {
		return
	}
	for _, region := range *regions {
		tii := NewTermIndexItem(convertRegionType(region), region)
		for _, name := range region.OrderedNameAndAlias() {
			tib.indexRoot.BuildIndex(name, 0, tii)
		}

		// 1. 为xx街道，建立xx镇、xx乡的别名索引项
		// 2. 为xx镇，建立xx乡的别名索引项
		// 3. 为xx乡，建立xx镇的别名索引项
		autoAlias := len(region.Name) <= 5 && len(region.Alias) == 0 &&
			(region.IsTown() || strings.HasSuffix(region.Name, "街道"))
		if autoAlias && len(region.Name) == 5 {
			switch string(region.Name[2]) {
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
				shortName = utils.Head(region.Name, len(region.Name)-1)
			} else {
				shortName = utils.Head(region.Name, len(region.Name)-2)
			}
			if len(shortName) >= 2 {
				tib.indexRoot.BuildIndex(shortName, 0, tii)
			}
			if strings.HasSuffix(region.Name, "街道") || strings.HasSuffix(region.Name, "镇") {
				tib.indexRoot.BuildIndex(shortName+"乡", 0, tii)
			}
			if strings.HasSuffix(region.Name, "街道") || strings.HasSuffix(region.Name, "乡") {
				tib.indexRoot.BuildIndex(shortName+"镇", 0, tii)
			}
		}
		// 递归
		if region.Children != nil {
			tib.indexRegions(&region.Children)
		}
	}
}

// 为忽略列表建立倒排索引
func (tib *TermIndexBuilder) indexIgnoring(ignoreList []string) {
	if len(ignoreList) == 0 {
		return
	}
	for _, v := range ignoreList {
		tib.indexRoot.BuildIndex(v, 0, NewTermIndexItem(IgnoreTerm, nil))
	}
}

// DeepMostQuery 深度优先匹配词条
func (tib *TermIndexBuilder) DeepMostQuery(text string, visitor TermIndexVisitor) {
	if len(text) == 0 {
		return
	}
	var pos int
	if strings.HasPrefix(text, "中国") || strings.HasPrefix(text, "天朝") {
		pos += 2
	}
	tib.DeepPosMostQuery(text, pos, visitor)
}

func (tib *TermIndexBuilder) DeepPosMostQuery(text string, pos int, visitor TermIndexVisitor) {
	if len(text) == 0 {
		return
	}
	// 开始匹配
	visitor.StartRound()
	tib.deepFirstQueryRound(text, pos, visitor)
	visitor.EndRound()
}

func (tib *TermIndexBuilder) deepFirstQueryRound(text string, pos int, visitor TermIndexVisitor) {
	// 获取索引对象
	if pos > len(text)-1 {
		return
	}
	entry, ok := tib.indexRoot.Children[text[pos]]
	if !ok {
		return
	}
	if pos+1 <= len(text)-1 {
		tib.deepFirstQueryRound(text, pos+1, visitor)
	}
	if len(entry.Items) > 0 {
		if visitor.Visit(entry, text, pos) {
			p := visitor.PositionAfterAcceptItem()
			if p+1 <= len(text)-1 {
				tib.DeepPosMostQuery(text, p+1, visitor)
			}
			visitor.EndVisit(entry, p)
		}
	}
}

func convertRegionType(region *Region) TermEnum {
	switch region.Types {
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
	default:
		return UndefinedTerm
	}
	return UndefinedTerm
}
