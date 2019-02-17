package main

import "github.com/fananchong/go-xserver/services/internal/protocol"

func (lobby *Lobby) onClientMsg(account string, cmd uint64, data []byte) {
	switch protocol.CMD_LOBBY_ENUM(cmd) {
	case protocol.CMD_LOBBY_QUERY_ROLELIST:
	default:
		Ctx.Log.Errorln("[LOBBY] Unknown cmd, cmd:", cmd)
	}
}
