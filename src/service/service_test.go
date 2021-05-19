package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"

	model2 "test3/src/proto"
)

func TestConvert(t *testing.T) {
	reward := model2.GeneralReward{}
	bytes := []byte{}
	err := proto.Unmarshal(bytes, &reward)
	if err == nil {
		fmt.Println("decode : ", reward)
	}
}
