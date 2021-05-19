package model

// ListReceived 领取列表
type ListReceived struct {
	UsersReceived string `json:"users_received"` //领取用户
	DateReceived  string `json:"date_received"`  //领取时间

}
