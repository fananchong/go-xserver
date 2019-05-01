package common

// IPlugin : 插件接口
type IPlugin interface {
	Start() bool
	Close()
}
