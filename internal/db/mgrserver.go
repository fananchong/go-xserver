/// -------------------------------------------------------------------------------
/// THIS FILE IS ORIGINALLY GENERATED BY redis2go.exe.
/// PLEASE DO NOT MODIFY THIS FILE.
/// -------------------------------------------------------------------------------

package db

import (
	"errors"
	"fmt"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/gomodule/redigo/redis"
)

type MgrServer struct {
	Key  uint32
	addr string
	port int32

	__dirtyData               map[string]interface{}
	__dirtyDataForStructFiled map[string]interface{}
	__isLoad                  bool
	__dbKey                   string
	__dbName                  string
	__expire                  uint
}

func NewMgrServer(dbName string, key uint32) *MgrServer {
	return &MgrServer{
		Key:                       key,
		__dbName:                  dbName,
		__dbKey:                   "MgrServer:" + fmt.Sprintf("%d", key),
		__dirtyData:               make(map[string]interface{}),
		__dirtyDataForStructFiled: make(map[string]interface{}),
	}
}

// 若访问数据库失败返回-1；若 key 存在返回 1 ，否则返回 0 。
func (this *MgrServer) HasKey() (int, error) {
	db := go_redis_orm.GetDB(this.__dbName)
	val, err := redis.Int(db.Do("EXISTS", this.__dbKey))
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (this *MgrServer) Load() error {
	if this.__isLoad == true {
		return errors.New("alreay load!")
	}
	db := go_redis_orm.GetDB(this.__dbName)
	val, err := redis.Values(db.Do("HGETALL", this.__dbKey))
	if err != nil {
		return err
	}
	if len(val) == 0 {
		return go_redis_orm.ERR_ISNOT_EXIST_KEY
	}
	var data struct {
		Addr string `redis:"addr"`
		Port int32  `redis:"port"`
	}
	if err := redis.ScanStruct(val, &data); err != nil {
		return err
	}
	this.addr = data.Addr
	this.port = data.Port
	this.__isLoad = true
	return nil
}

func (this *MgrServer) Save() error {
	if len(this.__dirtyData) == 0 && len(this.__dirtyDataForStructFiled) == 0 {
		return nil
	}
	for k, _ := range this.__dirtyDataForStructFiled {
		_ = k

	}
	db := go_redis_orm.GetDB(this.__dbName)
	if _, err := db.Do("HMSET", redis.Args{}.Add(this.__dbKey).AddFlat(this.__dirtyData)...); err != nil {
		return err
	}
	if this.__expire != 0 {
		if _, err := db.Do("EXPIRE", this.__dbKey, this.__expire); err != nil {
			return err
		}
	}
	this.__dirtyData = make(map[string]interface{})
	this.__dirtyDataForStructFiled = make(map[string]interface{})
	return nil
}

func (this *MgrServer) Delete() error {
	db := go_redis_orm.GetDB(this.__dbName)
	_, err := db.Do("DEL", this.__dbKey)
	if err == nil {
		this.__isLoad = false
		this.__dirtyData = make(map[string]interface{})
		this.__dirtyDataForStructFiled = make(map[string]interface{})
	}
	return err
}

func (this *MgrServer) IsLoad() bool {
	return this.__isLoad
}

func (this *MgrServer) Expire(v uint) {
	this.__expire = v
}

func (this *MgrServer) DirtyData() (map[string]interface{}, error) {
	for k, _ := range this.__dirtyDataForStructFiled {
		_ = k

	}
	data := make(map[string]interface{})
	for k, v := range this.__dirtyData {
		data[k] = v
	}
	this.__dirtyData = make(map[string]interface{})
	this.__dirtyDataForStructFiled = make(map[string]interface{})
	return data, nil
}

func (this *MgrServer) Save2(dirtyData map[string]interface{}) error {
	if len(dirtyData) == 0 {
		return nil
	}
	db := go_redis_orm.GetDB(this.__dbName)
	if _, err := db.Do("HMSET", redis.Args{}.Add(this.__dbKey).AddFlat(dirtyData)...); err != nil {
		return err
	}
	if this.__expire != 0 {
		if _, err := db.Do("EXPIRE", this.__dbKey, this.__expire); err != nil {
			return err
		}
	}
	return nil
}

func (this *MgrServer) GetAddr() string {
	return this.addr
}

func (this *MgrServer) GetPort() int32 {
	return this.port
}

func (this *MgrServer) SetAddr(value string) {
	this.addr = value
	this.__dirtyData["addr"] = string([]byte(value))
}

func (this *MgrServer) SetPort(value int32) {
	this.port = value
	this.__dirtyData["port"] = value
}
