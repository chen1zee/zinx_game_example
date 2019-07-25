package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zee.com/work/mmo_game/testprotobuf/protobuf"
)

func main() {
	person := &protobuf.Person{
		Name:   "aaa",
		Age:    16,
		Emails: []string{"ssss@123.com", "aaa@321.com"},
		Phones: []*protobuf.PhoneNumber{
			{Number: "13001919191", Type: protobuf.PhoneType_MOBILE},
			{Number: "13339292929", Type: protobuf.PhoneType_HOME},
			{Number: "32101918178", Type: protobuf.PhoneType_WORK},
		},
	}
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err: ", err)
		return
	}
	newData := &protobuf.Person{}
	err = proto.Unmarshal(data, newData)
	fmt.Println(newData)
	// TODO 五、MMO游戏的Proto3协议
}
