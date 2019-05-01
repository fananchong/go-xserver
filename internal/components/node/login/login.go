package nodelogin

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/internal/components/misc"
	nodecommon "github.com/fananchong/go-xserver/internal/components/node/common"
	nodenormal "github.com/fananchong/go-xserver/internal/components/node/normal"
	"github.com/fananchong/go-xserver/internal/db"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utils"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

// Login : 登陆模块
type Login struct {
	*nodenormal.Normal
	ctx              *common.Context
	verificationFunc common.FuncTypeAccountVerification
	allocServerType  []common.NodeType
	serverRedis      db.RedisAtomic
}

// NewLogin : 实例化登陆模块
func NewLogin(ctx *common.Context) *Login {
	login := &Login{
		ctx: ctx,
	}
	login.ctx.ILogin = login
	return login
}

// Start : 启动
func (login *Login) Start() bool {
	pluginType := misc.GetPluginType(login.ctx)
	if pluginType == common.Login {
		login.Normal = login.ctx.INode.(*nodenormal.Normal)
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
	login.allocServerType = append(login.allocServerType, types...)
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
	addressList, portList, ids, ok := login.selectServerList(account, login.allocServerType)
	if !ok {
		return "", nil, nil, nil, common.LoginSystemError
	}

	//生成 Token
	tempID := uuid.NewV4().String()
	tokenObj := db.NewToken(login.ctx.Config.DbToken.Name, account)
	tokenObj.Expire(7 * 86400) // 令牌过期时间 7 天
	to := tokenObj.GetToken(true)
	to.Token = tempID
	to.AllocServers = make(map[uint32]*protocol.SERVER_ID)
	for i := 0; i < len(login.allocServerType); i++ {
		to.AllocServers[uint32(login.allocServerType[i])] = ids[i]
	}
	if err := tokenObj.Save(); err != nil {
		login.ctx.Errorln(err, "account:", account)
		return "", nil, nil, nil, common.LoginSystemError
	}
	return tempID, addressList, portList, login.allocServerType, common.LoginSuccess
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
			login.ctx.Errorln(err, "account:", account)
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
	// db account
	err := go_redis_orm.CreateDB(
		login.ctx.Config.DbAccount.Name,
		login.ctx.Config.DbAccount.Addrs,
		login.ctx.Config.DbAccount.Password,
		login.ctx.Config.DbAccount.DBIndex)
	if err != nil {
		login.ctx.Errorln(err)
		return false
	}

	// db token
	err = go_redis_orm.CreateDB(
		login.ctx.Config.DbToken.Name,
		login.ctx.Config.DbToken.Addrs,
		login.ctx.Config.DbToken.Password,
		login.ctx.Config.DbToken.DBIndex)
	if err != nil {
		login.ctx.Errorln(err)
		return false
	}

	// db server
	c, err := redis.Dial("tcp", login.ctx.Config.DbServer.Addrs[0])
	if err != nil {
		login.ctx.Errorln(err)
		return false
	}
	if login.ctx.Config.DbServer.Password != "" {
		if _, err := c.Do("AUTH", login.ctx.Config.DbServer.Password); err != nil {
			login.ctx.Errorln(err)
			c.Close()
			return false
		}
	}
	if login.ctx.Config.DbServer.DBIndex != 0 {
		if _, err := c.Do("SELECT", login.ctx.Config.DbServer.DBIndex); err != nil {
			login.ctx.Errorln(err)
			c.Close()
			return false
		}
	}
	login.serverRedis.Cli = c
	return true
}

func (login *Login) selectServerList(account string, nodeType []common.NodeType) (addressList []string, portList []int32, serverIDs []*protocol.SERVER_ID, ok bool) {
	for i := 0; i < len(nodeType); i++ {
		dbObj, have := login.selectServer(account, nodeType[i])
		if !have {
			return
		}
		addressList = append(addressList, dbObj.Address)
		portList = append(portList, dbObj.Port)
		serverIDs = append(serverIDs, dbObj.ServerID)
	}
	ok = true
	return
}

func (login *Login) selectServer(account string, nodeType common.NodeType) (dbObj *db.AccountServer, ok bool) {
LOOP:
	dbObj = &db.AccountServer{}
	login.PrintNodeInfo(login.ctx, nodeType)
	node := login.GetNodeOne(nodeType)
	if node == nil {
		login.ctx.Errorln("Did not find the server. type:", nodeType, "account:", account)
		return
	}
	nodeID := node.GetID()
	dbObj.ServerID = nodecommon.NodeID2ServerID(nodeID)
	dbObj.Type = nodeType
	dbObj.Address = node.GetIP(utils.IPOUTER)
	dbObj.Port = node.GetPort(int(utils.PORTFORCLIENT))

	var data string
	var err error
	data, err = dbObj.Marshal()
	if err != nil {
		login.ctx.Errorln(err, "account:", account)
		return
	}
	login.ctx.Infoln("account:", account, "server:", data)
	var ret string
	key := db.GetKeyAllocServer(uint32(nodeType), account)
	ret, err = login.serverRedis.SetGetX(key, data, 365*86400) // 设置账号分配的服务器资源信息，过期时间 1 年
	if err != nil {
		login.ctx.Errorln(err, "account:", account)
		return
	}
	if ret != "" {
		dbObj.Unmarshal(ret)
		if login.HaveNode(nodecommon.ServerID2NodeID(dbObj.ServerID)) == false {
			if _, err = login.serverRedis.DelX(key, ret); err != nil {
				login.ctx.Errorln(err, "account:", account)
				return
			}
			login.ctx.Infoln("Try again to get one server, type:", nodeType, "account:", account)
			goto LOOP
		}
	}
	ok = true
	return
}
