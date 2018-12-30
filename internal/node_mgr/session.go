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
	id common.NodeID
	t  common.NodeType
}

// OnRecv : 接收到网络数据包，被触发
func (sess *Session) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	if sess.IsVerified() == false && sess.doVerify(cmd, data, flag) == false {
		return
	}
	switch cmd {
	case uint64(protocol.CMD_MGR_REGISTER_SERVER):
		sess.doRegister(data, flag)
	case uint64(protocol.CMD_MGR_PING):
		// do nothing
	default:
		common.XLOG.Errorln("unknow cmd, cmd =", cmd)
	}
}

// OnClose : 断开连接，被触发
func (sess *Session) OnClose() {
	sess.doLose(sess.id, sess.t)
}

func (sess *Session) doVerify(cmd uint64, data []byte, flag byte) bool {
	if cmd == uint64(protocol.CMD_MGR_REGISTER_SERVER) {
		msg := &protocol.MSG_MGR_REGISTER_SERVER{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Close()
			return false
		}
		if msg.GetToken() != common.XCONFIG.Common.IntranetToken {
			common.XLOG.Errorln("IntranetToken error!")
			common.XLOG.Errorln("Msg token:", msg.GetToken())
			common.XLOG.Errorln("Expect token:", common.XCONFIG.Common.IntranetToken)
			sess.Close()
			return false
		}
		sess.id = utility.ServerID2NodeID(msg.GetData().GetId())
		sess.t = common.NodeType(msg.GetData().GetType())
		sess.Verify()
		return true
	}
	common.XLOG.Errorln("Before message[CMD_MGR_REGISTER_SERVER], recv cmd:", cmd)
	sess.Close()
	return false
}

func (sess *Session) doRegister(data []byte, flag byte) {
	msg := &protocol.MSG_MGR_REGISTER_SERVER{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		sess.Close()
		return
	}
	common.XLOG.Infoln("Node register for me, node id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
	common.XLOG.Infoln(msg)

	xsessionmgr.register(msg.GetData())

	xsessionmgr.forAll(func(info *protocol.SERVER_INFO) {
		// TODO:
	})

	xsessionmgr.forAll(func(info *protocol.SERVER_INFO) {
		// TODO:
	})
}

func (sess *Session) doLose(nid common.NodeID, t common.NodeType) {
	xsessionmgr.lose(nid, t)
	xsessionmgr.forAll(func(info *protocol.SERVER_INFO) {
		// TODO:
	})
}
