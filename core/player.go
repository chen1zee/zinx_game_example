package core

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"zee.com/work/mmo_game/pb"
)

// 玩家对象
type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家的连接
	X    float32            // 平面X坐标
	Y    float32            // 高度
	Z    float32            // 平面y坐标 (注意不是 Y)
	V    float32            // 旋转0-360度
}

/*
	Player ID 生成器
*/
var PidGen int32 = 1  // 用来生成玩家ID的计数器
var IdLock sync.Mutex // 保护PidGen的互斥机制

// 创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个PID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在 160坐标点 基于X轴偏移若干坐标
		Y:    0,                            // 高度0
		Z:    float32(134 + rand.Intn(17)), // 随机在 134 坐标点， 基于Y轴偏移若干坐标
		V:    0,                            // 角度为 0 , 尚未实现
	}
	return p
}

func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	fmt.Printf("before Marshal data = %+v\n", data)
	// 将proto Message 结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	fmt.Printf("after Marshal data = %+v \n", msg)
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	// 调用 Zinx 框架的SendMsg 发包
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
	return
}

// 告知客户端pid, 同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组建MsgID 0 proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 发送给客户端
	p.SendMsg(1, data)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // Tp2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, msg)
}