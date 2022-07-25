package model

type Term struct {
	RoadNumIdf float64 // 7
	Types      byte
	Text       string
	Idf        float64
	Ref        *Term
}
