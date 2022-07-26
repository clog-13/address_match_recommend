package model

type iSegmenter interface {
	segment(text string) []string
}

type SimpleSegmenter struct {
}

func (s SimpleSegmenter) Segment(text string) []string {
	if len(text) == 0 {
		return nil
	}
	tokens := make([]string, len(text))
	var digitNum, ansiCharNum int
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= '0' && c <= '9' {
			if ansiCharNum > 0 {
				tokens = append(tokens, text[i-ansiCharNum:i-1])
				ansiCharNum = 0
			}
			digitNum++
			continue
		}
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			if digitNum > 0 {
				tokens = append(tokens, text[i-digitNum:i-1])
				digitNum = 0
			}
			ansiCharNum++
			continue
		}
		if digitNum > 0 || ansiCharNum > 0 { //digitNum, ansiCharNum中只可能一个大于0
			tokens = append(tokens, text[i-digitNum-ansiCharNum:i-1])
			digitNum, ansiCharNum = 0, 0
		}
		tokens = append(tokens, string(c))
	}
	if digitNum > 0 || ansiCharNum > 0 { //digitNum, ansiCharNum中只可能一个大于0
		tokens = append(tokens, text[len(text)-digitNum-ansiCharNum:])
		digitNum, ansiCharNum = 0, 0
	}
	return tokens
}
