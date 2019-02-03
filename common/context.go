package common

import (
	"context"
	"math/rand"
)

// Context : 应用程序上下文
type Context struct {
	Ctx             context.Context // 常驻功能
	Rand            *rand.Rand      // 常驻功能
	Config          *Config         // 常驻功能
	Log             ILogger         // 常驻功能
	Node            INode           // 常驻功能
	ServerForClient ITCPServer      // 注册该字段相应接口，才会开启
	Login           ILogin          // 节点类型为 Login，才会开启
	Gateway         IGateway        // 节点类型为 Gateway ，才会开启
}
