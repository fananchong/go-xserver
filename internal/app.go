package internal

import (
	"math/rand"
	"net/http"
	_ "net/http/pprof" // only need init function in pprof
	"os"
	"os/signal"
	"plugin"
	"runtime"
	"syscall"
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/spf13/viper"
)

// App : An application class
type App struct {
	signal chan os.Signal
	runner common.Plugin
}

// NewApp : Constructor function of class App
func NewApp() *App {
	app := &App{}
	return app
}

// Run : App instance launch
func (app *App) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	termination := false
	app.signal = make(chan os.Signal)
	signal.Notify(app.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
		syscall.SIGPIPE)

	common.XRAND = rand.New(rand.NewSource(time.Now().UnixNano()))

	// load config
	if err := loadConfig(); err != nil {
		return
	}

	// init log
	app.initLog()
	defer common.XLOG.Flush()

	// load plugin
	if err := app.initPlugin(); err != nil {
		common.XLOG.Errorln(err)
		return
	}

	// init pprof
	app.initPprof()

	// running
	app.onAppReady()
	for !termination {
		select {
		case sig := <-app.signal:
			switch sig {
			case syscall.SIGPIPE:
			default:
				termination = true
			}
			common.XLOG.Infoln("[app] recive signal. signal no =", sig)
		}
	}
	app.onAppShutDown()
}

func (app *App) onAppReady() {
	var nodeType common.NodeType
	var ok bool
	if nodeType, ok = app.runner.Init(); !ok {
		app.close()
		return
	}
	if app.initNode(nodeType) {
		go func() {
			if app.runner.Start() == false {
				app.close()
			}
		}()
	}
}

func (app *App) onAppShutDown() {
	app.runner.Close()
}

func (app *App) close() {
	close(app.signal)
}

func (app *App) initNode(nodeType common.NodeType) bool {
	node := NewNode(nodeType)
	if node.Init() == false {
		app.close()
		return false
	}
	common.XNODE = node
	return true
}

func (app *App) initLog() {
	common.XLOG = NewGLogger()
	logDir := common.XCONFIG.Common.LogDir
	if logDir != "" {
		os.MkdirAll(logDir, os.ModePerm)
	}
	common.XLOG.SetLogDir(logDir)
	common.XLOG.SetLogLevel(common.XCONFIG.Common.LogLevel)
}

func (app *App) initPlugin() error {
	appName := viper.GetString("app")
	if appName == "" {
		printUsage()
		os.Exit(1)
	}
	p, err := plugin.Open(appName + ".so")
	if err != nil {
		return err
	}
	obj, err := p.Lookup("PluginObj")
	if err != nil {
		return err
	}
	app.runner = *obj.(*common.Plugin)
	return nil
}

func (app *App) initPprof() {
	addr := common.XCONFIG.Common.Pprof
	if addr != "" {
		go func() {
			common.XLOG.Println("pprof listen :", addr)
			http.ListenAndServe(addr, nil)
		}()
	}
}
