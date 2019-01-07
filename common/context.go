package common

import (
	"context"
	"math/rand"
)

// Context : 应用程序上下文
type Context struct {
	Ctx               context.Context
	Rand              *rand.Rand
	Config            *Config
	Log               ILogger
	Node              INode
	ServerForClient   ITCPServer
	ServerForIntranet ITCPServer
	Login             ILogin
}
