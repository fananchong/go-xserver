package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	*gotcp.Session
	NID common.NodeID
}

// NewSession : 网络会话类的构造函数
func NewSession() *Session {
	nid := utility.NewNID()
	common.XLOG.Infoln("NODE ID:", utility.NodeID2UUID(nid).String())
	return &Session{
		NID: nid,
	}
}

// Start : 开始连接 Mgr Server
func (sess *Session) Start() bool {
	sess.connectMgrServer()
	return true
}

func (sess *Session) connectMgrServer() {
	sess.Session = &gotcp.Session{}
TRY_AGAIN:
	addr, port := getMgrInfoByBlock()
	if sess.Connect(fmt.Sprintf("%s:%d", addr, port), sess) == false {
		time.Sleep(1 * time.Second)
		goto TRY_AGAIN
	}
	sess.Verify()
	sess.registerSelf()
}

func (sess *Session) registerSelf() {
	msg := &protocol.MSG_MGR_REGISTER_SERVER{}
	msg.Data = &protocol.SERVER_INFO{}
	msg.Data.Id = utility.NodeID2ServerID(sess.NID)
	msg.Data.Type = uint32(common.XNODE.GetType())
	msg.Data.Addrs = []string{utility.GetIPInner(), utility.GetIPOuter()}
	msg.Data.Ports = common.XCONFIG.Network.Port

	// TODO: 后续支持
	// msg.Data.Overload
	// msg.Data.Version

	msg.Token = common.XCONFIG.Common.IntranetToken
	sess.SendMsg(uint64(protocol.CMD_MGR_REGISTER_SERVER), msg)
	common.XLOG.Infoln("register self to mgr server. self info:", msg.GetData())
}

// OnRecv : 接收到网络数据包，被触发
func (sess *Session) OnRecv(data []byte, flag byte) {
	cmd := gotcp.GetCmd(data)
	switch cmd {
	case uint64(protocol.CMD_MGR_REGISTER_SERVER):
		msg := &protocol.MSG_MGR_REGISTER_SERVER{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			return
		}
		common.XLOG.Infoln("one server register. id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
	case uint64(protocol.CMD_MGR_LOSE_SERVER):
		msg := &protocol.MSG_MGR_LOSE_SERVER{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			return
		}
		common.XLOG.Infoln("one server lose. id:", utility.ServerID2UUID(msg.GetId()).String())
	case uint64(protocol.CMD_MGR_PING):
		// do nothing
	default:
		common.XLOG.Errorln("unknow cmd, cmd =", cmd)
	}
}

// OnClose : 断开连接，被触发
func (sess *Session) OnClose() {
	go func() {
		time.Sleep(1 * time.Second)
		sess.connectMgrServer()
	}()
}
