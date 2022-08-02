package main

import (
	"fmt"
	"github.com/xiiv13/address_match_recommend/core"
	"github.com/xiiv13/address_match_recommend/models"
	"github.com/xiiv13/address_match_recommend/utils"
)

const (
	totalNumber     = 10000000
	falseDetectRate = 0.000001
)

var (
	persister = &models.AddressPersister{}
	bloom     = utils.NewCountingBloomFilter(totalNumber, falseDetectRate)
)

func main() {
	var queryAddr string
	//fmt.Scanln("输入地址: ", &queryAddr)
	queryAddr = "四川成都高新博士公馆"
	result := query(queryAddr, 5)
	fmt.Println(len(result.SimiDocs))
	for _, v := range result.SimiDocs {
		fmt.Println(persister.LoadAddr(v.Doc.Id).RawText)
	}
}

func query(text string, n int) models.Query {
	return core.FindsimilarAddress(text, n, true)

	//if bloom.BFTest([]byte(text)) { // 布隆过滤器判断存在
	//
	//} else {
	//	result := core.FindsimilarAddress(text, n, true)
	//
	//	bloom.BFSet([]byte(result.QueryAddr.AddressText))
	//}
}
