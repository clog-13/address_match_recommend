package models

// MatchedTerm 词条的匹配信息
type MatchedTerm struct {
	Term    Term    // 匹配的词条
	Coord   float64 // 匹配率
	Density float64 // 稠密度
	Boost   float64 // 权重
	TfIdf   float64 // 特征值 TF-IDF
}

func NewMatchedTerm(t Term) *MatchedTerm {
	return &MatchedTerm{
		Term: t,
	}
}
