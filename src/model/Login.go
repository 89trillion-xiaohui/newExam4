package model

type User struct {
	UID      string `bson:"uid"`
	GoldCoin int    `bson:"goldcoin"`
	Diamond  int    `bson:"diamond"`
}
