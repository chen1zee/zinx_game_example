package api

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"zee.com/work/mmo_game/core"
	"zee.com/work/mmo_game/pb"
)

// 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (*MoveApi) Handle(request ziface.IRequest) {
	// 1 将客户端传来的proto 协议解码
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}
	// 2 得知当前消息是从哪个玩家传递来的， 从连接属性 pid 中获取
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}
	fmt.Printf("user pid = %d , move (%f, %f, %f, %f)", pid, msg.X, msg.Y, msg.Z, msg.V)
	// 3 根据 pid 得到 player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 4 让player 对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
