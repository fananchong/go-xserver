package common

// IGateway 接口暴露框架层网关模块的使用方法
// 完整的网关，由框架层登陆模块、逻辑层交互模块共同完成

// 网关模块框架层负责的工作：
//    1. 提供服务器组内节点互连、接入验证
//    2. 提供 Client(s) <-> Gateway <-> Node 路径的消息转发
//    3. 提供 Node <-> Gateway <-> Other Node(s) 路径的消息转发
//    4. 框架层会接管几乎所有的网关的功能（网关会很纯粹，主要就是消息转发）

// 网关模块逻辑层负责的工作：
//    1. 自定义加解密算法

// 经过网关的消息规则，请参见：[规范-消息号.md](go-xserver/doc/规范-消息号.md)

// FuncTypeEncode : 数据加密函数声明
type FuncTypeEncode func(data []byte) []byte

// FuncTypeDecode : 数据解密函数声明
type FuncTypeDecode func(data []byte) []byte

// IGateway : 网关模块接口
type IGateway interface {
	RegisterEncodeFunc(f FuncTypeEncode)
	RegisterDecodeFunc(f FuncTypeDecode)
}
