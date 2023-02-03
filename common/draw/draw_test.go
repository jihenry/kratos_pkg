package draw

import (
	"fmt"
	"testing"
)

type PrizesType int
type DrawGifts struct {
	ItemID     int     `json:"itemID"`   //道具id
	ItemName   string  `json:"itemName"` //道具名称
	ItemNum    int     `json:"itemNum"`  //道具数量
	Weight     float64 `json:"weight"`   //权重
	GiftType   int     `json:"giftType"` //1 直接发放 2 提取链接 3 实物 收集用户地址 4权益 导出百果园id 发放兑换码
	NeedCoin   int     `json:"needCoin"`
	NeedTicket int     `json:"needTicket"` //用户需要的抽奖券
	Key        int     `json:"key"`        //中奖的游标
}

func TestRandDraw(t *testing.T) {
	maps := map[PrizesType]interface{}{
		//千元水果自由
		PrizesType(1): &DrawGifts{ItemName: "千元水果自由", ItemID: 1, ItemNum: 1, Weight: 0.0001, GiftType: 3, Key: 0},
		//滑板车
		PrizesType(2): &DrawGifts{ItemName: "COOGHI酷骑V1型号滑板车", ItemID: 2, ItemNum: 1, Weight: 0.0002, GiftType: 3, Key: 1},
		//20周年环保袋
		PrizesType(3): &DrawGifts{ItemName: "20周年环保袋", ItemID: 3, ItemNum: 1, Weight: 0.0003, GiftType: 3, Key: 2},
		//智云Smooth X手机云台
		PrizesType(4): &DrawGifts{ItemName: "智云Smooth X手机云台", ItemID: 4, ItemNum: 1, Weight: 0.0004, GiftType: 3, Key: 3},
		//拾柒迷你照片书定制
		PrizesType(5): &DrawGifts{ItemName: "拾柒迷你照片书定制", ItemID: 5, ItemNum: 1, Weight: 0.0005, GiftType: 2, Key: 4},
		//20周年杯盖
		PrizesType(6): &DrawGifts{ItemName: "20周年杯盖", ItemID: 6, ItemNum: 1, Weight: 0.0006, GiftType: 3, Key: 5},
		//熊猫大鲜纤+
		PrizesType(7): &DrawGifts{ItemName: "熊猫大鲜纤+", ItemID: 7, ItemNum: 1, Weight: 0.0007, GiftType: 3, Key: 6},
		//心的抱抱熊曲奇礼盒
		PrizesType(8): &DrawGifts{ItemName: "心的抱抱熊曲奇礼盒", ItemID: 8, ItemNum: 1, Weight: 0.0008, GiftType: 3, Key: 7},
		//20周年币
		PrizesType(9): &DrawGifts{ItemName: "周年币*2", ItemID: 9, ItemNum: 2, Weight: 1, GiftType: 1, Key: 8},
	}
	prizes := make([]*Prize, 0, len(maps))
	for _, v := range maps {
		if v != nil {
			if vv, ok := v.(*DrawGifts); ok {
				prizes = append(prizes, &Prize{
					PlayerID: int64(vv.ItemID),
					Weight:   int(vv.Weight * 10000),
					Key:      vv.Key,
				})
			}
		}
	}
	data := RandDraw(prizes)
	fmt.Println(data)
	//for i := 0; i < 600000; i++ {
	//
	//	if data.PlayerID != 9 {
	//		fmt.Println(data)
	//	}
	//
	//}

}
