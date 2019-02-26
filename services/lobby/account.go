package main

import (
	"fmt"
	"reflect"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// LimitRoleNum : 一个账号角色数限制，最大为 1 个。
const LimitRoleNum = 1

// Account : 账号类
type Account struct {
	*db.RoleList
	roles []*Role
	idgen db.IDGen
}

// NewAccount : 角色列表类构造函数
func NewAccount(account string) (*Account, error) {
	accountobj := &Account{}
	accountobj.RoleList = db.NewRoleList(Ctx.Config.DbAccount.Name, account)
	if err := accountobj.RoleList.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return nil, err
		}
		if accountobj.FirstInitialization() == false {
			return nil, fmt.Errorf("Account failed to initialize for the first time, account:%s", account)
		}
	}
	accountobj.roles = make([]*Role, LimitRoleNum)
	var ids [256]uint64
	ids[0] = accountobj.RoleList.GetSlot0()
	ids[1] = accountobj.RoleList.GetSlot1()
	ids[2] = accountobj.RoleList.GetSlot2()
	ids[3] = accountobj.RoleList.GetSlot3()
	ids[4] = accountobj.RoleList.GetSlot4()
	for i := 0; i < LimitRoleNum; i++ {
		roleID := ids[i]
		if roleID == 0 {
			// 没有角色，则生成角色ID
			id, err := lobby.NewID(db.IDGenTypeRole)
			if err != nil {
				return nil, err
			}
			roleID = id
			// 角色创建时，这里才用下`反射`。因此不要当心`反射`造成的效率问题
			v := reflect.ValueOf(accountobj.RoleList)
			f := v.MethodByName(fmt.Sprintf("SetSlot%d", i))
			f.Call([]reflect.Value{reflect.ValueOf(roleID)})
		}
		role, err := NewRole(roleID, account)
		if err != nil {
			return nil, err
		}
		accountobj.roles[i] = role
	}
	if err := accountobj.Save(); err != nil {
		return nil, err
	}
	return accountobj, nil
}

// FirstInitialization : 账号首次创建初始化
func (accountobj *Account) FirstInitialization() bool {
	return true
}

// GetRoleIDs : 获取账号对应的角色列表
func (accountobj *Account) GetRoles() []*Role {
	return accountobj.roles
}
