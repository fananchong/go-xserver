package common

var (
	// XTCPSVRFORCLIENT : 对外 TCP 服务
	XTCPSVRFORCLIENT ITCPServer

	// XTCPSVRFORINTRANET : 对内 TCP 服务
	XTCPSVRFORINTRANET ITCPServer
)

// ITCPServer : TCP 服务接口
type ITCPServer interface {
	RegisterSessType(v interface{})
}

// 框架层内置的 2 个 TCP 服务
// 对外 TCP 服务， IP 配置对应 --network-ipouter 参数； PORT 配置对应 --network-port 参数的第一个值
// 对内 TCP 服务， IP 配置对应 --network-ipinner 参数； PORT 配置对应 --network-port 参数的第二个值
// 逻辑层只有调用 ITCPServer.RegisterSessType ，注册 Session 类，框架层才会开启服务
