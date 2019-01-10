package common

// ILogin 接口暴露框架层登陆模块的使用方法
// 完整的登陆过程，由框架层登陆模块、逻辑层交互模块共同完成

// 登陆模块框架层负责的工作：
//    1. 提供缺省账号验证
//    2. 提供账号正常登陆
//        a. 提供账号正常上下线，不会回档等错误
//        b. 提供同账号多端同时登陆，不会异常
//        c. 提供账号服务器资源分配、比如分配 Gateway 资源
//    3. 提供自定义验证接口

// 登陆模块逻辑层负责的工作：
//    1. 自定义协议
//    2. 自定义客户端交互流程
//    3. 自定义账号验证过程

// go-xserver/services/login :
//    1. 一个缺省的 Login Server 实现例子
//    2. 主要展示逻辑层如何与框架层交互
//    3. 可以参考实现自己的项目要求的 Login Server

// LoginErrCode : 登陆错误
type LoginErrCode int

const (
	// LoginSuccess : 登陆成功
	LoginSuccess LoginErrCode = iota

	// VerifyFail : 验证错误
	VerifyFail

	// SystemError : 系统错误
	SystemError
)

// FuncTypeAccountVerification : 账号验证函数声明
type FuncTypeAccountVerification func(account, password string, userdata []byte) LoginErrCode

// ILogin : 登陆模块接口
type ILogin interface {
	Start() bool
	RegisterCustomAccountVerification(f FuncTypeAccountVerification)
	Login(account, password string, defaultMode bool, userdata []byte) (token, address string, port int32, errcode LoginErrCode)
	Close()
}
