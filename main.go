package main

import (
	"github.com/xiiv13/address_match_recommend/utils"
)

const (
	totalNumber     = 10000000
	falseDetectRate = 0.000001
)

var (
	bloom = utils.NewCountingBloomFilter(totalNumber, falseDetectRate)
)

//func main() {
//	var queryAddr string
//	fmt.Scanln("输入国家: ", &queryAddr)
//
//	result := query(queryAddr, 5)
//}
//
//func query(text string, n int) []string {
//	if bloom.BFTest([]byte(text)) { // 布隆过滤器判断存在
//
//	} else {
//		result := core.FindsimilarAddress(text, n, true)
//
//		bloom.BFSet([]byte(result.QueryAddr.AddressText))
//	}
//
//}
