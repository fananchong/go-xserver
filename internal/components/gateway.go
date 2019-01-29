package components

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/components/gateway"
	"github.com/fananchong/gotcp"
)

// Gateway : 网关服务器
type Gateway struct {
	ctx *common.Context
}

// NewGateway : 构造函数
func NewGateway(ctx *common.Context) *Gateway {
	gw := &Gateway{
		ctx: ctx,
	}
	gw.ctx.Gateway = gw
	return gw
}

// Start : 启动
func (gw *Gateway) Start() bool {
	if getPluginType(gw.ctx) == common.Gateway {
		gw.ctx.ServerForIntranet.(*gotcp.Server).SetUserData(gw.ctx)
		gw.ctx.ServerForIntranet.RegisterSessType(gateway.IntranetSession{})
	}
	return true
}

// Close : 关闭
func (gw *Gateway) Close() {

}

// OnRecvFromClient : 可自定义客户端交互协议。data 格式需转化为框架层可理解的格式。done 为 true ，表示框架层接管处理该消息
func (gw *Gateway) OnRecvFromClient(data []byte) (done bool) {
	return
}

// RegisterSendToClient : 可自定义客户端交互协议
func (gw *Gateway) RegisterSendToClient(f common.FuncTypeSendToClient) {

}

// RegisterEncodeFunc : 可自定义加解密算法
func (gw *Gateway) RegisterEncodeFunc(f common.FuncTypeEncode) {

}

// RegisterDecodeFunc : 可自定义加解密算法
func (gw *Gateway) RegisterDecodeFunc(f common.FuncTypeDecode) {

}
