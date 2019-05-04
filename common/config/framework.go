package config

// WARNNING : 请勿修改本文件，本文件由框架层程序来维护

// FrameworkConfig : 配置类
type FrameworkConfig struct {
	Common     frameworkConfigCommon  // 一些基础参数
	Network    frameworkConfigNetwork // 网络配置
	DbAccount  frameworkConfigRedis   // 帐号数据库（Redis）
	DbToken    frameworkConfigRedis   // Token 数据库（Redis）
	DbServer   frameworkConfigRedis   // Server 数据库（Redis）
	DbMgr      frameworkConfigRedis   // Mgr 数据库（Redis）
	DbRoleName frameworkConfigRedis   // 角色名-账号数据库（Redis）
	Role       frameworkConfigRole    // 角色相关配置
}

type frameworkConfigCommon struct {
	Version          string `default:"0.0.1" desc:"版本号"`
	LogDir           string `default:"./logs" desc:"Log 路径"`
	LogLevel         int    `default:"1" desc:"Log 等级， 1 infoLog; 2 warningLog; 3 errorLog; 4 fatalLog"`
	LogFlushInterval int    `default:"1000" desc:"Log 写入到文件的时间间隔，单位：毫秒"`
	Debug            bool   `default:"true" desc:"Debug 版本标志"`
	IntranetToken    string `default:"6d8f1f3a-739f-47fe-9ed1-ea39276cd10d" desc:"内部服务器验证 TOKEN"`
	MsgCmdOffset     int    `default:"1000" desc:"消息号 = 服务类型 * MsgCmdOffset + 数字"`
	Pprof            string `default:"" desc:"Http pprof 地址"`
}

type frameworkConfigNetwork struct {
	IPType  int     `default:"1" desc:"类型： 1 表示使用 IP；类型： 0 表示使用网卡名"`
	IPInner string  `default:"127.0.0.1" desc:"内网 IP"`
	IPOuter string  `default:"127.0.0.1" desc:"外网 IP"`
	Port    []int32 `default:"[7500, 30000]" desc:"第一个端口为对外提供服务的端口；第二个端口为服务器组内提供服务的端口；若有其他继续填充，自定义"`
}

type frameworkConfigRedis struct {
	Name     string   `desc:"Redis 数据库名称"`
	Addrs    []string `default:"[127.0.0.1:6379]" desc:"Redis 数据库地址"`
	Password string   `default:"123456" desc:"Redis 数据库密码"`
	DBIndex  int      `default:"1" desc:"Redis DB 索引"`
}

type frameworkConfigRole struct {
	IdleTime                int64 `default:"300" desc:"角色闲置该时间间隔后账号、角色对象内存中清除。单位：秒"`
	SessionAffinityInterval int64 `default:"300" desc:"账号登出后，分配的服务器资源信息保留时间。单位：秒"`
}
