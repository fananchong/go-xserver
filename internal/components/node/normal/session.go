package nodenormal

import (
	"fmt"
	"sync/atomic"
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/gomodule/redigo/redis"
)

// Session : 网络会话类
type Session struct {
	*nodecommon.SessionBase
	GWMgr       *IntranetSessionMgr
	shutdown    int32
	ReqServerID *protocol.SERVER_ID
}

// NewSession : 网络会话类的构造函数
func NewSession(ctx *common.Context) *Session {
	sess := &Session{}
	sess.SessionBase = nodecommon.NewSessionBase(ctx, sess)
	sess.SessMgr = nodecommon.NewSessionMgr(ctx)
	sess.GWMgr = NewIntranetSessionMgr(ctx)
	return sess
}

// Start : 开始连接 Mgr Server
func (sess *Session) Start() bool {
	sess.GWMgr.Start()
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
		sess.RegisterSelf(sess.GetID(), sess.GetType(), config.Mgr)
	}
}

// DoRegister : 某节点注册时处理
func (sess *Session) DoRegister(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
	sess.Ctx.Infoln("The service node registers information with me with ID", msg.GetData().GetId().GetID(), "type:", msg.GetData().GetType())

	tempSess := sess.SessMgr.GetByID(nodecommon.ServerID2NodeID(msg.GetData().GetId()))
	if tempSess == nil {
		// 本地保其他存节点信息
		targetSess := NewIntranetSession(sess.Ctx, sess.SessMgr, sess)
		targetSess.Info = msg.GetData()
		targetSess.RegisterFuncOnRelayMsg(sess.FuncOnRelayMsg())
		targetSess.RegisterFuncOnLoseAccount(sess.FuncOnLoseAccount())
		sess.SessMgr.Register(targetSess.SessionBase)

		sess.PrintNodeInfo(sess.Ctx, config.NodeType(msg.GetData().GetType()))

		// 如果存在互连关系的，开始互连逻辑。
		if sess.IsEnableMessageRelay() && targetSess.Info.GetType() == uint32(config.Gateway) {
			targetSess.Start()
		}
	} else {
		tempSess.Info = msg.GetData()
		sess.PrintNodeInfo(sess.Ctx, config.NodeType(msg.GetData().GetType()))
	}
}

// DoVerify : 验证时保存自己的注册消息
func (sess *Session) DoVerify(msg *protocol.MSG_MGR_REGISTER_SERVER, data []byte, flag byte) {
}

// DoLose : 节点丢失时处理
func (sess *Session) DoLose(msg *protocol.MSG_MGR_LOSE_SERVER, data []byte, flag byte) {
	sess.Ctx.Infoln("Service node connection lost, ID is", msg.GetId().GetID(), "type:", msg.GetType())

	// 如果存在互连关系的，关闭 TCP 连接
	if sess.IsEnableMessageRelay() && msg.GetType() == uint32(config.Gateway) {
		targetSess := sess.SessMgr.GetByID(nodecommon.ServerID2NodeID(msg.GetId()))
		if targetSess != nil {
			targetSess.Close()
		}
	}

	sess.SessMgr.Lose2(msg.GetId(), config.NodeType(msg.GetType()))
	sess.PrintNodeInfo(sess.Ctx, config.NodeType(msg.GetType()))
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
	ctx.Infoln("Try to get management server information ...")
	data := db.NewMgrServer(ctx.Config().DbMgr.Name, 0)
	for {
		if err := data.Load(); err == nil {
			break
		} else {
			ctx.Errorln("Get management server information fail. error:", err.Error())
			time.Sleep(1 * time.Second)
		}
	}
	ctx.Infoln("The address of the management server is", data.GetAddr())
	ctx.Infoln("The port of the management server is", data.GetPort())
	return data.GetAddr(), data.GetPort()
}

// Shutdown : 关闭
func (sess *Session) Shutdown() {
	if sess.GWMgr != nil {
		sess.GWMgr.Close()
		sess.GWMgr = nil
	}
	atomic.StoreInt32(&sess.shutdown, 1)
}

// DoRecv : 节点收到消息处理
func (sess *Session) DoRecv(cmd uint64, data []byte, flag byte) (done bool) {
	return
}

