package main

// Gateway : 网关服务器
type Gateway struct {
	*UserMgr
}

// NewGateway : 构造函数
func NewGateway() *Gateway {
	gw := &Gateway{}
	gw.UserMgr = NewUserMgr()
	return gw
}

// Start : 启动
func (gateway *Gateway) Start() bool {
	Ctx.RegisterSessType(User{})
	Ctx.RegisterSendToClient(gateway.sendToClient)
	Ctx.RegisterSendToAllClient(gateway.sendToAllClient)
	return true
}

// Close : 关闭
func (gateway *Gateway) Close() {

}

func (gateway *Gateway) sendToClient(account string, cmd uint64, data []byte) bool {
	if user := gateway.GetUser(account); user != nil {
		datalen := len(data) - 1
		if user.SendEx(int(cmd), data[:datalen], data[datalen]) {
			return true
		}
		Ctx.Warning("Sending message failed, account:", account, ", cmd:", cmd)
		return false
	}
	Ctx.Warning("The player was not found, account:", account)
	return false
}

func (gateway *Gateway) sendToAllClient(cmd uint64, data []byte) bool {
	datalen := len(data) - 1
	msg := data[:datalen]
	flag := data[datalen]
	gateway.Foreach(func(user *User) bool {
		if user.SendEx(int(cmd), msg, flag) == false {
			Ctx.Warning("Sending message failed, account:", user.GetAccount(), ", cmd:", cmd)
		}
		return true
	})
	return true
}
