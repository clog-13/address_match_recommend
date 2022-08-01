package test

import (
	"github.com/xiiv13/address_match_recommend/core"
	"testing"
)

func TestReadDat(t *testing.T) {
	core.ReadDat("../../resource/region_2021.dat")
}
func TestDecode(t *testing.T) {
	core.ReadDat("../../resource/region_2021.dat")
}
