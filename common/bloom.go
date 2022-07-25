package common

import (
	"hash"
	"hash/fnv"
	"math"
)

type CBF struct { // Bloom Filter
	h            hash.Hash32
	mod, hashCnt int // 模，hash函数数量
	list         []int8
}

func NewCountingBloomFilter(totalNumber int, falseDetectRate float64) *CBF {
	b := &CBF{h: fnv.New32()}
	b.estimateMK(totalNumber, falseDetectRate)
	b.list = make([]int8, b.mod)
	return b
}

func (b *CBF) estimateMK(number int, possibility float64) {
	// 根据概率公式计算合适的 模长 和 hash函数个数
	//mod = -1 * (n * lnP)/(ln2)^2
	nFloat := float64(number)
	ln2 := math.Log(2)
	b.mod = int(-1 * (nFloat * math.Log(possibility)) / math.Pow(ln2, 2))

	//hashCnt = mod/n * ln2
	b.hashCnt = int(math.Ceil(float64(b.mod) / nFloat * ln2))
}

func (b *CBF) hashFun(fnIdx int, data []byte) int {
	// 使用FNV-1算法计算的hash值
	b.h.Reset()            // 置为0
	_, _ = b.h.Write(data) // 计算hash
	hasInt := int(b.h.Sum32())
	return (hasInt + fnIdx) % b.mod // 加fnIdx，减少碰撞
}

func (b *CBF) BFSet(str []byte) {
	for i := 0; i < b.hashCnt; i++ {
		idx := b.hashFun(i, str)
		b.list[idx] = 1
	}
}

func (b *CBF) BFTest(str []byte) bool {
	for i := 0; i < b.hashCnt; i++ {
		idx := b.hashFun(i, str)
		if b.list[idx] == 0 { // 只要有一个不存在，则必不存在
			return false
		}
	}
	return true
}
