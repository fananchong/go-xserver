package main

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/services/internal/protocol"
)

func (lobby *Lobby) onClientMsg(sess common.INode, account string, cmd uint64, data []byte) {
	switch protocol.CMD_LOBBY_ENUM(cmd) {
	case protocol.CMD_LOBBY_QUERY_ROLELIST:
		lobby.onQueryRoleList(sess, account, data)
	default:
		Ctx.Log.Errorln("[LOBBY] Unknown cmd, cmd:", cmd)
	}
}

func (lobby *Lobby) onQueryRoleList(sess common.INode, account string, data []byte) {
	// TODO: 先调通流程，逻辑待完善、代码待整理
	Ctx.Log.Infoln("Query role list, account:", account)
	msg := &protocol.MSG_LOBBY_QUERY_ROLELIST_RESULT{}
	_, err := NewAccount(account)
	if err != nil {
		Ctx.Log.Errorln(err, "account:", account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		sess.SendClientMsgByRelay(account, uint64(protocol.CMD_LOBBY_QUERY_ROLELIST), msg)
		return
	}
	sess.SendClientMsgByRelay(account, uint64(protocol.CMD_LOBBY_QUERY_ROLELIST), msg)
}
