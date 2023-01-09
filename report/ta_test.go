package report

import (
	"context"
	"testing"
)

const (
	cstEventGrantAward EventType = "grant_award"
)

var taEventConfig = map[EventType]EventConfig{
	cstEventGrantAward: taGrantAwardEvent{},
}

type taGrantAwardEvent struct {
	ServerXID        string `mapstructure:"server_id"`
	ActivityXID      string `mapstructure:"activity_id"`
	PrizeID          uint64 `mapstructure:"award_id"`
	AwardType        int32  `mapstructure:"award_type"`
	AwardName        string `mapstructure:"award_name"`
	IsPlantTreeAward bool   `mapstructure:"is_plant_tree_award"`
	AwardFrom        string `mapstructure:"award_from"` //奖品来源
}

func TestEvent(t *testing.T) {
	client, err := NewThinkingDataClient("http://81.71.18.93", "6fd1b3853e2948b7868bd9b0953d8296", WithEventCfgs(taEventConfig), WithMode(ModeEach))
	if err != nil {
		t.Errorf("NewThinkingDataClient err:%s", err)
		return
	}
	err = client.ReportEvent(context.Background(), "cemls8pdrdn2479lp3s0", cstEventGrantAward, taGrantAwardEvent{
		ServerXID:        "cejs21fqq0m1lccujs5g",
		ActivityXID:      "1296",
		PrizeID:          0x13a72,
		AwardType:        2,
		AwardName:        "合水果游戏测试奖励03",
		IsPlantTreeAward: false,
		AwardFrom:        "game_watermelon",
	})
	if err != nil {
		t.Errorf("NewThinkingDataClient err:%s", err)
	}
}
