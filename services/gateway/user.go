package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
)

// User : 登录玩家类
type User struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (user *User) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if user.IsVerified() == false && user.doVerify(protocol.CMD_GATEWAY_ENUM(cmd), data, flag) == false {
		return
	}
	switch protocol.CMD_GATEWAY_ENUM(cmd) {
	default:
		Ctx.Log.Errorln("Unknown message number, message number is", cmd)
	}
}

// OnClose : 断开连接，被触发
func (user *User) OnClose() {

}

func (user *User) doVerify(cmd protocol.CMD_GATEWAY_ENUM, data []byte, flag byte) bool {
	return true
}
