package common

// XCONFIG : Global config function
var XCONFIG Config

type Config struct {
	Common    ConfigCommon // 一些基础参数
	DbAccount ConfigRedis  // 帐号数据库（Redis）
	DbToken   ConfigRedis  // Token数据库（Redis）
	DbServer  ConfigRedis  // Server数据库（Redis）
	DbMgr     ConfigRedis  // Mgr数据库（Redis）
	Login     ConfigLogin  // Login服务器配置
}

type ConfigCommon struct {
	Version       string // 版本号
	LogDir        string // log路径
	LogLevel      int    // log等级
	Debug         bool   // debug版本标志
	IntranetToken string // 内部服务器验证TOKEN
	MsgCmdOffset  int    // 消息号 = 服务类型 * MsgCmdOffset + 数字
	Pprof         string // http pprof 地址
}

type ConfigRedis struct {
	Name     string
	Addrs    []string
	Password string
	DBIndex  int
}

type ConfigLogin struct {
	Listen string
	Sign1  string
	Sign2  string
	Sign3  string
}
