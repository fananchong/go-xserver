package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
)

func (match *Match) onLobbyMsg(cmd uint64, data []byte) {
	switch protocol.CMD_MATCH_ENUM(cmd) {
	case protocol.CMD_MATCH_MATCH:
		req := &protocol.MSG_MATCH_MATCH{}
		if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], req) == nil {
			Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_MATCH_MATCH`(", int(protocol.CMD_MATCH_MATCH), ")")
			return
		}
		Ctx.Infoln("Request match, roleid:", req.GetRoleID())
		// TODO:
	}
}
