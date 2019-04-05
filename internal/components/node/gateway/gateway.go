package nodegateway

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
)

// Gateway : 网关节点
type Gateway struct {
	*nodecommon.Node
	ctx                 *common.Context
	funcSendToClient    common.FuncTypeSendToClient
	funcSendToAllClient common.FuncTypeSendToAllClient
	funcEncodeFunc      common.FuncTypeEncode
	funcDecodeFunc      common.FuncTypeDecode
	users               *UserMgr
}

// NewGateway : 网关节点实现类的构造函数
func NewGateway(ctx *common.Context) *Gateway {
	gw := &Gateway{
		ctx:  ctx,
		Node: nodecommon.NewNode(ctx, common.Gateway),
	}
	gw.ctx.Gateway = gw
	gw.users = NewUserMgr(ctx, gw)
	return gw
}

// Start : 启动
func (gw *Gateway) Start() bool {
	if misc.GetPluginType(gw.ctx.Ctx) == common.Gateway {
		if gw.initRedis() == false {
			return false
		}
		if gw.Node.Init(Session{}, []utils.IComponent{}) == false {
			return false
		}
		if gw.Node.Start() == false {
			return false
		}
		gw.users.Start()
	}
	return true
}

// Close : 关闭
func (gw *Gateway) Close() {
	if gw.Node != nil {
		gw.Node.Close()
		gw.users.Close()
		gw.Node = nil
	}
}

// VerifyToken : 令牌验证。返回值： 0 成功；1 令牌错误； 2 系统错误
func (gw *Gateway) VerifyToken(account, token string, clientSession common.IClientSesion) uint32 {
	tokenObj := db.NewToken(gw.ctx.Config.DbToken.Name, account)
	if err := tokenObj.Load(); err != nil {
		gw.ctx.Log.Errorln(err, "account:", account)
		return 2
	}
	tmpTokenObj := tokenObj.GetToken(false)
	if token != tmpTokenObj.Token {
		gw.ctx.Log.Errorf("Token verification failed, expecting token to be %s, but %s. account: %s\n", tmpTokenObj.Token, token, account)
		return 1
	}
	gw.users.AddUser(account, tmpTokenObj.GetAllocServers(), clientSession)
	return 0
}

// OnRecvFromClient : 可自定义客户端交互协议。data 格式需转化为框架层可理解的格式。done 为 true ，表示框架层接管处理该消息
func (gw *Gateway) OnRecvFromClient(account string, cmd uint32, data []byte) (done bool) {
	nodeType := common.NodeType(cmd / uint32(gw.ctx.Config.Common.MsgCmdOffset))
	if nodeType >= common.NodeTypeSize || nodeType <= common.Gateway {
		gw.ctx.Log.Errorln("Wrong message number. cmd:", cmd, "account:", account)
		return
	}

	// 是否需要状态中继
	nodeID, err := gw.users.GetServerAndActive(account, nodeType)
	if err != nil {
		gw.ctx.Log.Errorln(err, "account:", account, "cmd:", cmd)
		return
	}
	var target *nodecommon.SessionBase
	if nodeID != nil {
		target = gw.GetNode(*nodeID)
	} else {
		target = gw.GetNodeOne(nodeType)
	}
	if target == nil {
		gw.ctx.Log.Errorln("Target server not found. cmd:", cmd, "account:", account, "nodeType", nodeType)
		return
	}

	// Gateway 接管该消息，并开始中继
	done = true

	msg := &protocol.MSG_GW_RELAY_CLIENT_MSG{}
	msg.Account = account
	msg.CMD = cmd % uint32(gw.ctx.Config.Common.MsgCmdOffset)
	msg.Data = append(msg.Data, data...)
	if target.SendMsg(uint64(protocol.CMD_GW_RELAY_CLIENT_MSG), msg) == false {
		gw.ctx.Log.Errorln("Sending a message to the target server failed. cmd:", cmd, "account:", account, "nodeType", nodeType)
		return
	}
	return
}

// RegisterSendToClient : 可自定义客户端交互协议
func (gw *Gateway) RegisterSendToClient(f common.FuncTypeSendToClient) {
	gw.funcSendToClient = f
}

// GetSendToClient : 可自定义客户端交互协议
func (gw *Gateway) GetSendToClient() common.FuncTypeSendToClient {
	return gw.funcSendToClient
}

// RegisterSendToAllClient : 可自定义客户端交互协议
func (gw *Gateway) RegisterSendToAllClient(f common.FuncTypeSendToAllClient) {
	gw.funcSendToAllClient = f
}

// GetSendToAllClient : 可自定义客户端交互协议
func (gw *Gateway) GetSendToAllClient() common.FuncTypeSendToAllClient {
	return gw.funcSendToAllClient
}

// RegisterEncodeFunc : 可自定义加解密算法
func (gw *Gateway) RegisterEncodeFunc(f common.FuncTypeEncode) {
	gw.funcEncodeFunc = f
}

// RegisterDecodeFunc : 可自定义加解密算法
func (gw *Gateway) RegisterDecodeFunc(f common.FuncTypeDecode) {
	gw.funcDecodeFunc = f
}

func (gw *Gateway) initRedis() bool {
	// db token
	err := go_redis_orm.CreateDB(
		gw.ctx.Config.DbToken.Name,
		gw.ctx.Config.DbToken.Addrs,
		gw.ctx.Config.DbToken.Password,
		gw.ctx.Config.DbToken.DBIndex)
	if err != nil {
		gw.ctx.Log.Errorln(err)
		return false
	}

	// db server
	err = go_redis_orm.CreateDB(
		gw.ctx.Config.DbServer.Name,
		gw.ctx.Config.DbServer.Addrs,
		gw.ctx.Config.DbServer.Password,
		gw.ctx.Config.DbServer.DBIndex)
	if err != nil {
		gw.ctx.Log.Errorln(err)
		return false
	}
	gw.users.ServerRedisCli = go_redis_orm.GetDB(gw.ctx.Config.DbServer.Name)
	return true
}
