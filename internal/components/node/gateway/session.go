package nodegateway

import (
	gocontext "context"
	"net"

	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/gotcp"
)

// Session : 网络会话类
type Session struct {
	*nodecommon.SessionBase
	funcSendToClient    context.FuncTypeSendToClient
	funcSendToAllClient context.FuncTypeSendToAllClient
}

// Init : 初始化网络会话节点
func (sess *Session) Init(root gocontext.Context, conn net.Conn, derived gotcp.ISession, userdata interface{}) {
	ud := userdata.(*nodecommon.UserData)
	sess.SessionBase = nodecommon.NewSessionBase(ud.Ctx, sess)
	sess.SessionBase.Init(root, conn, derived)
	sess.SessMgr = ud.SessMgr
	sess.funcSendToClient = ud.Ctx.IGateway.(*Gateway).GetSendToClient()
	sess.funcSendToAllClient = ud.Ctx.IGateway.(*Gateway).GetSendToAllClient()
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Info = msg.GetData()
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	if nodecommon.EqualSID(sess.Info.GetId(), msg.GetData().GetId()) == false {
		sess.Ctx.Errorln("Service ID is different.")
		sess.Ctx.Errorln("sess.Info.GetId() :", sess.Info.GetId().GetID())
		sess.Ctx.Errorln("msg.GetData().GetId() :", msg.GetData().GetId().GetID())
		sess.Close()
		return
	}
	if msg.GetTargetServerType() != uint32(config.Gateway) {
		sess.Ctx.Errorln("Target server type different. Expectation is config.Gateway, but it is", msg.GetTargetServerType())
		sess.Close()
		return
	}
	sess.Info = msg.GetData()
	sess.SessMgr.Register(sess.SessionBase)
	sess.Ctx.Infoln("The service node registers with me, the node ID is", msg.GetData().GetId().GetID())
	sess.Ctx.Infoln(sess.Info)
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
	if sess.SessionBase == sessbase && sessbase.Info != nil {
		sess.SessMgr.Lose1(sessbase)
		sess.Ctx.Infoln("Service node loses connection, type:", sess.Info.GetType(), "id:", sess.Info.GetId().GetID())
	}
}

// DoRecv : 节点收到消息处理
func (sess *Session) DoRecv(cmd uint64, data []byte, flag byte) (done bool) {
	done = true
	switch protocol.CMD_GW_ENUM(cmd) {
	case protocol.CMD_GW_RELAY_CLIENT_MSG:
		msg := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_GW_RELAY_CLIENT_MSG`(", int(protocol.CMD_GW_RELAY_CLIENT_MSG), ")")
			done = false
			return
		}
		targetAccount := msg.GetAccount()
		if targetAccount != "" {
			sess.funcSendToClient(targetAccount,
				uint64(msg.GetCMD())+uint64(sess.Info.GetType())*uint64(sess.Ctx.Config().Common.MsgCmdOffset),
				msg.GetData(), uint8(msg.GetFlag()))
		} else {
			sess.funcSendToAllClient(
				uint64(msg.GetCMD())+uint64(sess.Info.GetType())*uint64(sess.Ctx.Config().Common.MsgCmdOffset),
				msg.GetData(), uint8(msg.GetFlag()))
		}
	case protocol.CMD_GW_RELAY_SERVER_MSG1:
		msg := &protocol.MSG_GW_RELAY_SERVER_MSG1{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_GW_RELAY_SERVER_MSG1`(", int(protocol.CMD_GW_RELAY_SERVER_MSG1), ")")
			done = false
			return
		}
		if msg.GetSendType() == protocol.RELAY_SERVER_MSG_TYPE_BROADCAST {
			sess.SessMgr.ForByType(config.NodeType(msg.GetTargetType()), func(targetSess *nodecommon.SessionBase) {
				targetSess.Send(data, flag)
			})
		} else if msg.GetSendType() == protocol.RELAY_SERVER_MSG_TYPE_RANDOM {
			targetSess := sess.SessMgr.SelectOne(config.NodeType(msg.GetTargetType()))
			if targetSess != nil {
				targetSess.Send(data, flag)
			} else {
				sess.Ctx.Errorln("No find server.", "cmd:", msg.GetCMD(), "targetType:", msg.GetTargetType())
				done = false
				return
			}
		} else {
			sess.Ctx.Errorln("Field 'send type' error, value:", msg.GetSendType(), "cmd:", msg.GetCMD(), "targetType:", msg.GetTargetType())
			done = false
			return
		}
	case protocol.CMD_GW_RELAY_SERVER_MSG2:
		msg := &protocol.MSG_GW_RELAY_SERVER_MSG2{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_GW_RELAY_SERVER_MSG2`(", int(protocol.CMD_GW_RELAY_SERVER_MSG2), ")")
			done = false
			return
		}
		id := nodecommon.ServerID2NodeID(msg.GetTargetID())
		targetSess := sess.SessMgr.GetByID(id)
		if targetSess != nil {
			targetSess.Send(data, flag)
		} else {
			sess.Ctx.Errorln("No find server.", "cmd:", msg.GetCMD(), "targetServerID:", id)
			done = false
			return
		}
	default:
		done = false
	}
	return
}
