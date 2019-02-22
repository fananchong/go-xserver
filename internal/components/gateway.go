package components

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	nodegateway "github.com/fananchong/go-xserver/internal/components/node/gateway"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
)

// Gateway : 网关服务器
type Gateway struct {
	*nodegateway.Node
	ctx              *common.Context
	funcSendToClient common.FuncTypeSendToClient
	funcEncodeFunc   common.FuncTypeEncode
	funcDecodeFunc   common.FuncTypeDecode
}

// NewGateway : 构造函数
func NewGateway(ctx *common.Context) *Gateway {
	gw := &Gateway{
		ctx: ctx,
	}
	gw.Node = nodegateway.NewNode(ctx)
	gw.ctx.Gateway = gw
	return gw
}

// Start : 启动
func (gw *Gateway) Start() bool {
	if getPluginType(gw.ctx) == common.Gateway {
		if gw.initRedis() == false {
			return false
		}
		if gw.Node.Init() == false {
			return false
		}
		if gw.Node.Start() == false {
			return false
		}
	}
	return true
}

// Close : 关闭
func (gw *Gateway) Close() {
	if gw.Node != nil {
		gw.Node.Close()
		gw.Node = nil
	}
}

// VerifyToken : 令牌验证。返回值： 0 成功；1 令牌错误； 2 系统错误
func (gw *Gateway) VerifyToken(account, token string) uint32 {
	tokenObj := db.NewToken(gw.ctx.Config.DbToken.Name, account)
	if err := tokenObj.Load(); err != nil {
		gw.ctx.Log.Errorln(err, "account:", account)
		return 2
	}
	if token != tokenObj.GetToken() {
		gw.ctx.Log.Errorf("Token verification failed, expecting token to be %s, but %s. account: %s\n", tokenObj.GetToken(), token, account)
		return 1
	}
	return 0
}

// OnRecvFromClient : 可自定义客户端交互协议。data 格式需转化为框架层可理解的格式。done 为 true ，表示框架层接管处理该消息
func (gw *Gateway) OnRecvFromClient(account string, cmd uint32, data []byte) (done bool) {
	nodeType := common.NodeType(cmd / uint32(gw.ctx.Config.Common.MsgCmdOffset))
	if nodeType >= common.NodeTypeSize || nodeType <= common.Gateway {
		gw.ctx.Log.Errorln("Wrong message number. cmd:", cmd, "account:", account)
		return
	}

	// TODO: 先调通消息，实际上要根据该账号是否有分配对应的目标类型服务，来决定是定向中继（即状态中继）、还是随机中继
	//       调通消息后，会出文档，并在这里实现。暂随机中继
	target := gw.GetNodeOne(nodeType)
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
	return true
}
