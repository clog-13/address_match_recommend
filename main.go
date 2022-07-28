package main

import (
	"address_match_recommend/core"
	"address_match_recommend/utils"
	"fmt"
	"strings"
)

func main() {
	var inputAddr strings.Builder
	var tmp string
	fmt.Scanln("输入国家: ", &tmp)
	inputAddr.WriteString(tmp)
	fmt.Scanln("输入省份: ", &inputAddr)
	inputAddr.WriteString(tmp)
	fmt.Scanln("输入市: ", &inputAddr)
	inputAddr.WriteString(tmp)
	fmt.Scanln("输入行政区: ", &inputAddr)
	inputAddr.WriteString(tmp)
	fmt.Scanln("输入具体地址: ", &inputAddr)
	inputAddr.WriteString(tmp)

	result := core.FindsimilarAddress(inputAddr.String(), 5, false)
	//addr := AddrApi(inputAddr.String())

	fmt.Println("用户最终收货地址：", result.SimiDocs)
}

const (
	totalNumber     = 10000000
	falseDetectRate = 0.000001
	// 相似判断范围
	level = 2 // 0:Country, 1:Province, 2:City, 3:Barrio ,4:Local
)

var (
	bloom = utils.NewCountingBloomFilter(totalNumber, falseDetectRate)
)

//func AddrApi(addrStr string) string {
//	userInputAddr := str2struct(addrStr)
//	if bloom.BFTest([]byte(addrStr)) { // 布隆过滤器判断存在
//		// 数据库查询
//		err := model.DB.Where(&userInputAddr).Find(&model.DeliveryAddr{}).Error
//		if err != nil {
//			if err == gorm.ErrRecordNotFound { // 数据库查询不存在
//				goto findNear
//			} else { // 数据库其他错误
//				return "-1"
//			}
//		}
//		return addrStr // 用户输入的地址存在，返回用户输入
//	}
//
//findNear: // 用户输入不存在
//	nearAddrs := []string{addrStr}
//	// 返回最相近的5个给用户供其选择
//	nearAddrs = append(nearAddrs, getNearAddr(addrStr, 5)...)
//	fmt.Println(nearAddrs) // 输出相似地址
//	var num int
//	fmt.Scanln(&num)
//
//	if num == 0 { // 用户坚持使用自己输入的地址，则将该地址存入系统的预存数据中
//		result := model.DB.Create(&userInputAddr)
//		if result.Error != nil {
//			return "-1" // 插入新地址失败
//		}
//		bloom.BFSet([]byte(addrStr))
//		return addrStr
//	} else { // 用户使用推荐的地址
//		return nearAddrs[num]
//	}
//}
//
//func getNearAddr(addrStr string, n int) []string {
//
//}
//
//func str2struct(str string) model.DeliveryAddr {
//	arr := strings.Split(str, "-")
//	return model.DeliveryAddr{
//		Country:  arr[0],
//		Province: arr[1],
//		City:     arr[2],
//		Barrio:   arr[3],
//		Local:    arr[4],
//	}
//}
