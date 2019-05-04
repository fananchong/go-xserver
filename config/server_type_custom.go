package config

// WARNNING : 逻辑层程序在该文件中，定义自己的服务类型
//            服务类型号码请勿与 server_type_framework.go 中的重复

const (

	// Lobby : 类型 4 ，大厅服
	Lobby NodeType = 4

	// Match : 类型 5 ，匹配服
	Match NodeType = 5

	// Room : 类型 6 ，房间服
	Room NodeType = 6
)
