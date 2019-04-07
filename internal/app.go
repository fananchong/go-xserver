package internal

import (
	"os"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/components"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodegateway "github.com/fananchong/go-xserver/internal/components/node/gateway"
	nodelogin "github.com/fananchong/go-xserver/internal/components/node/login"
	nodemgr "github.com/fananchong/go-xserver/internal/components/node/mgr"
	nodenormal "github.com/fananchong/go-xserver/internal/components/node/normal"
)

// App : 应用程序类
type App struct {
	ctx        *common.Context
	components []utils.IComponent
}

// NewApp : 应用程序类的构造函数
func NewApp() *App {
	app := &App{
		ctx: &common.Context{Ctx: misc.CreateContext()},
	}
	return app
}

// Run : 启动应用程序
func (app *App) Run() {
	// 注册组件
	app.components = []utils.IComponent{
		components.NewRand(app.ctx),
		components.NewConfig(app.ctx),
		components.NewLog(app.ctx),
		components.NewPprof(app.ctx),
		components.NewPlugin(app.ctx),
		components.NewRedis(app.ctx),
		components.NewRole2Account(app.ctx),
		components.NewTCPServer(app.ctx),
		nodenormal.NewNormal(app.ctx),
		nodemgr.NewMgr(app.ctx),
		nodelogin.NewLogin(app.ctx),
		nodegateway.NewGateway(app.ctx),
		components.NewSignal(app.ctx),
	}

	// 应用程序正式运行
	defer app.onAppShutDown()
	app.onAppReady()
}

func (app *App) onAppReady() {
	misc.SetComponentCount(app.ctx.Ctx, len(app.components))
	for i := 0; i < len(app.components); i++ {
		misc.OneComponentOK(app.ctx.Ctx)
		c := app.components[i]
		if !c.Start() {
			os.Exit(1)
		}
	}
}

func (app *App) onAppShutDown() {
	for i := len(app.components) - 1; i >= 0; i-- {
		c := app.components[i]
		c.Close()
	}
}
