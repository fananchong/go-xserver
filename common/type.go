package common

// NodeType : 节点类型
type NodeType int

const (

	// Client : 类型 0 ，客户端（虚拟）
	Client NodeType = iota

	// Mgr : 类型 1 ，管理服
	Mgr

	// Login : 类型 2 ，登录服
	Login

	// Gateway : 类型 3 ，网关服
	Gateway

	// Lobby : 类型 4 ，大厅服
	Lobby

	// Match : 类型 5 ，匹配服
	Match

	// Room : 类型 6 ，房间服
	Room

	// NodeTypeSize : 节点类型数量
	NodeTypeSize
)

// IPType : IP 类型
type IPType int

const (
	// IPINNER : 类型 0 ，内网 IP
	IPINNER IPType = iota

	// IPOUTTER : 类型 1 ，外网 IP
	IPOUTTER
)
