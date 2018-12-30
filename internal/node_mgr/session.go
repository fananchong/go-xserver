package nodemgr

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	defaultNodeInterfaceImpl
	gotcp.Session
	Info    *protocol.SERVER_INFO
	msgData []byte
	msgFlag byte
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
	sess.doLose()
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
		sess.Info = msg.GetData()
		sess.msgData = data
		sess.msgFlag = flag
		sess.SetID(utility.ServerID2NodeID(msg.GetData().GetId()))
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
	if utility.EqualSID(sess.Info.GetId(), msg.GetData().GetId()) == false {
		sess.Close()
		return
	}
	sess.Info = msg.GetData()
	sess.msgData = data
	sess.msgFlag = flag
	common.XLOG.Infoln("Node register for me, node id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
	common.XLOG.Infoln(sess.Info)

	xsessionmgr.register(sess)
	xsessionmgr.forAll(func(elem *Session) {
		if elem != sess {
			elem.Send(data, flag)
		}
	})
	xsessionmgr.forAll(func(elem *Session) {
		if elem != sess {
			sess.Send(elem.msgData, elem.msgFlag)
		}
	})
}

func (sess *Session) doLose() {
	if sess.Info != nil {
		xsessionmgr.lose(sess)
		msg := &protocol.MSG_MGR_LOSE_SERVER{}
		msg.Id = sess.Info.GetId()
		msg.Type = sess.Info.GetType()
		xsessionmgr.forAll(func(elem *Session) {
			elem.SendMsg(uint64(protocol.CMD_MGR_LOSE_SERVER), msg)
		})
		common.XLOG.Infoln("lose node, type:", msg.Type, "id:", utility.ServerID2UUID(msg.Id).String())
	}
}

// GetType : 获取节点类型
func (sess *Session) GetType() common.NodeType {
	if sess.Info != nil {
		return common.NodeType(sess.Info.GetType())
	}
	return common.Unknow
}

// GetIP : 获取本节点信息，IP
func (sess *Session) GetIP(i common.IPType) string {
	if sess.Info != nil {
		return sess.Info.GetAddrs()[i]
	}
	return ""
}

// GetPort : 获取本节点信息，端口
func (sess *Session) GetPort(i int) int32 {
	if sess.Info != nil {
		return sess.Info.GetPorts()[i]
	}
	return 0
}

// GetOverload : 获取本节点信息，负载
func (sess *Session) GetOverload(i int) uint32 {
	if sess.Info != nil {
		return sess.Info.GetOverload()[i]
	}
	return 0
}

// GetVersion : 获取本节点信息，版本号
func (sess *Session) GetVersion() string {
	if sess.Info != nil {
		return sess.Info.GetVersion()
	}
	return ""
}

// GetSID : 获取 SID
func (sess *Session) GetSID() *protocol.SERVER_ID {
	if sess.Info != nil {
		return sess.Info.GetId()
	}
	return &protocol.SERVER_ID{}
}
