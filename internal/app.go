package internal

import (
	"math/rand"
	"net/http"
	_ "net/http/pprof" // 只使用 pprof 包的 init 函数
	"os"
	"os/signal"
	"plugin"
	"runtime"
	"syscall"
	"time"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/gotcp"
	"github.com/spf13/viper"
)

// App : 应用程序类
type App struct {
	signal chan os.Signal
	runner common.Plugin
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
	termination := false
	app.signal = make(chan os.Signal)
	signal.Notify(app.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
		syscall.SIGPIPE)

	// 初始化随机种子
	common.XRAND = rand.New(rand.NewSource(time.Now().UnixNano()))

	// 加载配置
	if err := loadConfig(); err != nil {
		return
	}

	// 初始化 Log
	app.initLog()
	defer common.XLOG.Flush()

	// 加载插件
	if err := app.initPlugin(); err != nil {
		common.XLOG.Errorln(err)
		return
	}

	// 初始化性能分析工具
	app.initPprof()

	// 连接 redis mgr 数据库
	if err := initRedis(); err != nil {
		common.XLOG.Errorln(err)
		return
	}

	// 应用程序正式运行
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
	if app.initNode(nodeType) == false {
		app.close()
		return
	}
	app.runner.RegisterCallBack()
	if startNode(common.XNODE) == false {
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

func (app *App) close() {
	close(app.signal)
}

func (app *App) initNode(nodeType common.NodeType) bool {
	node := newNode(nodeType)
	if node == nil {
		app.close()
		return false
	}
	common.XNODE = node
	return true
}

func (app *App) initLog() {
	common.XLOG = newGLogger()
	logDir := common.XCONFIG.Common.LogDir
	if logDir != "" {
		os.MkdirAll(logDir, os.ModePerm)
	}
	common.XLOG.SetLogDir(logDir)
	common.XLOG.SetLogLevel(common.XCONFIG.Common.LogLevel)
	gotcp.SetLogger(common.XLOG)
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

func initRedis() error {
	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)
	return go_redis_orm.CreateDB(
		common.XCONFIG.DbMgr.Name,
		common.XCONFIG.DbMgr.Addrs,
		common.XCONFIG.DbMgr.Password,
		common.XCONFIG.DbMgr.DBIndex)
}
