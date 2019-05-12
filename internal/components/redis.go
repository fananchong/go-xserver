package components

import (
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
)

// Redis : Redis 组件
type Redis struct {
	ctx *common.Context
}

// NewRedis : 实例化
func NewRedis(ctx *common.Context) *Redis {
	redis := &Redis{ctx: ctx}

	// TODO: go_redis_orm 可以实例化，而非全局的
	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)

	cfgs := []config.FrameworkConfigRedis{
		redis.ctx.Config().DbMgr,
		redis.ctx.Config().DbRoleName,
		redis.ctx.Config().DbServer,
		redis.ctx.Config().DbAccount,
	}
	for _, cfg := range cfgs {
	LOOP:
		if err := go_redis_orm.CreateDB(cfg.Name, cfg.Addrs, cfg.Password, cfg.DBIndex); err != nil {
			redis.ctx.Errorln(err)
			time.Sleep(5 * time.Second)
			goto LOOP
		}
	}
	return redis
}

// Start : 实例化组件
func (redis *Redis) Start() bool {
	return true
}

// Close : 关闭组件
func (*Redis) Close() {
	// No need to do anything
}
