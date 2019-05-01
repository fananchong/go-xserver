package config

// NodeType : 节点类型
type NodeType int

const (

	// Client : 类型 0 ，客户端（虚拟）。 框架层会用到，请勿删及更改值
	Client NodeType = iota

	// Mgr : 类型 1 ，管理服。 框架层会用到，请勿删及更改值
	Mgr = 1

	// Login : 类型 2 ，登录服。 框架层会用到，请勿删及更改值
	Login = 2

	// Gateway : 类型 3 ，网关服。 框架层会用到，请勿删及更改值
	Gateway = 3

	// Lobby : 类型 4 ，大厅服
	Lobby = 4

	// Match : 类型 5 ，匹配服
	Match = 5

	// Room : 类型 6 ，房间服
	Room = 6

	// Unknow : 未知。 框架层会用到，请勿删及更改值
	Unknow = 9999
)
