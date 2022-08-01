package segment

type ISegmenter interface {
	Segment(text string) []string
}
