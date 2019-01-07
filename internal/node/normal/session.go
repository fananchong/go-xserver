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
	sess.ResetTCPSession()
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
	//nodecommon.XSESSIONMGR.Register()
	sess.Ctx.Log.Infoln("one server register. id:", utility.ServerID2UUID(msg.GetData().GetId()).String())
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
	//nodecommon.XSESSIONMGR.Lose()
	sess.Ctx.Log.Infoln("one server lose. id:", utility.ServerID2UUID(msg.GetId()).String())
}

// DoClose : 节点关闭时处理
func (sess *Session) DoClose(sessbase *nodecommon.SessionBase) {
	go func() {
		time.Sleep(1 * time.Second)
		sess.connectMgrServer()
	}()
}

func getMgrInfoByBlock(ctx *common.Context) (string, int32) {
	ctx.Log.Infoln("Try get mgr server info ...")
	data := db.NewMgrServer(ctx.Config.DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			ctx.Log.Debugln(err)
			time.Sleep(1 * time.Second)
		}
	}
	ctx.Log.Infoln("Mgr server address:", data.GetAddr())
	ctx.Log.Infoln("Mgr server port:", data.GetPort())
	return data.GetAddr(), data.GetPort()
}
