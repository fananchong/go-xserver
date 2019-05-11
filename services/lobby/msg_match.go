package main

import (
	"github.com/fananchong/go-xserver/services"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/fananchong/gotcp"
)

// onMatch : 请求匹配
func (accountObj *Account) onMatch(data []byte) {
	Ctx.Infoln("Match, account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	msg := &protocol.MSG_LOBBY_MATCH{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_LOBBY_MATCH`(", int(protocol.CMD_LOBBY_MATCH), "). account", accountObj.account, "roleid:", accountObj.currentRole.Key)
		return
	}

	// 匹配请求可重入， Match 服加判断忽略

	// 请求 Match 服
	matchMsg := &protocol.MSG_MATCH_MATCH{}
	matchMsg.RoleID = accountObj.currentRole.Key
	if _, err := utility.SendMsgToServer(Ctx, services.Match, uint64(protocol.CMD_MATCH_MATCH), matchMsg); err != nil {
		Ctx.Errorln("error:", err.Error(), "account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	}
}
