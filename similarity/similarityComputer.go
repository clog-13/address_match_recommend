package similarity

import (
	"address_match_recommend/common"
	"address_match_recommend/enum"
	. "address_match_recommend/model"
	"math"
	"strconv"
)

// TODO point

var (
	BoostM  = 1.0  // 正常权重
	BoostL  = 2.0  // 加权高值
	BoostXl = 4.0  // 加权高值
	BoostS  = 0.5  // 降权
	BoostXs = 0.25 // 降权

	MissingIdf = 4.0

	interpreter   AddressInterpreter
	segmenter     = new(SimpleSegmenter)
	defaultTokens = make([]string, 0)
	cacheFolder   string

	CacheVectorsInMemory = false
	VECTORS_CACHE        = make(map[string][]Document)
	IdfCache             = make(map[string]map[string]float64)

	timeBoost int64
)

// 分词，设置词条权重
func analyse(addr AddressEntity) Document {
	doc := NewDocument(addr.Id)

	//1. 分词。仅针对AddressEntity的text（地址解析后剩余文本）进行分词。
	tokens := make([]string, 0)
	if len(addr.Text) > 0 {
		tokens = segmenter.Segment(addr.Text)
	}
	terms := make([]Term, len(tokens)+4)

	//2. 生成term
	if !addr.Town.IsNil() {
		doc.Town = NewTerm(enum.Town, addr.Town.Name)
		terms = append(terms, doc.Town)
	}
	if !addr.Village.IsNil() {
		doc.Village = NewTerm(enum.Village, addr.Village.Name)
		terms = append(terms, doc.Village)
	}
	if len(addr.Road) > 0 {
		doc.Road = NewTerm(enum.Road, addr.Road)
		terms = append(terms, doc.Road)
	}
	if len(addr.RoadNum) > 0 {
		roadNumTerm := NewTerm(enum.RoadNum, addr.RoadNum)
		doc.RoadNum = roadNumTerm
		doc.RoadNumValue = translateRoadNum(addr.RoadNum)
		roadNumTerm.Ref = &doc.Road
		terms = append(terms, doc.RoadNum)
	}

	//2.2 地址文本分词后的token
	for _, v := range tokens {
		addTerm(v, enum.Text, terms)
	}

	idfs, ok := IdfCache[buildCacheKey(addr)]
	if ok {
		idf := float64(-1)
		for _, v := range terms {
			idf = idfs[generateIDFCacheEntryKey(v)]
			if idf == float64(-1) {
				v.Idf = MissingIdf
			} else {
				v.Idf = idf
			}
		}
	}
	doc.Terms = terms
	return doc
}

/**
 * 为所有文档的全部词条统计逆向引用情况。
 * @param docs 所有文档。
 * @return 全部词条的逆向引用情况，key为词条，value为该词条在多少个文档中出现过。
 */
func statInverseDocRefers(docs []Document) map[string]int {
	idrc := make(map[string]int)
	if docs == nil {
		return idrc
	}
	var key string
	for _, doc := range docs {
		if doc.Terms == nil {
			continue
		}
		for _, term := range doc.Terms {
			key = generateIDFCacheEntryKey(term)
			_, ok := idrc[key]
			if ok {
				idrc[key] = idrc[key] + 1
			} else {
				idrc[key] = 1
			}

		}
	}
	return idrc
}

func generateIDFCacheEntryKey(term Term) string {
	key := term.Text
	if enum.RoadNum == term.Types {
		num := translateRoadNum(key)
		if term.Ref == nil {
			key = ""
		} else {
			key = term.Ref.Text
		}
		key += "-" + strconv.Itoa(num)
	}
	return key
}

/**
 * 计算词条加权权重boost值。
 * @param forDoc true:为地址库文档词条计算boost；false:为查询文档词条计算boost。
 * @param qdoc 查询文档。
 * @param qterm 查询文档词条。
 * @param ddoc 地址库文档。
 * @param dterm 地址库文档词条。
 * @return
 */
