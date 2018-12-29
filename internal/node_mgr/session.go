package nodemgr

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	gotcp.Session
}

// OnRecv : 接收到网络数据包，被触发
func (sess *Session) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if sess.IsVerified() == false {
		if cmd == uint64(protocol.CMD_MGR_REGISTER_SERVER) {
			msg := &protocol.MSG_MGR_REGISTER_SERVER{}
			if gotcp.DecodeCmd(data, flag, msg) == nil {
				sess.Close()
				return
			}
			if msg.GetToken() != common.XCONFIG.Common.IntranetToken {
				common.XLOG.Errorln("IntranetToken error!")
				common.XLOG.Errorln("Msg token:", msg.GetToken())
				common.XLOG.Errorln("Expect token:", common.XCONFIG.Common.IntranetToken)
				sess.Close()
				return
			}
			sess.Verify()
		} else {
			common.XLOG.Errorln("Before message[CMD_MGR_REGISTER_SERVER], recv cmd:", cmd)
			sess.Close()
			return
		}
	}

	switch cmd {
	case uint64(protocol.CMD_MGR_REGISTER_SERVER):
		msg := &protocol.MSG_MGR_REGISTER_SERVER{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Close()
			return
		}
		common.XLOG.Infoln("Node register for me, node id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
		common.XLOG.Infoln(msg)
	}
}

// OnClose : 断开连接，被触发
func (sess *Session) OnClose() {
}
