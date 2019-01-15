package components

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/db"
	uuid "github.com/satori/go.uuid"
)

// Login : 登陆模块
type Login struct {
	ctx              *common.Context
	verificationFunc common.FuncTypeAccountVerification
	allocType        []common.NodeType
	idgen            db.IDGen
	redisAtomic      db.RedisAtomic
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
	var accountID uint64
	var err common.LoginErrCode
	if defaultMode {
		accountID, err = login.loginByDefault(account, password)
	} else {
		accountID, err = login.verificationFunc(account, password, userdata)
	}
	if err != common.LoginSuccess {
		return "", nil, nil, nil, err
	}

	login.ctx.Log.Infof("account:%s, account id:%d", account, accountID)

	// 分配服务资源列表
	serverList, ok := login.selectServerList(account, login.allocType)
	if !ok {
		return "", nil, nil, nil, common.LoginSystemError
	}

	//生成 Token
	tokenObj := db.NewToken(login.ctx.Config.DbToken.Name, account)
	tokenObj.SetToken(uuid.NewV4().String())
	tokenObj.SetAccountID(accountID)
	if err := tokenObj.Save(); err != nil {
		login.ctx.Log.Errorln(err)
		return "", nil, nil, nil, common.LoginSystemError
	}
	return "", serverList.AddressList, serverList.PortList, serverList.TypeList, common.LoginSuccess
}

// Close : 关闭
func (login *Login) Close() {

}

func (login *Login) loginByDefault(account, password string) (uint64, common.LoginErrCode) {
	var accountID uint64
	accountObj := db.NewAccount(login.ctx.Config.DbAccount.Name, account)
	if err := accountObj.Load(); err != nil {
		if err != go_redis_orm.ERR_ISNOT_EXIST_KEY {
			return 0, common.LoginSystemError
		}
		accountID, err = login.idgen.New(db.IDGenTypeAccount)
		if err != nil {
			login.ctx.Log.Errorln(err)
			return 0, common.LoginSystemError
		}
		accountObj.SetID(accountID)
		accountObj.SetPasswd(password)
		if err = accountObj.Save(); err != nil {
			login.ctx.Log.Errorln(err)
			return 0, common.LoginSystemError
		}
	} else {
		accountID = accountObj.GetID()
	}
	return accountID, common.LoginSuccess
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
	err = go_redis_orm.CreateDB(
		login.ctx.Config.DbServer.Name,
		login.ctx.Config.DbServer.Addrs,
		login.ctx.Config.DbServer.Password,
		login.ctx.Config.DbServer.DBIndex)
	if err != nil {
		login.ctx.Log.Errorln(err)
		return false
	}

	login.idgen.Cli = go_redis_orm.GetDB(login.ctx.Config.DbAccount.Name)
	login.redisAtomic.Cli = go_redis_orm.GetDB(login.ctx.Config.DbServer.Name)
	return true
}

func (login *Login) selectServerList(account string, nodeType []common.NodeType) (dbObj *db.AccountServers, ok bool) {
LOOP:
	dbObj = &db.AccountServers{}
	for _, v := range nodeType {
		node := login.ctx.Node.GetNodeOne(v)
		if node == nil {
			login.ctx.Log.Errorln("no find server. type =", v)
			return
		}
		dbObj.UIDList = append(dbObj.UIDList, node.GetID())
		dbObj.TypeList = append(dbObj.TypeList, v)
		dbObj.AddressList = append(dbObj.AddressList, node.GetIP(common.IPOUTER))
		dbObj.PortList = append(dbObj.PortList, node.GetPort(int(common.PORTFORCLIENT)))
	}
	var data string
	var err error
	data, err = dbObj.Marshal()
	if err != nil {
		login.ctx.Log.Errorln(err)
		return
	}
	var ret string
	ret, err = login.redisAtomic.SetX(account, data, 86400)
	if err != nil {
		login.ctx.Log.Errorln(err)
		return
	}
	if ret != "" {
		dbObj.Unmarshal(ret)
		for _, id := range dbObj.UIDList {
			if login.ctx.Node.HaveNode(id) == false {
				if _, err = login.redisAtomic.DelX(account, data); err != nil {
					login.ctx.Log.Errorln(err)
					return
				}
				login.ctx.Log.Infoln("try again to get server list")
				goto LOOP
			}
		}
	}
	ok = true
	return
}
