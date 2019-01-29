package gateway

import (
	"context"
	"net"

	"github.com/fananchong/go-xserver/internal/utility"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/gotcp"
)

// IntranetSession : 登录玩家类
type IntranetSession struct {
	gotcp.Session
	ctx *common.Context
}

// Init : 初始化网络会话节点
func (sess *IntranetSession) Init(root context.Context, conn net.Conn, derived gotcp.ISession, userdata interface{}) {
	sess.ctx = userdata.(*common.Context)
	sess.Session.Init(root, conn, derived)
}

// OnRecv : 接收到网络数据包，被触发
func (sess *IntranetSession) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if sess.IsVerified() == false && sess.doVerify(protocol.CMD_GW_ENUM(cmd), data, flag) == false {
		return
	}
	switch protocol.CMD_GW_ENUM(cmd) {
	default:
		sess.ctx.Log.Errorln("unknow cmd, cmd =", cmd)
	}
}

// OnClose : 断开连接，被触发
func (sess *IntranetSession) OnClose() {

}

func (sess *IntranetSession) doVerify(cmd protocol.CMD_GW_ENUM, data []byte, flag byte) bool {
	if cmd != protocol.CMD_GW_VERIFY_TOKEN {
		sess.ctx.Log.Errorln("The expected message number is `protocol.CMD_GW_VERIFY_TOKEN`(", protocol.CMD_GW_VERIFY_TOKEN, ")")
		sess.Close()
		return false
	}
	msg := &protocol.MSG_GW_VERIFY_TOKEN{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		sess.ctx.Log.Errorln("Message parsing failed, message number is`protocol.CMD_GW_VERIFY_TOKEN`(", protocol.CMD_GW_VERIFY_TOKEN, ")")
		sess.Close()
		return false
	}
	if sess.ctx.Config.Common.IntranetToken != msg.GetToken() {
		sess.ctx.Log.Errorln("Token verification failed with server ID", utility.ServerID2UUID(msg.GetId()).String())
		sess.Close()
		return false
	}
	sess.Verify()
	sess.ctx.Log.Infoln("Token verification succeeded with server ID", utility.ServerID2UUID(msg.GetId()).String())
	return true
}
