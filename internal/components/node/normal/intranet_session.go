package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/go-xserver/common"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
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
			node := sess.SessMgr.GetByID(utility.ServerID2NodeID(sess.Info.GetId()))
			if node == nil {
				// 目标节点已丢失，不用试图去连接啦
				break
			}
			address := fmt.Sprintf("%s:%d", sess.Info.GetAddrs()[common.IPINNER], sess.Info.GetPorts()[common.PORTFORINTRANET])
			sess.Ctx.Log.Infoln("Try to connect to the gateway server, address:", address, "node:", utility.ServerID2UUID(sess.Info.GetId()).String())
			if sess.Connect(address, sess) == false {
				time.Sleep(1 * time.Second)
				continue
			}
			sess.Verify()
			sess.RegisterSelf(sess.sourceSess.GetID(), common.Gateway)
			sess.Ctx.Log.Infoln("Successfully connected to the gateway server, address:", address, "node:", utility.ServerID2UUID(sess.Info.GetId()).String())
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
			sess.Ctx.Log.Errorln("There is no handler function, the handler function is `FuncOnRelayMsg`")
			return
		}
		msg := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Log.Errorln("Message parsing failed, message number is `protocol.CMD_GW_RELAY_CLIENT_MSG`(", int(protocol.CMD_GW_RELAY_CLIENT_MSG), ")")
			return
		}
		f(common.Client, sess, msg.GetAccount(), uint64(msg.GetCMD()), msg.GetData())
	case protocol.CMD_GW_LOSE_ACCOUNT:
		f := sess.FuncOnLoseAccount()
		if f == nil {
			sess.Ctx.Log.Errorln("There is no handler function, the handler function is `FuncOnLoseAccount`")
			return
		}
		msg := &protocol.MSG_GW_LOSE_ACCOUNT{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			sess.Ctx.Log.Errorln("Message parsing failed, message number is `protocol.CMD_GW_LOSE_ACCOUNT`(", int(protocol.CMD_GW_LOSE_ACCOUNT), ")")
			return
		}
		f(msg.Account)
	default:
		done = false
	}
	return
}

// DoSendClientMsgByRelay : 发送消息给客户端，通过 Gateway 中继
func (sess *IntranetSession) DoSendClientMsgByRelay(account string, cmd uint64, data []byte) bool {
	msgRelay := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
	msgRelay.Account = account
	msgRelay.CMD = uint32(cmd)
	msgRelay.Data = data
	return sess.SendMsg(uint64(protocol.CMD_GW_RELAY_CLIENT_MSG), msgRelay)
}
