package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstring(t *testing.T) {
	assert.Equal(t, "cde", Substring([]rune("abcdef"), 2, 4))
	assert.Equal(t, "def", Substring([]rune("abcdef"), 3, 5))
	assert.Equal(t, "a", Substring([]rune("abcdef"), 0, 0))
	assert.Equal(t, "d", Substring([]rune("abcdef"), 3, 3))
	assert.Equal(t, "f", Substring([]rune("abcdef"), 5, 5))
	assert.Equal(t, "def", Substring([]rune("abcdef"), 3, 9))
	assert.Equal(t, "", Substring([]rune("abcdef"), 9, 1))
	assert.Equal(t, "", Substring([]rune{}, 2, 5))
	assert.Equal(t, "", Substring([]rune("abcdef"), 7, 9))
	assert.Equal(t, "", Substring([]rune("零一二三四五"), 1, 3))
	assert.Equal(t, "二三四", Substring([]rune("零一二三四五"), 2, 4))
	assert.Equal(t, "三四五", Substring([]rune("零一二三四五"), 3, 5))
	assert.Equal(t, "零", Substring([]rune("零一二三四五"), 0, 0))
	assert.Equal(t, "三", Substring([]rune("零一二三四五"), 3, 3))
	assert.Equal(t, "五", Substring([]rune("零一二三四五"), 5, 5))
	assert.Equal(t, "三四五", Substring([]rune("零一二三四五"), 3, 9))
	assert.Equal(t, "", Substring([]rune("零一二三四五"), 9, 1))
	assert.Equal(t, "", Substring([]rune{}, 2, 5))
	assert.Equal(t, "", Substring([]rune("零一二三四五"), 7, 9))

}

func TestIsNumericChars(t *testing.T) {
	assert.False(t, IsNumericChars(""))
	assert.False(t, IsNumericChars("abc"))
	assert.False(t, IsNumericChars("19abc"))
	assert.False(t, IsNumericChars("19a771"))

	assert.True(t, IsNumericChars("0123"))
	assert.True(t, IsNumericChars("789"))
	assert.True(t, IsNumericChars("09"))
}

func TestRemove(t *testing.T) {
	s := Remove([]rune("江苏泰州兴化市昌荣镇【康琴网吧】 (昌荣镇附近)"), []rune("　 \r\n\t,，。·.．;；:：、！@$%*^`~=+&'\"|_-\\/"), "")
	fmt.Println(s)
	s = Remove([]rune("江苏泰州兴化市昌荣镇【康琴网吧】(昌荣镇附近)"), []rune("{}【】〈〉<>[]「」“”（）()"), "")
	fmt.Println(s)
}
