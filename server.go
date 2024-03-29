package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"zee.com/work/mmo_game/api"
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
	// 将当前新上线玩家添加到 worldManager 中
	core.WorldMgrObj.AddPlayer(player)
	// 将该连接绑定属性 Pid
	conn.SetProperty("pid", player.Pid)
	// 同步周边玩家 上线信息， 与 显示周边玩家信息
	player.SyncSurrounding()
	fmt.Println("====> Player pidId = ", player.Pid, " arrived =====")
}

func main() {
	// 创建服务器句柄
	s := znet.NewServer()
	// 注册路由
	s.AddRouter(2, &api.WorldChatApi{}) // 聊天
	s.AddRouter(3, &api.MoveApi{})      // 移动
	// 注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)

	// 启动服务
	s.Serve()
}
