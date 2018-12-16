package internal

import (
	"fmt"
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
		fmt.Println(err)
		return
	}

	// init log
	app.initLog()

	// load plugin
	if err := app.initPlugin(); err != nil {
		return
	}

	app.initProf()

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

func (app *App) initLog() {
	common.XLOG = common.NewGLogger()
	logDir := viper.GetString("common.LogDir")
	if logDir != "" {
		os.MkdirAll(logDir, os.ModeDir)
	}
	common.XLOG.SetLogDir(logDir)
	common.XLOG.SetLogLevel(viper.GetInt("common.LogLevel"))
}

func (app *App) initPlugin() error {
	appName := viper.GetString("app")
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

func (app *App) initProf() {
	addr := viper.GetString("common.Pprof")
	if addr != "" {
		go func() {
			common.XLOG.Println("pprof listen :", addr)
			http.ListenAndServe(addr, nil)
		}()
	}
}

func (app *App) close() {
	close(app.signal)
}

func (app *App) onAppReady() {
	if app.runner.Init() == false {
		app.close()
		return
	}
	go func() {
		if app.runner.Start() == false {
			app.close()
		}
	}()
}

func (app *App) onAppShutDown() {
	app.runner.Close()
}
