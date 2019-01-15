package db

import (
	"fmt"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/gomodule/redigo/redis"
)

// RedisAtomic : 自定义 Redis 原子操作
type RedisAtomic struct {
	Cli go_redis_orm.IClient
}

// SetX : 原子操作，封装 `不管 SETX 有没有设置成功，都重置过期时间为 n`
//        返回值， 空 表示设置成功； 非空 表示设置失败，并返回原先已设置的值
func (redisatomic *RedisAtomic) SetX(key, value string, expire int) (string, error) {
	lua := `local r=redis.call('SET',KEYS[1],ARGV[1],'NX');
redis.call('EXPIRE',KEYS[1],%d);
if not not r then return "" else redis.call('GET',KEYS[1]) end`
	cmd := fmt.Sprintf(lua, expire)
	reply, err := redisatomic.Cli.Do("eval %s 1 %s %s", cmd, key, value)
	if err != nil {
		return "", err
	}
	return redis.String(reply, nil)
}

// DelX : 原子操作，封装 `条件删除， if value = "xxx" then del key`
//        返回值， 1 有删除操作； 0 不需要删除
func (redisatomic *RedisAtomic) DelX(key, value string) (int, error) {
	cmd := `local v=redis.call('GET',KEYS[1]);
if ARGV[1]==v then return redis.call('DEL',KEYS[1]) else return 0 end`
	reply, err := redisatomic.Cli.Do("eval %s 1 %s %s", cmd, key, value)
	if err != nil {
		return -1, err
	}
	return redis.Int(reply, nil)
}
