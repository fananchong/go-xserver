package main

import (
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/go-xserver/services"
)

// Lobby : 大厅服务器
type Lobby struct {
	accountMgr *AccountMgr
}

// NewLobby : 构造函数
func NewLobby() *Lobby {
	lobby := &Lobby{
		accountMgr: NewAccountMgr(),
	}
	return lobby
}

// Start : 启动
func (lobby *Lobby) Start() bool {
	Ctx.EnableMessageRelay(true)
	Ctx.RegisterFuncOnRelayMsg(lobby.onRelayMsg)
	Ctx.RegisterFuncOnLoseAccount(lobby.onLoseAccount)
	return true
}

// Close : 关闭
func (lobby *Lobby) Close() {

}

func (lobby *Lobby) onRelayMsg(source config.NodeType, targetID context.NodeID, account string, cmd uint64, data []byte, flag uint8) {
	switch source {
	case config.Client:
		lobby.accountMgr.PostMsg(account, cmd, data, flag)
	case services.Match:
		lobby.onMatchMsg(targetID, cmd, data, flag)
	default:
		Ctx.Errorln("Unknown source, type:", source, "(", int(source), ")")
	}
}

func (lobby *Lobby) onLoseAccount(account string) {
	lobby.accountMgr.DelAccount(account)
}
