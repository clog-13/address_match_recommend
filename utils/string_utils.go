package utils

import "strings"

func Head(text string, length int) string {
	if len(text) == 0 || len(text) <= length {
		return text
	}
	if length <= 0 {
		return ""
	}
	return text[:length]
}

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

func Remove(text string, chars []byte, exclude string) string {
	if len(text) == 0 || len(chars) == 0 {
		return text
	}
	charSet := make(map[byte]struct{}, 0)
	for _, v := range chars {
		charSet[v] = struct{}{}
	}
	var sb strings.Builder
	var removed bool
	for _, v := range text {
		_, ok := charSet[byte(v)]
		if ok && !strings.Contains(exclude, string(v)) {
			removed = true
			continue
		}
		sb.WriteRune(v)
	}
	if removed {
		return sb.String()
	} else {
		return text
	}
}

func RemoveRepeatNum(text string, n int) string {
	if len(text) < n {
		return text
	}
	var sb strings.Builder
	var cnt int
	for i, v := range text {
		if v >= '0' && v <= '9' {
			cnt++
			continue
		}
		if cnt > 0 && cnt < n {
			sb.WriteString(text[i-cnt : i])
		}
		cnt = 0
		sb.WriteRune(v)
	}
	if cnt > 0 && cnt < n {
		sb.WriteString(text[len(text)-cnt:])
	}

	return sb.String()
}

func Tail(text string, length int) string {
	if len(text) == 0 {
		return text
	}
	if length <= 0 {
		return ""
	}
	return text[len(text)-length:]
}

func Take(text string, begin int) string {
	if len(text) == 0 {
		return text
	}
	if begin <= 0 {
		begin = 0
	}
	if begin > len(text)-1 {
		return ""
	}

	return text[begin:]
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
	if begin > end {
		return ""
	}

	return text[begin : end+1]
}
