package internal

import (
	"os"

	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/components"
)

// App : 应用程序类
type App struct {
	components []utils.IComponent
}

// NewApp : 应用程序类的构造函数
func NewApp() *App {
	app := &App{}
	return app
}

// Run : 启动应用程序
func (app *App) Run() {
	// 注册组件
	app.components = []utils.IComponent{
		&components.Rand{},
		&components.Config{},
		&components.Log{},
		&components.Pprof{},
		&components.Redis{},
		&components.TCPServer{},
		&components.Node{},
		&components.Login{},
		&components.Plugin{}, // 必须倒数第 2 个为 Plugin
		&components.Signal{}, // 必须最后 1 个为 Signal
	}

	// 应用程序正式运行
	defer app.onAppShutDown()
	app.onAppReady()
}

func (app *App) onAppReady() {
	for i := 0; i < len(app.components); i++ {
		c := app.components[i]
		if !c.Start() {
			os.Exit(-1)
		}
	}
}

func (app *App) onAppShutDown() {
	for i := len(app.components) - 1; i >= 0; i-- {
		c := app.components[i]
		c.Close()
	}
}
