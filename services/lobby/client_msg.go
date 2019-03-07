package main

import (
	"github.com/fananchong/go-xserver/services/internal/db"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/gogo/protobuf/proto"
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
		accountObj.onLogin(data)
	case protocol.CMD_LOBBY_ENTER_GAME:
		accountObj.onEnterGame(data)
	default:
		Ctx.Log.Errorln("[LOBBY] Unknown cmd, cmd:", cmd)
	}
}

func (accountObj *Account) onLogin(data []byte) {
	Ctx.Log.Infoln("Login, account:", accountObj.account)
	msg := &protocol.MSG_LOBBY_LOGIN_RESULT{}
	for i, role := range accountObj.GetRoles() {
		if role != nil {
			info := &protocol.ROLE_BASE_INFO{}
			info.RoleID = role.Key
			info.RoleName = role.GetName()
			msg.Roles = append(msg.Roles, info)
			Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, info.RoleID, info.RoleName)
		} else {
			Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, 0, "''")
		}
	}
	utility.SendMsgToClient(accountObj.sess, accountObj.account, uint64(protocol.CMD_LOBBY_LOGIN), msg)
}

func (accountObj *Account) onEnterGame(data []byte) {
	Ctx.Log.Infoln("Enter Game, account:", accountObj.account)
	req := &protocol.MSG_LOBBY_ENTER_GAME{}
	if err := proto.Unmarshal(data, req); err != nil {
		Ctx.Log.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_ENTER_GAME`(", int(protocol.CMD_LOBBY_ENTER_GAME), "). account", accountObj.account)
		return
	}

	if req.Slot >= LimitRoleNum {
		Ctx.Log.Errorln("Message field error, Slot is ", req.Slot, ", but expect less than", LimitRoleNum, ". account:", accountObj.account)
		// TODO: 发送错误消息
		return
	}

	role := accountObj.roles[req.Slot]
	if role == nil { // 创建角色
		// 没有角色，则生成角色ID
		roleID, err := lobby.NewID(db.IDGenTypeRole)
		if err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			// TODO: 发送错误消息
			return
		}
		role, err = NewRole(roleID, accountObj.account)
		if err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			// TODO: 发送错误消息
			return
		}
		accountObj.roles[req.Slot] = role
		dbIDs := accountObj.RoleList.GetRoles(true)
		if dbIDs.GetRoleIDs() == nil {
			dbIDs.RoleIDs = make(map[uint32]uint64)
		}
		dbIDs.RoleIDs[req.Slot] = roleID
		if err := accountObj.Save(); err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			// TODO: 发送错误消息
			return
		}
	}
	accountObj.currentRole = role

	// TODO: 发送成功消息
}
