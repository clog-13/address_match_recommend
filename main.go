package main

import (
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

	addr := AddrApi(inputAddr.String())
	fmt.Println("用户最终收货地址：", addr)
}
