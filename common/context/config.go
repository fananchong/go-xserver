package context

import "github.com/fananchong/go-xserver/config"

// IConfig : 配置类接口
type IConfig interface {
	LoadConfig(cfgfile string, cfgobj interface{}) bool // 逻辑层可以用该接口加载配置文件， cfgobj 为指针对象
	Config() *config.FrameworkConfig                    // 获取框架层配置
	PrintUsage()                                        // 打印命令行参数
}
