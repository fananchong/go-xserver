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
	return true
}

// Close : 关闭
func (gateway *Gateway) Close() {

}
