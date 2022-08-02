package core

import (
	"fmt"
	"testing"
)

func TestFindsimilarAddress(t *testing.T) {
	querys := FindsimilarAddress("北京海淀区丹棱街18号创富大厦1106", 5, true)
	fmt.Println(querys)
	for _, v := range querys.SimiDocs {
		fmt.Println(v.Similarity)
		fmt.Println(v.Doc)
	}
}
