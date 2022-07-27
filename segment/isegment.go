package segment

type iSegmenter interface {
	Segment(text string) []string
}
