package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/fananchong/gotcp"
)

// onChat : 聊天
func (accountObj *Account) onChat(data []byte) {
	Ctx.Log.Infoln("Chat, account:", accountObj.account)
	msg := &protocol.MSG_LOBBY_CHAT{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], msg) == nil {
		Ctx.Log.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_CHAT`(", int(protocol.CMD_LOBBY_CHAT), "). account", accountObj.account)
		return
	}

	// TODO: 待实现
	//           1. 当 msg.GetTo() 不为空，查找其所在的 Gateway ，并发送 msg
	//           2. 当 msg.GetTo() 为空，向所有 Gateway 发送 msg
	//       以上具有代表性，框架层可以提供以下接口或功能：
	//           1. 根据 account / roleID / roleName ，向指定的服务类型发送消息
	//           2. 向某类型服务广播消息
	//           3. Gateway 处理全服广播

	msg.From = accountObj.GetRole().GetName()
	if msg.GetTo() == "" {
		// 全服广播
		utility.BroadcastMsgToClient(Ctx, uint64(protocol.CMD_LOBBY_CHAT), msg)
	} else {
		// 私聊
	}
	// utility.SendMsgToClient(accountObj.sess, accountObj.account, uint64(protocol.CMD_LOBBY_CHAT), msg)
}
