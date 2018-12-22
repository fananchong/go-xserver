package common

import "github.com/fananchong/go-xserver/common/custom"

// XCONFIG : 全局配置类对象
var XCONFIG Config

// Config : 配置类
type Config struct {
	Common    ConfigCommon       // 一些基础参数
	DbAccount ConfigRedis        // 帐号数据库（Redis）
	DbToken   ConfigRedis        // Token 数据库（Redis）
	DbServer  ConfigRedis        // Server 数据库（Redis）
	DbMgr     ConfigRedis        // Mgr 数据库（Redis）
	Login     custom.ConfigLogin // Login 服务器配置
}

// ConfigCommon : 配置 common 节
type ConfigCommon struct {
	Version       string `default:"0.0.1" desc:"版本号"`
	LogDir        string `default:"./logs" desc:"Log 路径"`
	LogLevel      int    `default:"0" desc:"Log 等级"`
	Debug         bool   `default:"false" desc:"Debug 版本标志"`
	IntranetToken string `default:"6d8f1f3a-739f-47fe-9ed1-ea39276cd10d" desc:"内部服务器验证 TOKEN"`
	MsgCmdOffset  int    `default:"1000" desc:"消息号 = 服务类型 * MsgCmdOffset + 数字"`
	Pprof         string `default:"" desc:"Http pprof 地址"`
}

// ConfigRedis : 配置 redis 相关节
type ConfigRedis struct {
	Name     string   `desc:"Redis 数据库名称"`
	Addrs    []string `default:"[127.0.0.1:6379]" desc:"Redis 数据库地址"`
	Password string   `default:"" desc:"Redis 数据库密码"`
	DBIndex  int      `default:"0" desc:"Redis DB 索引"`
}