// SendMsgToClient : 发送消息给客户端，通过 Gateway 中继
func (sess *Session) SendMsgToClient(account string, cmd uint64, data []byte) bool {
	targetSess := sess.GWMgr.GetAndActive(account)
	if targetSess != nil {
		msgRelay := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		msgRelay.Account = account
		msgRelay.CMD = uint32(cmd)
		msgRelay.Data = data
		return targetSess.SendMsg(uint64(protocol.CMD_GW_RELAY_CLIENT_MSG), msgRelay)
	}
	// 非本服务节点上的账号，则查找对应的 Gateway ID ，再发送
	dbaccess := go_redis_orm.GetDB(sess.Ctx.Config().DbServer.Name)
	key := db.GetKeyAllocServer(uint32(config.Gateway), account)
	val, err := redis.String(dbaccess.Do("GET", key))
	if err != nil {
		sess.Ctx.Errorln(err, "account:", account, ", cmd:", cmd)
		return false
	}
	dbObj := &db.AccountServer{}
	if err := dbObj.Unmarshal(val); err != nil {
		sess.Ctx.Errorln(err, "account:", account, ", cmd:", cmd)
		return false
	}
	ttl, err := redis.Uint64(dbaccess.Do("TTL", key))
	if err != nil {
		sess.Ctx.Errorln(err, "account:", account, ", cmd:", cmd)
		return false
	}
	if int64(ttl) <= sess.Ctx.Config().Role.SessionAffinityInterval {
		sess.Ctx.Infoln("Target account offline", "account:", account, ", cmd:", cmd)
		return false
	}
	targetSess = sess.SessMgr.GetByID(nodecommon.ServerID2NodeID(dbObj.ServerID))
	sess.GWMgr.AddUser(account, targetSess) // sess 加入缓存
	msgRelay := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
	msgRelay.Account = account
	msgRelay.CMD = uint32(cmd)
	msgRelay.Data = data
	return targetSess.SendMsg(uint64(protocol.CMD_GW_RELAY_CLIENT_MSG), msgRelay)
}

// BroadcastMsgToClient : 广播消息给客户端，通过 Gateway 中继
func (sess *Session) BroadcastMsgToClient(cmd uint64, data []byte) bool {
	sess.SessMgr.ForByType(config.Gateway, func(targetSess *nodecommon.SessionBase) {
		msgRelay := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
		msgRelay.Account = ""
		msgRelay.CMD = uint32(cmd)
		msgRelay.Data = data
		targetSess.SendMsg(uint64(protocol.CMD_GW_RELAY_CLIENT_MSG), msgRelay)
	})
	return true
}

// SendMsgToServer : 发送消息给某类型服务（随机一个）
func (sess *Session) SendMsgToServer(t config.NodeType, cmd uint64, data []byte) bool {
	gw := sess.SessMgr.SelectOne(config.Gateway)
	if gw != nil {
		msgRelay := &protocol.MSG_GW_RELAY_SERVER_MSG1{}
		msgRelay.SourceID = nodecommon.NodeID2ServerID(sess.GetID())
		msgRelay.SourceType = uint32(sess.GetType())
		msgRelay.TargetType = uint32(t)
		msgRelay.SendType = protocol.RELAY_SERVER_MSG_TYPE_RANDOM
		msgRelay.CMD = uint32(cmd)
		msgRelay.Data = data
		return gw.SendMsg(uint64(protocol.CMD_GW_RELAY_SERVER_MSG1), msgRelay)
	}
	return false
}

// ReplyMsgToServer : 回发消息给请求服务器
func (sess *Session) ReplyMsgToServer(cmd uint64, data []byte) bool {
	if sess.ReqServerID == nil {
		return false
	}
	gw := sess.SessMgr.SelectOne(config.Gateway)
	if gw != nil {
		msgRelay := &protocol.MSG_GW_RELAY_SERVER_MSG2{}
		msgRelay.SourceID = nodecommon.NodeID2ServerID(sess.GetID())
		msgRelay.SourceType = uint32(sess.GetType())
		msgRelay.TargetID = sess.ReqServerID
		msgRelay.CMD = uint32(cmd)
		msgRelay.Data = data
		return gw.SendMsg(uint64(protocol.CMD_GW_RELAY_SERVER_MSG2), msgRelay)
	}
	return false
}

// BroadcastMsgToServer : 广播消息给某类型服务
func (sess *Session) BroadcastMsgToServer(t config.NodeType, cmd uint64, data []byte) bool {
	gw := sess.SessMgr.SelectOne(config.Gateway)
	if gw != nil {
		msgRelay := &protocol.MSG_GW_RELAY_SERVER_MSG1{}
		msgRelay.SourceID = nodecommon.NodeID2ServerID(sess.GetID())
		msgRelay.SourceType = uint32(sess.GetType())
		msgRelay.TargetType = uint32(t)
		msgRelay.SendType = protocol.RELAY_SERVER_MSG_TYPE_BROADCAST
		msgRelay.CMD = uint32(cmd)
		msgRelay.Data = data
		return gw.SendMsg(uint64(protocol.CMD_GW_RELAY_SERVER_MSG1), msgRelay)
	}
	return false
}
