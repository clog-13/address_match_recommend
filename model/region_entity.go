package model

import (
	"strings"
)

const (
	Undefinded = 0
	Country    = 10
	Street     = 450
)

/**
 * 行政区域实体。标准行政区域说明：
 *
 * <ul style="color:red;">
 * <li>直辖市：采用【北京 -&gt; 北京市 -&gt; 下属区县】、【天津 -&gt; 天津市 -&gt; 下属区县】形式表示；</li>
 * <li>省直辖县级行政区划：例如【湖北省 -&gt; 潜江市】，parent_id为湖北省，其下没有区县数据。
 * 在匹配地址时需注意，有的地址库会采用【湖北省 -&gt; 潜江 -&gt; 潜江市】方式表示，不做特殊处理将无法匹配上。</li>
 * <li>街道乡镇：所有的街道乡镇都使用{@link RegionType#Street}存储，父级ID为区县，包括街道、乡镇，以及各种特殊的街道一级行政区域。</li>
 * <li>附加乡镇：不在标准行政区域体系中，由历史地址数据中通过文本匹配出来的乡镇，都使用{@link RegionType#Town}存储，父级ID为区县。</li>
 * <li>附加村庄：不在标准行政区域体系中，由历史地址数据中通过文本匹配出来的村庄，都使用{@link RegionType#Town}存储，父级ID为区县。</li>
 * <li>平台相关的特殊区域划分：主要纳入了京东的特殊4级地址，例如【三环内】，都使用{@link RegionType#PlatformL4}存储，父级ID为区县。</li>
 * </ul>
 *
 * <p>
 * Table: <strong>bas_region</strong></p>
 * <p>
 * <table class="er-mapping" cellspacing=0 cellpadding=0 style="border:solid 1 #666;padding:3px;">
 *   <tr style="background-color:#ddd;Text-align:Left;">
 *     <th nowrap>属性名</th><th nowrap>属性类型</th><th nowrap>字段名</th><th nowrap>字段类型</th><th nowrap>说明</th>
 *   </tr>
 *   <tr><td>id</td><td>{@link Integer}</td><td>id</td><td>int</td><td>&nbsp;</td></tr>
 *   <tr><td>parentId</td><td>{@link Integer}</td><td>parent_id</td><td>int</td><td>&nbsp;</td></tr>
 *   <tr><td>name</td><td>{@link String}</td><td>name</td><td>varchar</td><td>&nbsp;</td></tr>
 *   <tr><td>type</td><td>{@link Integer}</td><td>type</td><td>int</td><td>&nbsp;</td></tr>
 *   <tr><td>zip</td><td>{@link String}</td><td>zip</td><td>varchar</td><td>&nbsp;</td></tr>
 * </table></p>
 *
 * @author Richie 刘志斌 yudi@sina.com
 * @since 2016/9/4 1:22:41
 */

// RegionEntity 行政区域实体
type RegionEntity struct {
	serialVersionUID int64 // -111163973997033386L

	id           int64
	parentId     int64
	name         string
	alias        string
	types        int // RegionType enum
	zip          string
	children     []RegionEntity
	orderedNames []string
}

func (r RegionEntity) IsTown() bool {
	switch r.types {
	case Country:
		return true
	case Street:
		if r.name == "" {
			return false
		}
		return len(r.name) <= 4 &&
			(string(r.name[len(r.name)-1]) == "镇" || string(r.name[len(r.name)-1]) == "乡")
	}
	return false
}

// OrderedNameAndAlias 获取所有名称和别名列表，按字符长度倒排序。
func (r RegionEntity) OrderedNameAndAlias() []string {
	if r.orderedNames == nil {
		return r.orderedNames
	}
	r.buildOrderedNameAndAlias()
	return r.orderedNames
}

func (r RegionEntity) buildOrderedNameAndAlias() {
	if r.orderedNames != nil {
		return
	}
	tokens := make([]string, 0)
	if r.alias != "" && len(strings.TrimSpace(r.alias)) > 0 {
		tokens = strings.Split(strings.TrimSpace(r.alias), ";")
	}
	if tokens == nil || len(tokens) <= 0 {
		r.orderedNames = make([]string, 1)
	} else {
		r.orderedNames = make([]string, len(tokens)+1)
	}
	r.orderedNames = append(r.orderedNames, r.name)
	if tokens != nil {
		for _, v := range tokens {
			if v == "" || len(strings.TrimSpace(v)) <= 0 {
				continue
			}
			r.orderedNames = append(r.orderedNames, strings.TrimSpace(v))
		}
	}

	exchanged := true
	endIndex := len(r.orderedNames) - 1
	for exchanged && endIndex > 0 {
		exchanged = false
		for i := 0; i < endIndex; i++ {
			if len(r.orderedNames[i]) < len(r.orderedNames[i+1]) {
				temp := r.orderedNames[i]
				r.orderedNames[i] = r.orderedNames[i+1]
				r.orderedNames[i+1] = temp
				exchanged = true
			}
		}
		endIndex--
	}
}