func getBoostValue(forDoc bool, qdoc Document, qterm Term, ddoc Document, dterm Term) float64 {
	value := BoostM
	var types byte
	if forDoc {
		types = dterm.Types
	} else {
		types = qterm.Types
	}
	switch types {
	case enum.Province:
	case enum.City:
	case enum.District:
		value = BoostXl // 省市区、道路出现频次高，IDF值较低，但重要程度最高，因此给予比较高的加权权重
	case enum.StreetTerm:
		value = BoostXs //一般人对于城市街道范围概念不强，在地址中随意选择街道的可能性较高，因此降权处理
	case enum.Text:
		value = BoostM
	case enum.Town:
	case enum.Village:
		value = BoostXs
		if enum.Town == types { //乡镇
			//查询文档和地址库文档都有乡镇，为乡镇加权。注意：存在乡镇相同、不同两种情况。
			//  乡镇相同：查询文档和地址库文档都加权BOOST_L，提高相似度
			//  乡镇不同：只有查询文档的词条加权BOOST_L，地址库文档的词条因无法匹配不会进入该函数。结果是拉开相似度的差异
			if !qdoc.Town.IsNil() && !ddoc.Town.IsNil() {
				value = BoostL
			}
		} else { //村庄
			//查询文档和地址库文档都有乡镇且乡镇相同，且查询文档和地址库文档都有村庄时，为村庄加权
			//与上述乡镇类似，存在村庄相同和不同两种情况
			if !qdoc.Village.IsNil() && !ddoc.Village.IsNil() && !qdoc.Town.IsNil() {
				if qdoc.Town == ddoc.Town {
					if qdoc.Village == ddoc.Village {
						value = BoostXl
					} else {
						value = BoostL
					}
				} else if !ddoc.Town.IsNil() {
					if !forDoc {
						value = BoostL
					} else {
						value = BoostS
					}
				}
			}
		}
	case enum.Road:
	case enum.RoadNum:
		if qdoc.Town.IsNil() || qdoc.Village.IsNil() { // 有乡镇有村庄，不再考虑道路、门牌号的加权
			if enum.Road == types { //道路
				if !qdoc.Road.IsNil() && !ddoc.Road.IsNil() {
					value = BoostL
				}
			} else { // 门牌号。注意：查询文档和地址库文档的门牌号都会进入此处执行，这一点跟Road、Town、Village不同。
				if qdoc.RoadNumValue > 0 && ddoc.RoadNumValue > 0 && !qdoc.Road.IsNil() && qdoc.Road == ddoc.Road {
					if qdoc.RoadNumValue == ddoc.RoadNumValue {
						value = float64(3)
					} else {
						if forDoc {
							value = (1 / math.Sqrt(math.Sqrt(math.Abs(float64(qdoc.RoadNumValue-ddoc.RoadNumValue))+1))) * BoostL
						} else {
							value = float64(3)
						}
					}
				}
			}
		}
	}
	return value
}

/**
 * 将道路门牌号中的数字提取出来
 * @param text 道路门牌号，例如40号院、甲一号院等
 * @return 返回门牌号数字
 */
func translateRoadNum(text string) int {
	if len(text) == 0 {
		return 0
	}
	var sb string
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= '0' && c <= '9' { // ANSI数字字符
			sb += string(c)
			continue
		}
		// TODO
		switch string(c) { // 中文全角数字字符
		case "０":
			sb += "0"
		case "１":
			sb += "1"
		case "２":
			sb += "2"
		case "３":
			sb += "3"
		case "４":
			sb += "4"
		case "５":
			sb += "5"
		case "６":
			sb += "6"
		case "７":
			sb += "7"
		case "８":
			sb += "8"
		case "９":
			sb += "9"
		}
	}

	if len(sb) > 0 {
		ri, _ := strconv.Atoi(sb)
		return ri
	}
	isTen := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		if isTen {
			pre := len(sb) > 0
			sc := string(c)
			post := sc == "一" || sc == "二" || sc == "三" || sc == "四" || sc == "五" || sc == "六" || sc == "七" || sc == "八" || sc == "九"
			if pre {
				if !post {
					sb += "0"
				}
			} else {
				if post {
					sb += "1"
				} else {
					sb += "10"
				}
			}
			isTen = false
		}

		switch string(c) {
		case "一":
			sb += "1"
		case "二":
			sb += "2"
		case "三":
			sb += "3"
		case "四":
			sb += "4"
		case "五":
			sb += "5"
		case "六":
			sb += "6"
		case "七":
			sb += "7"
		case "八":
			sb += "8"
		case "九":
			sb += "9"
		case "十":
			isTen = true
		}
		if len(sb) > 0 {
			break
		}
	}

	if isTen {
		if len(sb) > 0 {
			sb += "0"
		} else {
			sb += "10"
		}
	}
	if len(sb) > 0 {
		rs, _ := strconv.Atoi(sb)
		return rs
	}
	return 0
}

func addTerm(text string, types byte, terms []Term) Term {
	if len(text) == 0 {
		return Term{}
	}
	termText := text
	for _, v := range terms {
		if v.Text == termText {
			return v
		}
	}
	newTerm := NewTerm(types, termText)
	terms = append(terms, newTerm)
	return newTerm
}

