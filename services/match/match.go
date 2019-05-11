package main

import (
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/services"
)

// Match : 匹配服务器
type Match struct {
}

// NewMatch : 构造函数
func NewMatch() *Match {
	match := &Match{}
	return match
}

// Start : 启动
func (match *Match) Start() bool {
	Ctx.EnableMessageRelay(true)
	Ctx.RegisterFuncOnRelayMsg(match.onRelayMsg)
	return true
}

// Close : 关闭
func (match *Match) Close() {
}

func (match *Match) onRelayMsg(source config.NodeType, _ string, cmd uint64, data []byte) {
	switch source {
	case services.Lobby:
		match.onLobbyMsg(cmd, data)
	default:
		Ctx.Errorln("Unknown source, type:", source, "(", int(source), ")")
	}
}
