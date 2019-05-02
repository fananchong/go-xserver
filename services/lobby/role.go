package main

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// Role : 角色类
type Role struct {
	*db.RoleBase
	account string
	inited  bool
}

// NewRole : 角色类构造函数
func NewRole(roleID uint64, account string) (*Role, error) {
	role := &Role{}
	role.RoleBase = db.NewRoleBase(Ctx.Config().DbAccount.Name, roleID)
	role.account = account
	if err := role.RoleBase.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return nil, err
		}
	}
	return role, nil
}

// FirstInitialization : 角色首次创建初始化
func (role *Role) FirstInitialization() bool {
	return true
}

// InitFromDB : 从数据库加载数据
func (role *Role) InitFromDB() bool {
	role.inited = true
	return true
}