func FindsimilarAddress(addressText string, topN int, explain bool) *Query {
	query := NewQuery(topN)

	// TODO
	queryAddr := FormatAddressEntity(addressText)

	//从文件缓存或内存缓存获取所有文档。
	allDocs := loadDocunentsFromCache(queryAddr)

	//为词条计算特征值
	queryDoc := analyse(queryAddr)
	query.QueryAddr = queryAddr
	query.QueryDoc = queryDoc

	//对应地址库中每条地址计算相似度，并保留相似度最高的topN条地址
	var similarity float64
	for _, v := range allDocs {
		similarity = computeDocSimilarity(query, v, topN, explain)
		if topN == 1 && similarity == 1 {
			break
		}
	}

	//按相似度从高到低排序
	if topN > 1 {
		query.SortSimilarDocs()
	}
	return query
}

func loadDocunentsFromCache(address AddressEntity) []Document {
	cacheKey := buildCacheKey(address)
	if len(cacheKey) == 0 {
		return nil
	}
	docs := make([]Document, 0)
	if !CacheVectorsInMemory { // 从文件读取
		docs = loadDocumentsFromDatabase(cacheKey)
		return docs
	} else { // 从内存读取，如果未缓存到内存，则从文件加载到内存中
		docs = VECTORS_CACHE[cacheKey]
		if docs == nil {
			// TODO
			docs = VECTORS_CACHE[cacheKey]
			if docs == nil {
				docs = loadDocumentsFromDatabase(cacheKey)
				if docs == nil {
					docs = make([]Document, 0)
					VECTORS_CACHE[cacheKey] = docs
				}
			}

			// 为所有词条计算IDF并缓存
			idfs := IdfCache[cacheKey]
			if idfs == nil {
				// TODO
				idfs = IdfCache[cacheKey]
				if idfs == nil {
					termReferences := statInverseDocRefers(docs)
					idfs = make(map[string]float64, len(termReferences))
					for k, v := range termReferences {
						idf := 0.0
						if common.IsAnsiChars(k) || common.IsNumericChars(k) {
							idf = 2.0
						} else {
							idf = math.Log(float64(len(docs) / (v + 1)))
						}
						if idf < 0.0 {
							idf = 0.0
						}
						idfs[k] = idf
					}
					IdfCache[cacheKey] = idfs
				}
			}

			for _, doc := range docs {
				if !doc.Town.IsNil() {
					doc.Town.Idf = idfs[generateIDFCacheEntryKey(doc.Town)]
				}
				if !doc.Village.IsNil() {
					doc.Village.Idf = idfs[generateIDFCacheEntryKey(doc.Village)]
				}
				if !doc.Road.IsNil() {
					doc.Road.Idf = idfs[generateIDFCacheEntryKey(doc.Road)]
				}
				if !doc.RoadNum.IsNil() {
					doc.RoadNum.Idf = idfs[generateIDFCacheEntryKey(doc.RoadNum)]
				}
				for _, term := range doc.Terms {
					term.Idf = idfs[generateIDFCacheEntryKey(term)]
				}
			}
		}
	}
	return docs
}

// TODO
func loadDocumentsFromDatabase(key string) []Document {
	docs := make([]Document, 0)

	return docs
	//List<Document> docs = new ArrayList<Document>();
	//
	//String filePath = getCacheFolder() + "/" + key + ".vt";
	//File file = new File(filePath);
	//if(!file.exists()) return docs;
	//try {
	//	FileInputStream fsr = new FileInputStream(file);
	//	InputStreamReader sr = new InputStreamReader(fsr, "utf8");
	//	BufferedReader br = new BufferedReader(sr);
	//	String line = null;
	//	while((line = br.readLine()) != null){
	//	Document doc = deserialize(line);
	//	if(doc==null) continue;
	//	docs.add(doc);
	//}
	//	br.close();
	//	sr.close();
	//	fsr.close();
	//} catch (Exception ex) {
	//	LOG.error("[doc-vec] [cache] [error] Error in reading file: " + filePath, ex);
	//}
	//
	//return docs;
}

