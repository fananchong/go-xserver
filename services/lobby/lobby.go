package main

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
)

// Lobby : 大厅服务器
type Lobby struct {
}

// NewLobby : 构造函数
func NewLobby() *Lobby {
	return &Lobby{}
}

// Start : 启动
func (lobby *Lobby) Start() bool {
	if lobby.initRedis() == false {
		return false
	}
	Ctx.Node.EnableMessageRelay(true)
	Ctx.Node.RegisterFuncOnRelayMsg(lobby.onRelayMsg)
	return true
}

// Close : 关闭
func (lobby *Lobby) Close() {

}

func (lobby *Lobby) onRelayMsg(source common.NodeType, sess common.INode, account string, cmd uint64, data []byte) {
	switch source {
	case common.Client:
		lobby.onClientMsg(sess, account, cmd, data)
	default:
		Ctx.Log.Errorln("Unknown source, type:", source, "(", int(source), ")")
	}
}

func (lobby *Lobby) initRedis() bool {
	// db account
	err := go_redis_orm.CreateDB(
		Ctx.Config.DbAccount.Name,
		Ctx.Config.DbAccount.Addrs,
		Ctx.Config.DbAccount.Password,
		Ctx.Config.DbAccount.DBIndex)
	if err != nil {
		Ctx.Log.Errorln(err)
		return false
	}
	return true
}
