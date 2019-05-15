package main

import (
	"github.com/fananchong/go-xserver/common/context"
	"github.com/fananchong/go-xserver/services"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/fananchong/gotcp"
)

// onMatch : 请求匹配
func (accountObj *Account) onMatch(data []byte, flag uint8) {
	Ctx.Infoln("Match, account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	msg := &protocol.MSG_LOBBY_MATCH{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_LOBBY_MATCH`(", int(protocol.CMD_LOBBY_MATCH), "). account", accountObj.account, "roleid:", accountObj.currentRole.Key)
		return
	}

	// 匹配请求可重入， Match 服加判断忽略

	// 请求 Match 服
	matchMsg := &protocol.MSG_MATCH_MATCH{}
	matchMsg.Account = accountObj.account
	matchMsg.RoleID = accountObj.currentRole.Key
	if _, err := utility.SendMsgToServer(Ctx, services.Match, uint64(protocol.CMD_MATCH_MATCH), matchMsg); err != nil {
		Ctx.Errorln("error:", err.Error(), "account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	}
}

func (accountObj *Account) onMatchResult(data []byte, flag uint8) {
	Ctx.Infoln("MatchResult, account:", accountObj.account, "roleid:", accountObj.currentRole.Key)
	msg := &protocol.MSG_LOBBY_MATCH_RESULT{}
	if gotcp.DecodeCmd(data, flag, msg) == nil {
		Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_LOBBY_MATCH_RESULT`(", int(protocol.CMD_LOBBY_MATCH_RESULT), ")")
		return
	}
	utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_MATCH_RESULT), msg)
}

func (lobby *Lobby) onMatchMsg(targetID context.NodeID, cmd uint64, data []byte, flag uint8) {
	switch protocol.CMD_MATCH_ENUM(cmd) {
	case protocol.CMD_MATCH_MATCH:
		msg := &protocol.MSG_MATCH_MATCH_RESULT{}
		if gotcp.DecodeCmd(data, flag, msg) == nil {
			Ctx.Errorln("Message parsing failed, message number is`protocol.CMD_MATCH_MATCH`(", int(protocol.CMD_MATCH_MATCH), ")")
			return
		}
		tmpMsg := &protocol.MSG_LOBBY_MATCH_RESULT{}
		tmpMsg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_ENUM(msg.GetErr())
		tmpMsg.Roles = append(tmpMsg.Roles, msg.GetRoles()...)
		// 暂时序列化下，更好的，应该有传递 msg 的接口
		data, flag, err := gotcp.EncodeCmd(uint64(protocol.CMD_LOBBY_MATCH_RESULT), tmpMsg)
		if err != nil {
			Ctx.Errorf("error:%s, account:%s, roleid:%d", err.Error(), msg.GetAccount(), msg.GetRoleID())
			return
		}
		lobby.accountMgr.PostMsg(msg.GetAccount(), uint64(protocol.CMD_LOBBY_MATCH_RESULT), data, flag)
	default:
		Ctx.Errorln("Unknown cmd:", cmd)
	}
}
