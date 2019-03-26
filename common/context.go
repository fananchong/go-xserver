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
	Role2Account    IRole2Account   // 常驻功能
	ServerForClient ITCPServer      // 注册该字段相应接口，才会开启（可选模块）
	Login           ILogin          // 常驻功能。节点类型为 Login，才会开启
	Gateway         IGateway        // 常驻功能。节点类型为 Gateway ，才会开启
}
