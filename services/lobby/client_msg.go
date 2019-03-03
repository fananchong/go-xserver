package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
)

// ChanMsg : 账号消息
type ChanMsg struct {
	Cmd  uint64
	Data []byte
}

// PostMsg : 推送消息
func (accountObj *Account) PostMsg(cmd uint64, data []byte) {
	datacopy := make([]byte, len(data))
	copy(datacopy, data)
	accountObj.chanMsg <- ChanMsg{cmd, datacopy}
}

// ProcessMsg : 处理消息
func (accountObj *Account) processMsg(cmd uint64, data []byte) {
	switch protocol.CMD_LOBBY_ENUM(cmd) {
	case protocol.CMD_LOBBY_LOGIN:
		accountObj.onQueryRoleList(data)
	default:
		Ctx.Log.Errorln("[LOBBY] Unknown cmd, cmd:", cmd)
	}
}

func (accountObj *Account) onQueryRoleList(data []byte) {
	Ctx.Log.Infoln("Query role list, account:", accountObj.account)
	msg := &protocol.MSG_LOBBY_LOGIN_RESULT{}
	for i, role := range accountObj.GetRoles() {
		info := &protocol.ROLE_BASE_INFO{}
		info.RoleID = role.Key
		info.RoleName = role.GetName()
		msg.Roles = append(msg.Roles, info)
		Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, info.RoleID, info.RoleName)
	}
	utility.SendMsgToClient(accountObj.sess, accountObj.account, uint64(protocol.CMD_LOBBY_LOGIN), msg)
}
