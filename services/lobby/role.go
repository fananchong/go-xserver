package main

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// Role : 角色类
type Role struct {
	*db.Role
}

// NewRole : 角色类构造函数
func NewRole(roleID uint64) (*Role, error) {
	role := &Role{}
	role.Role = db.NewRole(Ctx.Config.DbAccount.Name, roleID)
	if err := role.Role.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return nil, err
		}
	}
	return role, nil
}
