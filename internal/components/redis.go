package components

import (
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
)

// Redis : Redis 组件
type Redis struct {
	ctx *common.Context
}

// NewRedis : 实例化
func NewRedis(ctx *common.Context) *Redis {
	// TODO: go_redis_orm 可以实例化，而非全局的
	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	return &Redis{ctx: ctx}
}

// Start : 实例化组件
func (redis *Redis) Start() bool {
LOOP0:
	if err := go_redis_orm.CreateDB(
		redis.ctx.Config.DbMgr.Name,
		redis.ctx.Config.DbMgr.Addrs,
		redis.ctx.Config.DbMgr.Password,
		redis.ctx.Config.DbMgr.DBIndex); err != nil {
		redis.ctx.Log.Errorln(err)
		time.Sleep(5 * time.Second)
		goto LOOP0
	}

	if pluginType != common.Login {
	LOOP1:
		if err := go_redis_orm.CreateDB(
			redis.ctx.Config.DbServer.Name,
			redis.ctx.Config.DbServer.Addrs,
			redis.ctx.Config.DbServer.Password,
			redis.ctx.Config.DbServer.DBIndex); err != nil {
			redis.ctx.Log.Errorln(err)
			time.Sleep(5 * time.Second)
			goto LOOP1
		}
	}

	return true
}

// Close : 关闭组件
func (*Redis) Close() {
	// No need to do anything
}
