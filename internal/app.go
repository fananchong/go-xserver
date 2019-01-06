package internal

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/components"
)

// App : 应用程序类
type App struct {
	signal     chan os.Signal
	components []utils.IComponent
}

// NewApp : 应用程序类的构造函数
func NewApp() *App {
	app := &App{}
	return app
}

// Run : 启动应用程序
func (app *App) Run() {
	// 设置最大线程数
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 设置信号量
	app.signal = make(chan os.Signal)
	signal.Notify(app.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
		syscall.SIGPIPE)

	// 注册组件
	app.components = []utils.IComponent{
		&components.Rand{},
		&components.Config{},
		&components.Log{},
		&components.Pprof{},
		&components.Redis{},
		//&components.TCPServer{},
		&components.Node{},
		&components.Login{},
		&components.Plugin{}, // 最后 1 个，为 Plugin
	}

	// 应用程序正式运行
	app.onAppReady()
	defer app.onAppShutDown()
	for {
		select {
		case sig := <-app.signal:
			switch sig {
			case syscall.SIGPIPE:
			default:
				common.XLOG.Infoln("[app] recive signal:", sig)
				return
			}
		}
	}
}

func (app *App) onAppReady() {
	for i := 0; i < len(app.components); i++ {
		c := app.components[i]
		if !c.Start() {
			app.close()
			return
		}
	}
}

func (app *App) onAppShutDown() {
	for i := len(app.components) - 1; i >= 0; i-- {
		c := app.components[i]
		c.Close()
	}
}

func (app *App) close() {
	close(app.signal)
}
