package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"zee.com/work/mmo_game/core"
)

// 当客户端建立连接的时候 hook 函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)
	// 同步当前PlayerID 给客户端， 走 MsgID:1 信息
	player.SyncPid()
	// 同步当前玩家的初始化坐标信息给客户端， 走 MsgID:200 消息
	player.BroadCastStartPosition()
	fmt.Println("====> Player pidId = ", player.Pid, " arrived =====")
}

func main() {
	// 创建服务器句柄
	s := znet.NewServer()

	// 注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)

	// 启动服务
	s.Serve()
}
