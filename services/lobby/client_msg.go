package main

import (
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
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
	Ctx.Log.Infoln("Query role list, account:", account)
	msg := &protocol.MSG_LOBBY_QUERY_ROLELIST_RESULT{}
	accountObj, err := NewAccount(account)
	if err != nil {
		Ctx.Log.Errorln(err, "account:", account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		utility.SendMsgToClient(sess, account, uint64(protocol.CMD_LOBBY_QUERY_ROLELIST), msg)
		return
	}
	for i, role := range accountObj.GetRoles() {
		info := &protocol.ROLE_BASE_INFO{}
		info.RoleID = role.Key
		info.RoleName = role.GetName()
		msg.Roles = append(msg.Roles, info)
		Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, info.RoleID, info.RoleName)
	}
	utility.SendMsgToClient(sess, account, uint64(protocol.CMD_LOBBY_QUERY_ROLELIST), msg)
}
