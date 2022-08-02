package models

type Query struct {
	TopN      int
	QueryAddr Address
	QueryDoc  Document

	SimiDocs []*SimilarDocument
}

// AddSimiDoc 添加一个相似文档, 只保留相似度最高的top N条相似文档,相似度最低的从simiDocs中删除
func (q *Query) AddSimiDoc(simiDoc *SimilarDocument) {
	if simiDoc.Similarity <= 0 {
		return
	}
	if q.SimiDocs == nil {
		q.SimiDocs = make([]*SimilarDocument, 0)
	}

	if len(q.SimiDocs) < q.TopN {
		q.SimiDocs = append(q.SimiDocs, simiDoc)
		return
	}
	minSimilarityIndex := 0
	for i := 1; i < q.TopN; i++ {
		if q.SimiDocs[i].Similarity < q.SimiDocs[minSimilarityIndex].Similarity {
			minSimilarityIndex = i
		}
	}
	if simiDoc.Similarity > q.SimiDocs[minSimilarityIndex].Similarity {
		q.SimiDocs[minSimilarityIndex] = simiDoc
	}
}

func (q *Query) AddSimiDocs(doc Document, simi float64) {
	if simi <= 0 {
		return
	}

	if q.SimiDocs == nil {
		q.SimiDocs = make([]*SimilarDocument, 0)
		simiDoc := NewSimilarDocument(doc)
		simiDoc.Similarity = simi
		q.SimiDocs = append(q.SimiDocs, simiDoc)
		return
	}
	if q.SimiDocs[0].Similarity < simi {
		simiDoc := NewSimilarDocument(doc)
		simiDoc.Similarity = simi
		q.SimiDocs[0] = simiDoc
	}
}