func computeDocSimilarity(query *Query, doc Document, topN int, explain bool) float64 {
	var dterm Term
	//=====================================================================
	//计算text类型词条的稠密度、匹配率
	//1. Text类型词条匹配情况

	qTextTermCount := 0                                    //查询文档Text类型词条数量
	dTextTermMatchCount, matchStart, matchEnd := 0, -1, -1 //地址库文档匹配上的Text词条数量
	for _, v := range query.QueryDoc.Terms {
		if v.Types != enum.Text { //仅针对Text类型词条计算 词条稠密度、词条匹配率
			continue
		}
		qTextTermCount++
		for i := 0; i < len(doc.Terms); i++ {
			term := doc.Terms[i]
			if term.Types != enum.Text { //仅针对Text类型词条计算 词条稠密度、词条匹配率
				continue
			}
			if term.Text == v.Text {
				dTextTermMatchCount++
				if matchStart == -1 {
					matchStart = i
					matchEnd = i
					break
				}
				if i > matchEnd {
					matchEnd = i
				} else if i < matchStart {
					matchStart = i
				}
				break
			}
		}
	}
	//2. 计算稠密度、匹配率
	textTermDensity, textTermCoord := float64(1), float64(1)
	if qTextTermCount > 0 {
		textTermCoord = math.Sqrt(float64(dTextTermMatchCount/qTextTermCount))*0.5 + 0.5
	}
	//词条稠密度：
	// 查询文档a的文本词条为：【翠微西里】
	// 地址库文档词条为：【翠微北里12号翠微嘉园B座西801】
	// 地址库词条能匹配上【翠微西里】的每一个词条，但不是连续匹配，中间间隔了其他词条，稠密度不够，这类文档应当比能够连续匹配上查询文档的权重低
	//稠密度 = 0.7 + (匹配上查询文档的词条数量 / 匹配上的词条起止位置间总词条数量) * 0.3
	//   乘以0.3是为了将稠密度对相似度结果的影响限制在 0 - 0.3 的范围内。
	//假设：查询文档中Text类型的词条为：翠, 微, 西, 里。地址库中有如下两个文档，Text类型的词条为：
	//1: 翠, 微, 西, 里, 10, 号, 楼
	//2: 翠, 微, 北, 里, 89, 号, 西, 2, 楼
	//则：
	// density1 = 0.7 + ( 4/4 ) * 0.3 = 0.7 + 0.3 = 1
	// density2 = 0.7 + ( 4/7 ) * 0.3 = 0.7 + 0.17143 = 0.87143
	// 文档2中 [翠、微、西、里] 4个词匹配上查询文档词条，这4个词条之间共包含7个词条。
	if qTextTermCount >= 2 && dTextTermMatchCount >= 2 {
		textTermDensity = math.Sqrt(float64(dTextTermMatchCount/(matchEnd-matchStart+1)))*0.5 + 0.5
	}
	var simiDoc SimilarDocument
	if explain && topN > 1 {
		simiDoc = NewSimilarDocument(doc)
	}

	//=====================================================================
	//计算TF-IDF和相似度所需的中间值
	var sumQD, sumQQ, sumDD, qtfidf, dtfidf float64 = 0, 0, 0, 0, 0
	var dboost, qboost float64 = 0, 0
	for _, v := range query.QueryDoc.Terms {
		qboost = getBoostValue(false, query.QueryDoc, v, doc, Term{})
		qtfidf = v.Idf * qboost
		dterm = doc.GetTerm(v.Text)
		if dterm == (Term{}) && enum.RoadNum == v.Types {
			if doc.RoadNum != (Term{}) && doc.Road != (Term{}) && &doc.Road == v.Ref {
				dterm = doc.RoadNum
			}
		}
		if dterm == (Term{}) {
			dboost = 0
		} else {
			dboost = getBoostValue(true, query.QueryDoc, v, doc, dterm)
		}
		coord, density := float64(1), float64(1)
		if dterm != (Term{}) && enum.Text == dterm.Types {
			coord = textTermCoord
			density = textTermDensity
		}
		if dterm != (Term{}) {
			dtfidf = dterm.Idf
		} else {
			dtfidf = v.Idf
		}
		dtfidf *= dboost * coord * density

		if explain && topN > 1 && dterm != (Term{}) {
			mt := new(MatchedTerm)
			mt.Boost = dboost
			mt.Tfidf = dtfidf
			if dterm.Types == enum.Text {
				mt.Density = density
				mt.Coord = coord
			} else {
				mt.Density = float64(-1)
				mt.Coord = float64(-1)
			}
			simiDoc.AddMatchedTerm(mt)
		}
		sumQQ += qtfidf * qtfidf
		sumQD += qtfidf * dtfidf
		sumDD += dtfidf * dtfidf
	}
	if sumDD == 0 || sumQQ == 0 {
		return 0
	}
	similarity := sumQD / (math.Sqrt(sumQQ * sumDD))
	if explain && topN > 1 {
		simiDoc.Similarity = similarity
		query.AddSimiDoc(simiDoc)
	} else {
		query.AddSimiDocs(doc, similarity)
	}
	return similarity
}

func buildCacheKey(address AddressEntity) string {
	if address.IsNil() || !address.Province.IsNil() || !address.City.IsNil() {
		return ""
	}
	//var res string
	res := strconv.Itoa(int(address.Province.Id)) + "-" + strconv.Itoa(int(address.City.Id))
	if address.City.Children != nil {
		res += "-" + strconv.Itoa(int(address.District.Id))
	}
	return res
}
