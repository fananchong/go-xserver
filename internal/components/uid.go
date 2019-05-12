package components

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/gomodule/redigo/redis"
)

// UID : UID 组件
type UID struct {
	ctx *common.Context
	cli go_redis_orm.IClient
}

// NewUID : 实例化
func NewUID(ctx *common.Context) *UID {
	uid := &UID{
		ctx: ctx,
		cli: go_redis_orm.GetDB(ctx.Config().DbAccount.Name),
	}
	ctx.IUID = uid
	return uid
}

// Start : 实例化组件
func (uid *UID) Start() bool {
	return true
}

// Close : 关闭组件
func (*UID) Close() {
	// No need to do anything
}

// GetUID : 根据 KEY , 获取 UID
func (uid *UID) GetUID(key string) (uint64, error) {
	return redis.Uint64(uid.cli.Do("HINCRBY", "idgen", key, 1))
}
