package common

// IGateway 接口暴露框架层网关模块的使用方法
// 完整的网关，由框架层登陆模块、逻辑层交互模块共同完成

// 网关模块框架层负责的工作：
//    1. 提供服务器组内节点互连、接入验证
//    2. 提供 Client(s) <-> Gateway <-> Node 路径的消息转发
//    3. 提供 Node <-> Gateway <-> Other Node(s) 路径的消息转发
//    4. 框架层会接管几乎所有的网关的功能（网关会很纯粹，主要就是消息转发）

// 网关模块逻辑层负责的工作：
//    1. 自定义客户端交互协议（如使用 TCP 、 HTTP ； 如使用 proto 消息、 struct 消息 等等）
//    2. 自定义加解密算法

// 经过网关的消息规则，请参见：[规范-消息号.md](go-xserver/doc/规范-消息号.md)

// FuncTypeEncode : 数据加密函数声明
type FuncTypeEncode func(data []byte) []byte

// FuncTypeDecode : 数据解密函数声明
type FuncTypeDecode func(data []byte) []byte

// FuncTypeSendToClient : 发送客户端数据函数声明
type FuncTypeSendToClient func(account string, cmd uint64, data []byte) bool

// IGateway : 网关模块接口
type IGateway interface {
	VerifyToken(account, token string) uint32                             // 令牌验证。返回值： 0 成功；1 令牌错误； 2 系统错误
	OnRecvFromClient(account string, cmd uint32, data []byte) (done bool) // 可自定义客户端交互协议。done 为 true ，表示框架层接管处理该消息
	RegisterSendToClient(f FuncTypeSendToClient)                          // 可自定义客户端交互协议
	GetSendToClient() FuncTypeSendToClient                                // 可自定义客户端交互协议
	RegisterEncodeFunc(f FuncTypeEncode)                                  // 可自定义加解密算法
	RegisterDecodeFunc(f FuncTypeDecode)                                  // 可自定义加解密算法
}
