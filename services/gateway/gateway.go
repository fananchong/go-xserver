package main

// Gateway : 登陆服务器
type Gateway struct {
}

// NewGateway : 构造函数
func NewGateway() *Gateway {
	return &Gateway{}
}

// Start : 启动
func (gateway *Gateway) Start() bool {
	Ctx.ServerForClient.RegisterSessType(User{})
	Ctx.ServerForIntranet.RegisterSessType(IntranetNode{})
	return true
}

// Close : 关闭
func (gateway *Gateway) Close() {

}
