package model

// GiftInfo 礼包内容
type GiftInfo struct {
	GoldCoin int `json:"gold_coin"` //金币
	Diamond  int `json:"diamond"`   //钻石
	Props    int `json:"props"`     //道具
	Legend   int `json:"legend"`    //英雄
	Pawn     int `json:"pawn"`      //小兵
}
