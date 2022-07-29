package segment

type ISegmenter interface {
	Segment(text string) []string
}

func NewSegment(types string) ISegmenter {
	switch types {
	case "simple":
		return SimpleSegmenter{}
	}
	return SimpleSegmenter{}
}
