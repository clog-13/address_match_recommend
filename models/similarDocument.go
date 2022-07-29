package models

type SimilarDocument struct {
	Doc          Document
	MatchedTerms map[string]*MatchedTerm
	Similarity   float64
}

func NewSimilarDocument(doc Document) *SimilarDocument {
	return &SimilarDocument{Doc: doc}
}

func (s SimilarDocument) AddMatchedTerm(value *MatchedTerm) {
	if s.MatchedTerms == nil {
		s.MatchedTerms = make(map[string]*MatchedTerm)
	}
	s.MatchedTerms[value.Term.Text] = value
}
