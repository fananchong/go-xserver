package common

// Plugin : 插件接口
type Plugin interface {
	Init() bool
	Start() bool
	Close()
}
