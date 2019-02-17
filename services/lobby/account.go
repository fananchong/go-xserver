package main

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// LimitRoleNum : 一个账号角色数限制，最大为 1 个。
const LimitRoleNum = 1

// Account : 账号类
type Account struct {
	*db.RoleList
	roles []*Role
}

// NewAccount : 角色列表类构造函数
func NewAccount(account string) (*Account, error) {
	accountobj := &Account{}
	accountobj.RoleList = db.NewRoleList(Ctx.Config.DbAccount.Name, account)
	if err := accountobj.RoleList.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return nil, err
		}
	}
	accountobj.roles = make([]*Role, LimitRoleNum)
	var ids [16]uint64
	ids[0] = accountobj.RoleList.GetSlot1()
	ids[1] = accountobj.RoleList.GetSlot2()
	ids[2] = accountobj.RoleList.GetSlot3()
	ids[3] = accountobj.RoleList.GetSlot4()
	ids[4] = accountobj.RoleList.GetSlot5()
	for i := 0; i < LimitRoleNum; i++ {
		roleID := ids[i]
		if roleID != 0 {
			role, err := NewRole(roleID)
			if err != nil {
				return nil, err
			}
			accountobj.roles[i] = role
		}
	}
	return accountobj, nil
}
