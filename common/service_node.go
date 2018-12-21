package common

// XSERVICENODE : 本服务节点对象
var XSERVICENODE IServerNode

// ServerIDType : 服务器ID类型
type ServerIDType [16]byte

// IServerNode : 服务节点类接口
type IServerNode interface {

	// 获取本节点信息
	GetId() ServerIDType    // 服务ID
	GetIpInner() string     // 内网IP
	GetIpOuter() string     // 外网IP
	GetPort(i int) int      // 端口
	GetOverload(i int) uint // 负载
	GetVersion() string     // 版本号

	// 根据服务类型，随机选择 1 节点
	SelectOne(t ServerType) IServerNode

	// 发送消息
	SendOne(t ServerType, data []byte) //
	// TODO： 发送系列待续
}
