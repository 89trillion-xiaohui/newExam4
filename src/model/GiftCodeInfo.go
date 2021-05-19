package model

// GiftCodeInfo 礼品码信息
type GiftCodeInfo struct {
	GiftCode         string       `json:"gift_code"`                 //礼品码
	Description      int          `json:"description,string"`        //礼品描述
	Times            int          `json:"times,string"`              //可领取次数
	ExpiryDate       string       `json:"expiry_date,string"`        //礼品码有效期
	User             string       `json:"user,string"`               //创建人
	Date             string       `json:"date,string"`               //创建时间
	TimesHasReceived int          `json:"times_has_received,string"` //已领取次数
	GiftText         GiftInfo     `json:"gift_text,string"`          //礼包内容
	ListReceived     ListReceived `json:"list_received"`             //领取列表
}
