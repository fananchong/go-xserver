package common

// Plugin : 插件接口
type Plugin interface {
	Init() (nodeType NodeType, ok bool)
	RegisterCallBack()
	Start() bool
	Close()
}

// Plugin 插件接口类连接着框架层与逻辑层
// 负责程序启动时，框架层、逻辑层间的数据传递
// 主要流程如下：
//     1. go-xserver 启动，初始化框架层基础功能模块
//     2. go-xserver 加载 Plugin 对象，并调用 Plugin.Init
//     3. Plugin.Init 告知 go-xserver `节点类型`
//     4. go-xserver 根据`节点类型`，初始化依赖`节点类型`信息的功能模块
//     5. go-xserver 调用 Plugin.RegisterCallBack 把逻辑层自定义回调注册进框架层
//     6. go-xserver 启动框架层所有模块
//     7. go-xserver 调用 Plugin.Start 启动逻辑层模块
//     8. 逻辑层功能开始运作
