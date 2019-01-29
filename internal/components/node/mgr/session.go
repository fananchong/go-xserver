package nodemgr

import (
	"context"
	"net"

	"github.com/fananchong/go-xserver/common"
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
	sess.SessionBase = nodecommon.NewSessionBase(userdata.(*common.Context), sess)
	sess.SessionBase.Init(root, conn, derived)
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Info = msg.GetData()
	sess.MsgData = make([]byte, len(data))
	copy(sess.MsgData, data)
	sess.MsgFlag = flag
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	if utility.EqualSID(sess.Info.GetId(), msg.GetData().GetId()) == false {
		sess.Close()
		return
	}
	sess.Info = msg.GetData()
	sess.MsgData = make([]byte, len(data))
	copy(sess.MsgData, data)
	sess.MsgFlag = flag
	sess.Ctx.Log.Infoln("Node register for me, node id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
	sess.Ctx.Log.Infoln(sess.Info)

	nodecommon.GetSessionMgr().Register(sess.SessionBase)
	nodecommon.GetSessionMgr().ForAll(func(elem *nodecommon.SessionBase) {
		if elem != sess.SessionBase {
			elem.Send(data, flag)
		}
	})
	nodecommon.GetSessionMgr().ForAll(func(elem *nodecommon.SessionBase) {
		if elem != sess.SessionBase {
			sess.Send(elem.MsgData, elem.MsgFlag)
		}
	})
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
	if sess.SessionBase == sessbase && sessbase.Info != nil {
		msg := &protocol.MSG_MGR_LOSE_SERVER{}
		msg.Id = sess.Info.GetId()
		msg.Type = sess.Info.GetType()
		nodecommon.GetSessionMgr().ForAll(func(elem *nodecommon.SessionBase) {
			elem.SendMsg(uint64(protocol.CMD_MGR_LOSE_SERVER), msg)
		})
		sess.Ctx.Log.Infoln("lose node, type:", msg.Type, "id:", utility.ServerID2UUID(msg.Id).String())
	}
}
