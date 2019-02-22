package nodegateway

import (
	"context"
	"net"

	"github.com/fananchong/go-xserver/common"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
	"github.com/gogo/protobuf/proto"
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
		sess.Ctx.Log.Errorln("Service ID is different.")
		sess.Ctx.Log.Errorln("sess.Info.GetId() :", utility.ServerID2UUID(sess.Info.GetId()).String())
		sess.Ctx.Log.Errorln("msg.GetData().GetId() :", utility.ServerID2UUID(msg.GetData().GetId()).String())
		sess.Close()
		return
	}
	if msg.GetTargetServerType() != uint32(common.Gateway) {
		sess.Ctx.Log.Errorln("Target server type different. Expectation is common.Gateway, but it is ", msg.GetTargetServerType())
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

// DoRecv : 节点收到消息处理
func (sess *Session) DoRecv(cmd uint64, data []byte, flag byte) (done bool) {
	done = true
	switch protocol.CMD_GW_ENUM(cmd) {
	case protocol.CMD_GW_RELAY_CLIENT_MSG:
		msg := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Log.Errorln("Message parsing failed, message number is`protocol.CMD_GW_RELAY_CLIENT_MSG`(", int(protocol.CMD_GW_RELAY_CLIENT_MSG), ")")
			done = false
			return
		}
		sess.Ctx.Gateway.GetSendToClient()(msg.GetAccount(), uint64(msg.GetCMD()), msg.GetData())
	default:
		done = false
	}
	return
}

// DoSendClientMsgByRelay : 发送消息给客户端，通过 Gateway 中继
func (sess *Session) DoSendClientMsgByRelay(account string, cmd uint64, msg proto.Message) bool {
	panic("")
}
