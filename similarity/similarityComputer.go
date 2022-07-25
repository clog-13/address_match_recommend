package similarity

import (
	"address_match_recommend/enum"
	"address_match_recommend/model"
	"math"
)

func findsimilarAddress(addressText string, topN int, explain bool) *model.Query {
	query := model.NewQuery(topN)

	//解析地址 TODO
	queryAddr := model.FormatAddressEntity(addressText)

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

func loadDocunentsFromCache(queryAddr model.AddressEntity) []model.Document {

}

func analyse(queryAddr model.AddressEntity) model.Document {

}

func computeDocSimilarity(query *model.Query, doc model.Document, topN int, explain bool) float64 {
	var dterm model.Term
		//=====================================================================
		//计算text类型词条的稠密度、匹配率
		//1. Text类型词条匹配情况


		 qTextTermCount := 0 //查询文档Text类型词条数量
		dTextTermMatchCount, matchStart, matchEnd :=0,-1, -1 //地址库文档匹配上的Text词条数量
		for _, v := range query.QueryDoc.Terms {
			if v.Types!=enum.TEXT { //仅针对Text类型词条计算 词条稠密度、词条匹配率
				continue
			}
			qTextTermCount++
			for i:=0; i< len(doc.Terms);i++ {
				term := doc.Terms[i]
				if term.Types != enum.TEXT { //仅针对Text类型词条计算 词条稠密度、词条匹配率
					continue
				}
				if term.Text==v.Text {
					dTextTermMatchCount++
					if matchStart ==-1{
						matchStart =i
							matchEnd = i
							break
					}
					if i>matchEnd {
						matchEnd = i
					}else if i< matchStart{
						matchStart = i
					}
					break
				}
			}
		}
		//2. 计算稠密度、匹配率
		textTermDensity, textTermCoord := float64(1), float64(1)
		if qTextTermCount>0 {
			textTermCoord = math.Sqrt(float64(dTextTermMatchCount/qTextTermCount))*0.5+0.5
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
		if qTextTermCount>=2 && dTextTermMatchCount >=2 {
			textTermDensity = math.Sqrt(float64(dTextTermMatchCount/(matchEnd-matchStart+1)))*0.5+0.5
		}
		var simiDoc model.SimilarDocument
		if explain && topN>1 {
			simiDoc = model.NewSimilarDocument(doc)
		}

		//=====================================================================
		//计算TF-IDF和相似度所需的中间值
		var sumQD, sumQQ,sumDD,qtfidf,dtfidf float64=0,0,0,0,0
		var dboost, qboost float64 = 0, 0
		for _, v := range query.QueryDoc.Terms {
			qboost = getBoostValue(false, query.QueryDoc, v, doc, model.Term{})
			qtfidf = v.Idf*qboost
			dterm = doc.GetTerm(v.Text)
			if dterm==model.Term{} && enum.
		}

		for(Term qterm : query.getQueryDoc().getTerms()) {
			qboost = getBoostValue(false, query.getQueryDoc(), qterm, doc, null);
			qtfidf = qterm.getIdf() * qboost;
			dterm = doc.getTerm(qterm.getText());
			if(dterm==null && TermType.RoadNum==qterm.getType()){
				//从b中找门牌号词条
				if(doc.getRoadNum()!=null && doc.getRoad()!=null && doc.getRoad().equals(qterm.getRef()))
					dterm = doc.getRoadNum();
			}
			dboost = dterm==null ? 0 : getBoostValue(true, query.getQueryDoc(), qterm, doc, dterm);
			double coord = (dterm!=null && TermType.Text==dterm.getType()) ? textTermCoord : 1;
			double density = (dterm!=null && TermType.Text==dterm.getType()) ? textTermDensity : 1;
			dtfidf = (dterm!=null ? dterm.getIdf() : qterm.getIdf()) * dboost * coord * density;

			if(explain && topN>1 && dterm!=null){
				MatchedTerm mt = null;
				mt = new MatchedTerm(dterm);
				mt.setBoost(dboost);
				mt.setTfidf(dtfidf);
				if(dterm.getType()==TermType.Text){
					mt.setDensity(density);
					mt.setCoord(coord);
				}else{
					mt.setDensity(-1);
					mt.setCoord(-1);
				}
				simiDoc.addMatchedTerm(mt);
			}

			sumQQ += qtfidf * qtfidf;
			sumQD += qtfidf * dtfidf;
			sumDD += dtfidf * dtfidf;
		}
		if(sumDD==0 || sumQQ==0) return 0;

		double similarity = sumQD / ( Math.sqrt(sumQQ * sumDD) );
		if(explain && topN>1){
			simiDoc.setSimilarity(similarity);
			query.addSimiDoc(simiDoc);
		}else query.addSimiDoc(doc, similarity);
		return similarity;
}

func getBoostValue(forDoc bool, qdoc model.Document, qterm model.Term, ddoc model.Document ,dterm model.Term) float64 {

}
