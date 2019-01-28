package nodenormal

import (
	"fmt"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/db"
	nodecommon "github.com/fananchong/go-xserver/internal/node/common"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

// Session : 网络会话类
type Session struct {
	*nodecommon.SessionBase
}

// NewSession : 网络会话类的构造函数
func NewSession(ctx *common.Context) *Session {
	sess := &Session{}
	sess.SessionBase = nodecommon.NewSessionBase(ctx, sess)
	return sess
}

// Start : 开始连接 Mgr Server
func (sess *Session) Start() bool {
	sess.connectMgrServer()
	return true
}

func (sess *Session) connectMgrServer() {
TRY_AGAIN:
	addr, port := getMgrInfoByBlock(sess.Ctx)
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
	msg.Data.Id = utility.NodeID2ServerID(sess.GetID())
	msg.Data.Type = uint32(sess.Ctx.Node.GetType())
	msg.Data.Addrs = []string{utils.GetIPInner(sess.Ctx), utils.GetIPOuter(sess.Ctx)}
	msg.Data.Ports = sess.Ctx.Config.Network.Port

	// TODO: 后续支持
	// msg.Data.Overload
	// msg.Data.Version

	msg.Token = sess.Ctx.Config.Common.IntranetToken
	sess.Info = msg.GetData()
	sess.SendMsg(uint64(protocol.CMD_MGR_REGISTER_SERVER), msg)
	sess.Ctx.Log.Infoln("register self to mgr server. self info:", msg.GetData())
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Ctx.Log.Infoln("one server register. id:", utility.ServerID2UUID(msg.GetData().GetId()).String())

	// 本地保其他存节点信息
	targetSess := NewIntranetSession(sess.Ctx)
	targetSess.Info = msg.GetData()
	nodecommon.GetSessionMgr().Register(targetSess.SessionBase)

	// 如果存在互连关系的，开始互连逻辑。
	if sess.IsEnableMessageRelay() && targetSess.Info.GetType() == uint32(common.Gateway) {
		targetSess.Start()
	}
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
	sess.Ctx.Log.Infoln("one server lose. id:", utility.ServerID2UUID(msg.GetId()).String(), "type:", msg.GetType())

	// 如果存在互连关系的，关闭 TCP 连接
	if sess.IsEnableMessageRelay() && msg.GetType() == uint32(common.Gateway) {
		targetSess := nodecommon.GetSessionMgr().GetByID(utility.ServerID2NodeID(msg.GetId()))
		if targetSess != nil {
			targetSess.Close()
		}
	}

	nodecommon.GetSessionMgr().Lose2(msg.GetId(), common.NodeType(msg.GetType()))

	sess.Ctx.Log.Infof("left node in type[%d]:\n", msg.GetType())
	nodecommon.GetSessionMgr().ForByType(common.NodeType(msg.GetType()), func(sessbase *nodecommon.SessionBase) {
		sess.Ctx.Log.Infof("\t%s\n", utility.ServerID2UUID(sessbase.GetSID()).String())
	})
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
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
	ctx.Log.Infoln("Try get mgr server info ...")
	data := db.NewMgrServer(ctx.Config.DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			ctx.Log.Errorln(err)
			time.Sleep(1 * time.Second)
		}
	}
	ctx.Log.Infoln("Mgr server address:", data.GetAddr())
	ctx.Log.Infoln("Mgr server port:", data.GetPort())
	return data.GetAddr(), data.GetPort()
}
