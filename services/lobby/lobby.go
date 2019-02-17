package main

import "github.com/fananchong/go-xserver/common"

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
	Ctx.Node.RegisterFuncOnRelayMsg(lobby.onRelayMsg)
	return true
}

// Close : 关闭
func (lobby *Lobby) Close() {

}

func (lobby *Lobby) onRelayMsg(source common.NodeType, account string, cmd uint64, data []byte) {
	switch source {
	case common.Client:
		lobby.onClientMsg(account, cmd, data)
	default:
		Ctx.Log.Errorln("Unknown source, type:", source, "(", int(source), ")")
	}
}
