package main

import "github.com/fananchong/go-xserver/services/internal/protocol"

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
		accountObj.onLogin(data)
	case protocol.CMD_LOBBY_CREATE_ROLE:
		accountObj.onCreateRole(data)
	case protocol.CMD_LOBBY_ENTER_GAME:
		accountObj.onEnterGame(data)
	default:
		if accountObj.GetRole() == nil {
			Ctx.Errorln("[LOBBY] Login not completed. account", accountObj.account, ",cmd:", cmd)
			return
		}
		switch protocol.CMD_LOBBY_ENUM(cmd) {
		case protocol.CMD_LOBBY_CHAT:
			accountObj.onChat(data)
		default:
			Ctx.Errorln("[LOBBY] Unknown cmd, cmd:", cmd)
		}
	}
}
