package main

// Lobby : 大厅服务器
type Lobby struct {
}

// NewLobby : 构造函数
func NewLobby() *Lobby {
	return &Lobby{}
}

// Start : 启动
func (lobby *Lobby) Start() bool {
	Ctx.Node.EnableMessageRelay(true)
	return true
}

// Close : 关闭
func (lobby *Lobby) Close() {

}
