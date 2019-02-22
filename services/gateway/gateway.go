package main

// Gateway : 网关服务器
type Gateway struct {
}

// NewGateway : 构造函数
func NewGateway() *Gateway {
	return &Gateway{}
}

// Start : 启动
func (gateway *Gateway) Start() bool {
	Ctx.ServerForClient.RegisterSessType(User{})
	Ctx.Gateway.RegisterSendToClient(gateway.sendToClient)
	return true
}

// Close : 关闭
func (gateway *Gateway) Close() {

}

func (gateway *Gateway) sendToClient(account string, cmd uint64, data []byte) bool {
	Ctx.Log.Infoln("1111111111111111111111")
	// TODO: 根据 account，找到对应 user ， 发送之，待写下，稍后
	return false
}
