package index

import (
	"fmt"
	"github.com/xiiv13/address_match_recommend/models"
	"testing"
)

func TestIsAcceptableItemType(t *testing.T) {
	fmt.Println(IsAcceptableItemType('1'))
	fmt.Println(IsAcceptableItemType(49))
	fmt.Println(IsAcceptableItemType(models.ProvinceTerm))
}
