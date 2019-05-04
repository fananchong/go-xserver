package config

// WARNNING : 请勿修改本文件，本文件由框架层程序来维护

// NodeType : 节点类型
type NodeType int

const (

	// Client : 类型 0 ，客户端（虚拟）
	Client NodeType = 0

	// Mgr : 类型 1 ，管理服
	Mgr NodeType = 1

	// Login : 类型 2 ，登录服
	Login NodeType = 2

	// Gateway : 类型 3 ，网关服
	Gateway NodeType = 3

	// Unknow : 未知
	Unknow NodeType = 9999
)
