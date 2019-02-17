package nodenormal

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/fananchong/go-xserver/common"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

// Session : 网络会话类
type Session struct {
	*nodecommon.SessionBase
	shutdown int32
}

// NewSession : 网络会话类的构造函数
func NewSession(ctx *common.Context) *Session {
	sess := &Session{}
	sess.SessionBase = nodecommon.NewSessionBase(ctx, sess)
	sess.SessMgr = nodecommon.NewSessionMgr()
	return sess
}

// Start : 开始连接 Mgr Server
func (sess *Session) Start() bool {
	atomic.StoreInt32(&sess.shutdown, 0)
	sess.connectMgrServer()
	return true
}

func (sess *Session) connectMgrServer() {
	if shutdown := atomic.LoadInt32(&sess.shutdown); shutdown == 0 {
	TRY_AGAIN:
		addr, port := getMgrInfoByBlock(sess.Ctx)
		if sess.Connect(fmt.Sprintf("%s:%d", addr, port), sess) == false {
			time.Sleep(1 * time.Second)
			goto TRY_AGAIN
		}
		sess.Verify()
		sess.RegisterSelf(common.Mgr)
	}
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Ctx.Log.Infoln("The service node registers information with me with ID", utility.ServerID2UUID(msg.GetData().GetId()).String())

	tempSess := sess.SessMgr.GetByID(utility.ServerID2NodeID(msg.GetData().GetId()))
	if tempSess == nil {
		// 本地保其他存节点信息
		targetSess := NewIntranetSession(sess.Ctx, sess.SessMgr)
		targetSess.Info = msg.GetData()
		targetSess.RegisterFuncOnRelayMsg(sess.FuncOnRelayMsg())
		sess.SessMgr.Register(targetSess.SessionBase)

		// 如果存在互连关系的，开始互连逻辑。
		if sess.IsEnableMessageRelay() && targetSess.Info.GetType() == uint32(common.Gateway) {
			targetSess.Start()
		}
	} else {
		tempSess.Info = msg.GetData()
	}
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
	sess.Ctx.Log.Infoln("Service node connection lost, ID is", utility.ServerID2UUID(msg.GetId()).String(), "type:", msg.GetType())

	// 如果存在互连关系的，关闭 TCP 连接
	if sess.IsEnableMessageRelay() && msg.GetType() == uint32(common.Gateway) {
		targetSess := sess.SessMgr.GetByID(utility.ServerID2NodeID(msg.GetId()))
		if targetSess != nil {
			targetSess.Close()
		}
	}

	sess.SessMgr.Lose2(msg.GetId(), common.NodeType(msg.GetType()))

	sess.Ctx.Log.Infof("Remaining list of service nodes of this type[%d]:\n", msg.GetType())
	sess.SessMgr.ForByType(common.NodeType(msg.GetType()), func(sessbase *nodecommon.SessionBase) {
		sess.Ctx.Log.Infof("\t%s\n", utility.ServerID2UUID(sessbase.GetSID()).String())
	})
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
	if sessbase.Info != nil {
		sess.SessMgr.Lose1(sessbase)
	}
	go func() {
		time.Sleep(1 * time.Second)
		sess.connectMgrServer()
	}()
}

// Ping : ping
func (sess *Session) Ping() {
	msg := &protocol.MSG_MGR_PING{}
	sess.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
}

func getMgrInfoByBlock(ctx *common.Context) (string, int32) {
	ctx.Log.Infoln("Try to get management server information ...")
	data := db.NewMgrServer(ctx.Config.DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			ctx.Log.Errorln(err)
			time.Sleep(1 * time.Second)
		}
	}
	ctx.Log.Infoln("The address of the management server is", data.GetAddr())
	ctx.Log.Infoln("The port of the management server is", data.GetPort())
	return data.GetAddr(), data.GetPort()
}

// Shutdown : 关闭
func (sess *Session) Shutdown() {
	atomic.StoreInt32(&sess.shutdown, 1)
}

// DoRecv : 节点收到消息处理
func (sess *Session) DoRecv(cmd uint64, data []byte, flag byte) (done bool) {
	return
}
