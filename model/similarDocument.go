package model

type SimilarDocument struct {
	Doc          Document
	MatchedTerms map[string]MatchedTerms
	Similarity   float64
}

func NewSimilarDocument(doc Document) SimilarDocument {
	return SimilarDocument{Doc: doc}
}
