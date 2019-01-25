package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
)

// IntranetNode : 登录玩家类
type IntranetNode struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (node *IntranetNode) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if node.IsVerified() == false && node.doVerify(protocol.CMD_GATEWAY_ENUM(cmd), data, flag) == false {
		return
	}
	switch protocol.CMD_GATEWAY_ENUM(cmd) {
	default:
		Ctx.Log.Errorln("unknow cmd, cmd =", cmd)
	}
}

// OnClose : 断开连接，被触发
func (node *IntranetNode) OnClose() {

}

func (node *IntranetNode) doVerify(cmd protocol.CMD_GATEWAY_ENUM, data []byte, flag byte) bool {
	return true
}
