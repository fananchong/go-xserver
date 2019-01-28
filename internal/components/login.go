package components

import (
	"fmt"

	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

// Login : 登陆模块
type Login struct {
	ctx              *common.Context
	verificationFunc common.FuncTypeAccountVerification
	allocType        []common.NodeType
	idgen            db.IDGen
	serverRedis      db.RedisAtomic
}

// NewLogin : 实例化登陆模块
func NewLogin(ctx *common.Context) *Login {
	login := &Login{
		ctx: ctx,
	}
	login.ctx.Login = login
	return login
}

// Start : 启动
func (login *Login) Start() bool {
	if getPluginType(login.ctx) == common.Login {
		if !login.initRedis() {
			return false
		}
	}
	OneComponentOK(login.ctx.Ctx)
	return true
}

// RegisterCustomAccountVerification : 注册自定义账号验证处理
func (login *Login) RegisterCustomAccountVerification(f common.FuncTypeAccountVerification) {
	login.verificationFunc = f
}

// RegisterAllocationNodeType : 注册要分配的服务器资源类型
func (login *Login) RegisterAllocationNodeType(types []common.NodeType) {
	login.allocType = append(login.allocType, types...)
}

// Login : 登陆处理
func (login *Login) Login(account, password string, defaultMode bool, userdata []byte) (string,
	[]string, []int32, []common.NodeType, common.LoginErrCode) {

	//账号验证
	var err common.LoginErrCode
	if defaultMode {
		err = login.loginByDefault(account, password)
	} else {
		err = login.verificationFunc(account, password, userdata)
	}
	if err != common.LoginSuccess {
		return "", nil, nil, nil, err
	}

	// 分配服务资源列表
	addressList, portList, ok := login.selectServerList(account, login.allocType)
	if !ok {
		return "", nil, nil, nil, common.LoginSystemError
	}

	//生成 Token
	tokenObj := db.NewToken(login.ctx.Config.DbToken.Name, account)
	tokenObj.SetToken(uuid.NewV4().String())
	if err := tokenObj.Save(); err != nil {
		login.ctx.Log.Errorln(err, "account:", account)
		return "", nil, nil, nil, common.LoginSystemError
	}
	return tokenObj.GetToken(), addressList, portList, login.allocType, common.LoginSuccess
}

// Close : 关闭
func (login *Login) Close() {
	if login.serverRedis.Cli != nil {
		login.serverRedis.Cli.Close()
		login.serverRedis.Cli = nil
	}
}

func (login *Login) loginByDefault(account, password string) common.LoginErrCode {
	accountObj := db.NewAccount(login.ctx.Config.DbAccount.Name, account)
	if err := accountObj.Load(); err != nil {
		// 新建账号
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return common.LoginSystemError
		}
		accountObj.SetPasswd(password)
		if err = accountObj.Save(); err != nil {
			login.ctx.Log.Errorln(err, "account:", account)
			return common.LoginSystemError
		}
	} else {
		// 验证密码
		if accountObj.GetPasswd() != password {
			return common.LoginVerifyFail
		}
	}
	return common.LoginSuccess
}

func (login *Login) initRedis() bool {
	go_redis_orm.SetNewRedisHandler(go_redis_orm.NewDefaultRedisClient)

	// db account
	err := go_redis_orm.CreateDB(
		login.ctx.Config.DbAccount.Name,
		login.ctx.Config.DbAccount.Addrs,
		login.ctx.Config.DbAccount.Password,
		login.ctx.Config.DbAccount.DBIndex)
	if err != nil {
		login.ctx.Log.Errorln(err)
		return false
	}

	// db token
	err = go_redis_orm.CreateDB(
		login.ctx.Config.DbToken.Name,
		login.ctx.Config.DbToken.Addrs,
		login.ctx.Config.DbToken.Password,
		login.ctx.Config.DbToken.DBIndex)
	if err != nil {
		login.ctx.Log.Errorln(err)
		return false
	}

	// db server
	c, err := redis.Dial("tcp", login.ctx.Config.DbServer.Addrs[0])
	if err != nil {
		login.ctx.Log.Errorln(err)
		return false
	}
	if login.ctx.Config.DbServer.Password != "" {
		if _, err := c.Do("AUTH", login.ctx.Config.DbServer.Password); err != nil {
			login.ctx.Log.Errorln(err)
			c.Close()
			return false
		}
	}
	if login.ctx.Config.DbToken.DBIndex != 0 {
		if _, err := c.Do("SELECT", login.ctx.Config.DbServer.DBIndex); err != nil {
			login.ctx.Log.Errorln(err)
			c.Close()
			return false
		}
	}

	login.idgen.Cli = go_redis_orm.GetDB(login.ctx.Config.DbAccount.Name)
	login.serverRedis.Cli = c
	return true
}

func (login *Login) selectServerList(account string, nodeType []common.NodeType) (addressList []string, portList []int32, ok bool) {
	for _, v := range nodeType {
		dbObj, have := login.selectServer(account, v)
		if !have {
			return
		}
		addressList = append(addressList, dbObj.Address)
		portList = append(portList, dbObj.Port)
	}
	ok = true
	return
}

func (login *Login) selectServer(account string, nodeType common.NodeType) (dbObj *db.AccountServer, ok bool) {
LOOP:
	dbObj = &db.AccountServer{}
	node := login.ctx.Node.GetNodeOne(nodeType)
	if node == nil {
		login.ctx.Log.Errorln("no find server. type:", nodeType, "account:", account)
		return
	}
	dbObj.NodeID = node.GetID()
	dbObj.Type = nodeType
	dbObj.Address = node.GetIP(common.IPOUTER)
	dbObj.Port = node.GetPort(int(common.PORTFORCLIENT))

	var data string
	var err error
	data, err = dbObj.Marshal()
	if err != nil {
		login.ctx.Log.Errorln(err, "account:", account)
		return
	}
	login.ctx.Log.Infoln("account:", account, "server:", data)
	var ret string
	key := fmt.Sprintf("srv%d_%s", nodeType, account)
	ret, err = login.serverRedis.SetX(key, data, 365*86400)
	if err != nil {
		login.ctx.Log.Errorln(err, "account:", account)
		return
	}
	if ret != "" {
		dbObj.Unmarshal(ret)
		if login.ctx.Node.HaveNode(dbObj.NodeID) == false {
			if _, err = login.serverRedis.DelX(key, ret); err != nil {
				login.ctx.Log.Errorln(err, "account:", account)
				return
			}
			login.ctx.Log.Infoln("try again to get server, type:", nodeType, "account:", account)
			goto LOOP
		}
	}
	ok = true
	return
}
