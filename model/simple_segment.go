package model

type iSegmenter interface {
	Segment(text string) []string
}

type SimpleSegmenter struct {
}

// TODO test

// Segment 简单分词器, 直接按单个字符切分, 连续出现的数字、英文字母会作为一个词条
func (s SimpleSegmenter) Segment(text string) []string {
	if len(text) == 0 {
		return nil
	}
	tokens := make([]string, 0)
	var digitNum, ansiCharNum int
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= '0' && c <= '9' {
			if ansiCharNum > 0 {
				tokens = append(tokens, text[i-ansiCharNum:i])
				ansiCharNum = 0
			}
			digitNum++
			continue
		}
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			if digitNum > 0 {
				tokens = append(tokens, text[i-digitNum:i])
				digitNum = 0
			}
			ansiCharNum++
			continue
		}

		if digitNum > 0 || ansiCharNum > 0 { // digitNum, ansiCharNum中只可能一个大于0
			tokens = append(tokens, text[i-digitNum-ansiCharNum:i])
			digitNum, ansiCharNum = 0, 0
		}
		tokens = append(tokens, string(c))
	}

	if digitNum > 0 || ansiCharNum > 0 { // digitNum, ansiCharNum中只可能一个大于0
		tokens = append(tokens, text[len(text)-digitNum-ansiCharNum:])
		digitNum, ansiCharNum = 0, 0
	}
	return tokens
}
