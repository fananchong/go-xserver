package common

// ServerType : 基础服务类型
type ServerType int

const (

	// Client : 类型为客户端 0
	Client ServerType = iota

	// Mgr : 类型为管理服 1
	Mgr

	// Login : 类型为登录服 2
	Login

	// Gateway : 类型为网关服 3
	Gateway

	// ServerTypeSize : 服务类型数量
	ServerTypeSize
)
