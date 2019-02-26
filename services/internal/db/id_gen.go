package db

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/gomodule/redigo/redis"
)

// IDGenType : ID 类型
type IDGenType int

const (
	// IDGenTypeAccount :
	IDGenTypeAccount = iota
	// IDGenTypeRole :
	IDGenTypeRole
)

// IDGen : ID 生成器（服务器组内唯一）
type IDGen struct {
	Cli go_redis_orm.IClient
}

// NewID : 获取一个 ID
func (idgen *IDGen) NewID(typ IDGenType) (uint64, error) {
	return redis.Uint64(idgen.Cli.Do("HINCRBY", "idgen", typ, 1))
}
