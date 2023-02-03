package draw

import (
	"math/rand"
	"sort"
	"time"
)

type Prize struct {
	PlayerID int64
	Weight   int
	Key      int
	Data     interface{} //用户原始参数
}

func RandDraw(prizes []*Prize) *Prize {
	var (
		//权重累加求和
		weightTotal int
	)
	for _, v := range prizes {
		weightTotal += v.Weight
	}
	//生成一个权重随机数,介于0-总权重和之间
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(weightTotal)
	//权重组重组并排序
	sort.Slice(prizes, func(i, j int) bool {
		return prizes[i].Weight > prizes[j].Weight
	})
	var (
		key int
	)
	for k := range prizes {
		weight := prizes[k].Weight
		if randomNum <= weight {
			key = k
			break
		}
		randomNum -= weight
	}
	// 去除对应奖项 从奖项数组中取出本次抽奖结果
	res := prizes[key]
	return res
}
