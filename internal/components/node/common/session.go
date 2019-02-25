package nodecommon

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
	"github.com/fananchong/gotcp"
)

// ISessionDerived : SessionBase 派生类接口定义
type ISessionDerived interface {
	DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte)
	DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte)
	DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte)
	DoClose(sessbase *SessionBase)
	DoRecv(cmd uint64, data []byte, flag byte) (done bool)
	DoSendClientMsgByRelay(account string, cmd uint64, data []byte) bool
}

// SessionBase : 网络会话类
type SessionBase struct {
	DefaultNodeInterfaceImpl
	*gotcp.Session
	Ctx     *common.Context
	MsgData []byte
	MsgFlag byte
	derived ISessionDerived
}

// NewSessionBase : 网络会话类的构造函数
func NewSessionBase(ctx *common.Context, derived ISessionDerived) *SessionBase {
	return &SessionBase{
		Ctx:     ctx,
		Session: &gotcp.Session{},
		derived: derived,
	}
}

// OnRecv : 接收到网络数据包，被触发
func (sessbase *SessionBase) OnRecv(data []byte, flag byte) {
	cmd := protocol.CMD_MGR_ENUM(gotcp.GetCmd(data))
	if sessbase.IsVerified() == false && sessbase.doVerify(cmd, data, flag) == false {
		return
	}
	switch cmd {
	case protocol.CMD_MGR_REGISTER_SERVER:
		sessbase.doRegister(data, flag)
	case protocol.CMD_MGR_LOSE_SERVER:
		sessbase.doLose(data, flag)
	case protocol.CMD_MGR_PING:
		// No need to do anything
	default:
		if sessbase.derived.DoRecv(uint64(cmd), data, flag) == false {
			sessbase.Ctx.Log.Errorln("Unknown message number, message number is", cmd)
		}
	}
}

// OnClose : 断开连接，被触发
func (sessbase *SessionBase) OnClose() {
	sessbase.derived.DoClose(sessbase)
}

func (sessbase *SessionBase) doVerify(cmd protocol.CMD_MGR_ENUM, data []byte, flag byte) bool {
	if cmd == protocol.CMD_MGR_REGISTER_SERVER {
		msg := &protocol.MSG_MGR_REGISTER_SERVER{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sessbase.Ctx.Log.Errorln("Message parsing failed, message number is`protocol.CMD_MGR_REGISTER_SERVER`(", int(protocol.CMD_MGR_REGISTER_SERVER), ")")
			sessbase.Close()
			return false
		}
		if msg.GetToken() != sessbase.Ctx.Config.Common.IntranetToken {
			sessbase.Ctx.Log.Errorln("Token verification failed.",
				"Msg token:", msg.GetToken(),
				"Expect token:", sessbase.Ctx.Config.Common.IntranetToken)
			sessbase.Close()
			return false
		}
		sessbase.derived.DoVerify(msg, data, flag)
		sessbase.Verify()
		return true
	}
	sessbase.Ctx.Log.Errorln("The expected message number is `protocol.CMD_MGR_REGISTER_SERVER`(", int(protocol.CMD_MGR_REGISTER_SERVER), "), but", cmd, "(", int(cmd), ")")
	sessbase.Close()
	return false
}

func (sessbase *SessionBase) doRegister(data []byte, flag byte) {
	msg := &protocol.MSG_MGR_REGISTER_SERVER{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		sessbase.Ctx.Log.Errorln("Message parsing failed, message number is`protocol.CMD_MGR_REGISTER_SERVER`(", int(protocol.CMD_MGR_REGISTER_SERVER), ")")
		sessbase.Close()
		return
	}
	sessbase.derived.DoRegister(msg, data, flag)
}

func (sessbase *SessionBase) doLose(data []byte, flag byte) {
	msg := &protocol.MSG_MGR_LOSE_SERVER{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		sessbase.Ctx.Log.Errorln("Message parsing failed, message number is`protocol.CMD_MGR_LOSE_SERVER`(", int(protocol.CMD_MGR_LOSE_SERVER), ")")
		return
	}
	sessbase.derived.DoLose(msg, data, flag)
}

// RegisterSelf : 注册自己
func (sessbase *SessionBase) RegisterSelf(targetServerType common.NodeType) {
	msg := &protocol.MSG_MGR_REGISTER_SERVER{}
	msg.Data = &protocol.SERVER_INFO{}
	msg.Data.Id = utility.NodeID2ServerID(sessbase.GetID())
	msg.Data.Type = uint32(sessbase.Ctx.Node.GetType())
	msg.Data.Addrs = []string{utils.GetIPInner(sessbase.Ctx), utils.GetIPOuter(sessbase.Ctx)}
	msg.Data.Ports = sessbase.Ctx.Config.Network.Port

	// TODO: 后续支持
	// msg.Data.Overload
	// msg.Data.Version

	msg.Token = sessbase.Ctx.Config.Common.IntranetToken
	msg.TargetServerType = uint32(targetServerType)
	sessbase.Info = msg.GetData()
	sessbase.SendMsg(uint64(protocol.CMD_MGR_REGISTER_SERVER), msg)

	if targetServerType == common.Mgr {
		sessbase.Ctx.Log.Infoln("Register your information with the management server, info:", msg.GetData())
	} else if targetServerType == common.Gateway {
		sessbase.Ctx.Log.Infoln("Register your information with the gateway server, info:", msg.GetData())
	} else {
		sessbase.Ctx.Log.Errorln("Register your information with the server(", targetServerType, "), info:", msg.GetData())
	}
}

// SendClientMsgByRelay : 发送消息给客户端，通过 Gateway 中继
func (sessbase *SessionBase) SendClientMsgByRelay(account string, cmd uint64, data []byte) bool {
	return sessbase.derived.DoSendClientMsgByRelay(account, cmd, data)
}
