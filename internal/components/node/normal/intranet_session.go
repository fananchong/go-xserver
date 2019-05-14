package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utils"
	"github.com/fananchong/gotcp"
)

// IntranetSession : 网络会话类（ 服务器组内 Gateway 客户端会话类 ）
type IntranetSession struct {
	*nodecommon.SessionBase
	sourceSess *Session
}

// NewIntranetSession : 网络会话类的构造函数
func NewIntranetSession(ctx *common.Context, sessMgr *nodecommon.SessionMgr, sourceSess *Session) *IntranetSession {
	sess := &IntranetSession{}
	sess.SessionBase = nodecommon.NewSessionBase(ctx, sess)
	sess.SessMgr = sessMgr
	sess.sourceSess = sourceSess
	return sess
}

// Start : 启动
func (sess *IntranetSession) Start() {
	go func() {
		for {
			node := sess.SessMgr.GetByID(nodecommon.ServerID2NodeID(sess.Info.GetId()))
			if node == nil {
				// 目标节点已丢失，不用试图去连接啦
				break
			}
			address := fmt.Sprintf("%s:%d", sess.Info.GetAddrs()[utils.IPINNER], sess.Info.GetPorts()[utils.PORTFORINTRANET])
			sess.Ctx.Infoln("Try to connect to the gateway server, address:", address, "node:", sess.Info.GetId().GetID())
			if sess.Connect(address, sess) == false {
				time.Sleep(1 * time.Second)
				continue
			}
			sess.Verify()
			sess.RegisterSelf(sess.sourceSess.GetID(), sess.sourceSess.GetType(), config.Gateway)
			sess.Ctx.Infoln("Successfully connected to the gateway server, address:", address, "node:", sess.Info.GetId().GetID())
			break
		}
	}()
}

// DoRegister : 某节点注册时处理
func (sess *IntranetSession) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoVerify : 验证时保存自己的注册消息
func (sess *IntranetSession) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *IntranetSession) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {

}

// DoClose : 节点关闭时处理
func (sess *IntranetSession) DoClose(sessbase *nodecommon.SessionBase) {
}

// DoRecv : 节点收到消息处理
func (sess *IntranetSession) DoRecv(cmd uint64, data []byte, flag byte) (done bool) {
	done = true
	switch protocol.CMD_GW_ENUM(cmd) {
	case protocol.CMD_GW_RELAY_CLIENT_MSG:
		f := sess.FuncOnRelayMsg()
		if f == nil {
			sess.Ctx.Errorln("There is no handler function, the handler function is `FuncOnRelayMsg`")
			return
		}
		msg := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is `protocol.CMD_GW_RELAY_CLIENT_MSG`(", int(protocol.CMD_GW_RELAY_CLIENT_MSG), ")")
			return
		}
		f(config.Client, context.NodeID(0), msg.GetAccount(), uint64(msg.GetCMD()), msg.GetData())
	case protocol.CMD_GW_RELAY_SERVER_MSG1:
		f := sess.FuncOnRelayMsg()
		if f == nil {
			sess.Ctx.Errorln("There is no handler function, the handler function is `FuncOnRelayMsg`")
			return
		}
		msg := &protocol.MSG_GW_RELAY_SERVER_MSG1{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is `protocol.CMD_GW_RELAY_SERVER_MSG1`(", int(protocol.CMD_GW_RELAY_SERVER_MSG1), ")")
			return
		}
		if msg.GetTargetType() != uint32(sess.sourceSess.GetType()) {
			sess.Ctx.Errorln("Field 'TargetType' error. TargetType:", msg.GetTargetType(), "SessType:", sess.sourceSess.GetType())
			return
		}
		f(config.NodeType(msg.GetSourceType()), nodecommon.ServerID2NodeID(msg.GetSourceID()), "", uint64(msg.GetCMD()), msg.GetData())
	case protocol.CMD_GW_RELAY_SERVER_MSG2:
		f := sess.FuncOnRelayMsg()
		if f == nil {
			sess.Ctx.Errorln("There is no handler function, the handler function is `FuncOnRelayMsg`")
			return
		}
		msg := &protocol.MSG_GW_RELAY_SERVER_MSG2{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is `protocol.CMD_GW_RELAY_SERVER_MSG2`(", int(protocol.CMD_GW_RELAY_SERVER_MSG2), ")")
			return
		}
		if !nodecommon.EqualSID(msg.GetTargetID(), nodecommon.NodeID2ServerID(sess.sourceSess.GetID())) {
			sess.Ctx.Errorln("Field 'TargetID' error. TargetID:", msg.GetTargetID(), "SessID:", sess.sourceSess.GetID())
			return
		}
		f(config.NodeType(msg.GetSourceType()), nodecommon.ServerID2NodeID(msg.GetSourceID()), "", uint64(msg.GetCMD()), msg.GetData())
	case protocol.CMD_GW_REGISTER_ACCOUNT:
		msg := &protocol.MSG_GW_REGISTER_ACCOUNT{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is `protocol.CMD_GW_REGISTER_ACCOUNT`(", int(protocol.CMD_GW_REGISTER_ACCOUNT), ")")
			return
		}
		sess.sourceSess.GWMgr.AddUser(msg.Account, sess.SessionBase)
	case protocol.CMD_GW_LOSE_ACCOUNT:
		f := sess.FuncOnLoseAccount()
		if f == nil {
			sess.Ctx.Errorln("There is no handler function, the handler function is `FuncOnLoseAccount`")
			return
		}
		msg := &protocol.MSG_GW_LOSE_ACCOUNT{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Errorln("Message parsing failed, message number is `protocol.CMD_GW_LOSE_ACCOUNT`(", int(protocol.CMD_GW_LOSE_ACCOUNT), ")")
			return
		}
		f(msg.Account)
		sess.sourceSess.GWMgr.DelUser(msg.Account)
	default:
		done = false
	}
	return
}
