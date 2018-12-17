# go-xserver
Go服务器框架（go-x.v2）

现在看 go-x ，制作时没有从使用者角度来做。因此重新写个。

## 代码框架

以 go-xserver 加载逻辑插件的方式，把框架层与逻辑层代码严格分离

从代码层看，主要有 3 个目录：

- common
- services
- internal

## services 目录

每个具体服务应用，都会是该目录下一个子文件夹

它会用 common 目录中的一些公共代码、接口

internal 目录在不在都不会影响 services 目录中代码的编译、运行

即 internal 目录可以对 services 目录不可见

## common 目录

通常为 interface 接口声明为主，细节实现可以放入 internal 目录

其实现，会被编译进 go-xserver 程序内

## internal 目录

框架层代码。可以对 services 目录不可见

## 例子参考

比如对 log 的封装代码：

- [common/log.go](common/log.go)
- [internal/log.go](internal/log.go)

## 将要实现的功能

- 框架层功能
    - 服务器组内互联
    - 服务器组内通信
    
- 逻辑层功能
    - 登陆服务
    - 网关服务
    - 大厅服务
    - 匹配服务
    - 房间服务
    - 管理服务

