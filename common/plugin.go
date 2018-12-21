package common

// Plugin : 插件接口
type Plugin interface {
	Init() (nodeType NodeType, ok bool)
	Start() bool
	Close()
}
