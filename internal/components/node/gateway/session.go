package nodegateway

import (
	"context"
	"net"

	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	*nodecommon.SessionBase
}

// Init : 初始化网络会话节点
func (sess *Session) Init(root context.Context, conn net.Conn, derived gotcp.ISession, userdata interface{}) {
	ud := userdata.(*nodecommon.UserData)
	sess.SessionBase = nodecommon.NewSessionBase(ud.Ctx, sess)
	sess.SessionBase.Init(root, conn, derived)
	sess.SessMgr = ud.SessMgr
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Info = msg.GetData()
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	if utility.EqualSID(sess.Info.GetId(), msg.GetData().GetId()) == false {
		sess.Close()
		return
	}
	sess.Info = msg.GetData()
	sess.SessMgr.Register(sess.SessionBase)
	sess.Ctx.Log.Infoln("The service node registers with me, the node ID is ", utility.ServerID2UUID(msg.GetData().GetId()).String())
	sess.Ctx.Log.Infoln(sess.Info)
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
	if sess.SessionBase == sessbase && sessbase.Info != nil {
		sess.SessMgr.Lose1(sessbase)
		sess.Ctx.Log.Infoln("Service node loses connection, type:", sess.Info.GetType(), "id:", utility.ServerID2UUID(sess.Info.GetId()).String())
	}
}
