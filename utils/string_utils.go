package utils

func IsAnsiChars(text string) bool {
	if len(text) == 0 {
		return false
	}
	for i := 0; i < len(text); i++ {
		c := text[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return true
}

func IsNumericChars(text string) bool {
	if len(text) == 0 {
		return false
	}
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func Head(text string, length int) string {
	if len(text) == 0 || len(text) <= length {
		return text
	}
	if length <= 0 {
		return ""
	}
	return text[:length]
}

func Substring(text string, begin, end int) string {
	if len(text) == 0 {
		return text
	}
	if begin < 0 {
		begin = 0
	}
	if end >= len(text) {
		end = len(text)
	}
	return text[begin:end]
}
