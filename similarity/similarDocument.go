package similarity

import "address_match_recommend/models"

type SimilarDocument struct {
	Doc          models.Document
	MatchedTerms map[string]*MatchedTerm
	Similarity   float64
}

func NewSimilarDocument(doc models.Document) *SimilarDocument {
	return &SimilarDocument{Doc: doc}
}

func (s SimilarDocument) AddMatchedTerm(value *MatchedTerm) {
	if s.MatchedTerms == nil {
		s.MatchedTerms = make(map[string]*MatchedTerm)
	}
	s.MatchedTerms[value.Term.Text] = value
}
