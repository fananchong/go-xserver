package main

import (
	"github.com/fananchong/go-xserver/services/internal/db"
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/fananchong/gotcp"
)

// onLogin : 获取角色列表（登录大厅服务）
func (accountObj *Account) onLogin(data []byte) {
	Ctx.Log.Infoln("Login, account:", accountObj.account)
	msg := &protocol.MSG_LOBBY_LOGIN_RESULT{}
	for i, role := range accountObj.GetRoles() {
		info := &protocol.ROLE_BASE_INFO{}
		if role != nil {
			info.RoleID = role.Key
			info.RoleName = role.GetName()
			Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, info.RoleID, info.RoleName)
		} else {
			Ctx.Log.Infof("\t[role%d] roleid:%d, rolename:%s\n", i, 0, "''")
		}
		msg.Roles = append(msg.Roles, info)
	}
	utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_LOGIN), msg)
}

// onCreateRole : 创建角色
func (accountObj *Account) onCreateRole(data []byte) {
	Ctx.Log.Infoln("Create role, account:", accountObj.account)
	req := &protocol.MSG_LOBBY_CREATE_ROLE{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], req) == nil {
		Ctx.Log.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_CREATE_ROLE`(", int(protocol.CMD_LOBBY_CREATE_ROLE), "). account", accountObj.account)
		return
	}

	msg := &protocol.MSG_LOBBY_CREATE_ROLE_RESULT{}
	if req.Slot >= LimitRoleNum {
		Ctx.Log.Errorln("Message field error, Slot is ", req.Slot, ", but expect less than", LimitRoleNum, ". account:", accountObj.account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
		return
	}

	if req.GetInfo() == nil {
		Ctx.Log.Errorln("Message field error, Info is nil. account:", accountObj.account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
		return
	}

	// TODO: 重名检查

	role := accountObj.roles[req.Slot]
	if role == nil { // 创建角色
		// 没有角色，则生成角色ID
		roleID, err := lobby.NewID(db.IDGenTypeRole)
		if err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
			utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
			return
		}
		// 生成角色
		role, err = NewRole(roleID, accountObj.account)
		if err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
			utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
			return
		}
		role.SetName(req.GetInfo().GetRoleName())
		if err = role.Save(); err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
			utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
			return
		}
		// 关联账号
		accountObj.roles[req.Slot] = role
		dbIDs := accountObj.RoleList.GetRoles(true)
		if dbIDs.GetRoleIDs() == nil {
			dbIDs.RoleIDs = make(map[uint32]uint64)
		}
		dbIDs.RoleIDs[req.Slot] = roleID
		if err := accountObj.Save(); err != nil {
			Ctx.Log.Errorln(err, "account:", accountObj.account)
			msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
			utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
			return
		}
	}
	utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_CREATE_ROLE), msg)
}

// onEnterGame : 获取角色详细信息（进入游戏）
func (accountObj *Account) onEnterGame(data []byte) {
	Ctx.Log.Infoln("Enter Game, account:", accountObj.account)
	req := &protocol.MSG_LOBBY_ENTER_GAME{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], req) == nil {
		Ctx.Log.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_ENTER_GAME`(", int(protocol.CMD_LOBBY_ENTER_GAME), "). account", accountObj.account)
		return
	}

	msg := &protocol.MSG_LOBBY_ENTER_GAME_RESULT{}
	if req.Slot >= LimitRoleNum {
		Ctx.Log.Errorln("Message field error, Slot is ", req.Slot, ", but expect less than", LimitRoleNum, ". account:", accountObj.account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_ENTER_GAME), msg)
		return
	}

	role := accountObj.roles[req.Slot]
	if role == nil {
		// 没有角色
		Ctx.Log.Errorln("No role found, Slot is ", req.Slot, ", account:", accountObj.account)
		msg.Err = protocol.ENUM_LOBBY_COMMON_ERROR_SYSTEM_ERROR
		utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_ENTER_GAME), msg)
		return
	}
	accountObj.currentRole = role
	msg.DetailInfo = &protocol.ROLE_DETAIL_INFO{}
	msg.DetailInfo.BaseInfo = &protocol.ROLE_BASE_INFO{}
	msg.DetailInfo.BaseInfo.RoleID = role.Key
	msg.DetailInfo.BaseInfo.RoleName = role.GetName()

	// TODO: 未加载角色各细节数据，则加载之

	utility.SendMsgToClient(Ctx, accountObj.account, uint64(protocol.CMD_LOBBY_ENTER_GAME), msg)
}
