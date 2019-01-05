package components

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
)

// Redis : Redis 组件
type Redis struct {
}

// Start : 实例化组件
func (*Redis) Start() bool {
	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	if err := go_redis_orm.CreateDB(
		common.XCONFIG.DbMgr.Name,
		common.XCONFIG.DbMgr.Addrs,
		common.XCONFIG.DbMgr.Password,
		common.XCONFIG.DbMgr.DBIndex); err != nil {
		common.XLOG.Errorln(err)
		return false
	}
	return true
}

// Close : 关闭组件
func (*Redis) Close() {
	// do nothing
}
