package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/gotcp"
)

// onMatch : 请求匹配
func (accountObj *Account) onMatch(data []byte) {
	Ctx.Infoln("Match, account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	msg := &protocol.MSG_LOBBY_MATCH{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_MATCH`(", int(protocol.CMD_LOBBY_MATCH), "). account", accountObj.account, "roleid:", accountObj.currentRole.Key)
		return
	}

	// 匹配请求可重入， Match 服加判断忽略

	// TODO: 转发给 Match 服

}
