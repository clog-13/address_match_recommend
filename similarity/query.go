package similarity

import "address_match_recommend/models"

type Query struct {
	topN      int
	QueryAddr models.AddressEntity
	QueryDoc  Document

	SimiDocs []SimilarDocument
}

func NewQuery(topN int) *Query {
	return &Query{
		topN: topN,
	}
}

// SortSimilarDocs 将相似文档按相似度从高到低排序。
func (q Query) SortSimilarDocs() {
	if len(q.SimiDocs) == 0 {
		return
	}
	exchanged := true
	endIndex := len(q.SimiDocs) - 1
	for exchanged {
		exchanged = false
		for i := 1; i <= endIndex; i++ {
			if q.SimiDocs[i-1].Similarity < q.SimiDocs[i].Similarity {
				temp := q.SimiDocs[i-1]
				q.SimiDocs[i-1] = q.SimiDocs[i]
				q.SimiDocs[i] = temp
				exchanged = true
			}
		}
		endIndex--
	}
}

// AddSimiDoc 添加一个相似文档, 只保留相似度最高的top N条相似文档,相似度最低的从simiDocs中删除
func (q Query) AddSimiDoc(simiDoc SimilarDocument) bool {
	if simiDoc.Similarity <= 0 {
		return false
	}
	if q.SimiDocs == nil {
		q.SimiDocs = make([]SimilarDocument, q.topN)
	}
	if len(q.SimiDocs) < q.topN {
		q.SimiDocs = append(q.SimiDocs, simiDoc)
		return true
	}
	minSimilarityIndex := 0
	for i := 1; i < q.topN; i++ {
		if q.SimiDocs[i].Similarity < q.SimiDocs[minSimilarityIndex].Similarity {
			minSimilarityIndex = i
		}
	}
	if q.SimiDocs[minSimilarityIndex].Similarity < simiDoc.Similarity {
		q.SimiDocs[minSimilarityIndex] = simiDoc
		return true
	}
	return false
}

func (q Query) AddSimiDocs(doc Document, similarity float64) bool {
	if similarity <= 0 {
		return false
	}
	if q.SimiDocs == nil {
		q.SimiDocs = make([]SimilarDocument, q.topN)
		simiDoc := NewSimilarDocument(doc)
		simiDoc.Similarity = similarity
		q.SimiDocs = append(q.SimiDocs, simiDoc)
		return true
	}
	if q.SimiDocs[0].Similarity < similarity {
	}
	simiDoc := NewSimilarDocument(doc)
	simiDoc.Similarity = similarity
	q.SimiDocs[0] = simiDoc
	return true
}

func (q Query) GetSimilarDoce() []SimilarDocument {
	if q.SimiDocs == nil {
		q.SimiDocs = make([]SimilarDocument, 0)
	}
	return q.SimiDocs
}
