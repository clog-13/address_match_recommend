package model

type MatchedTerm struct {
	Term    Term
	Coord   float64
	Density float64
	Boost   float64
	Tfidf   float64
}
