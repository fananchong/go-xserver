# go-xserver

**go-xserver 是一个 Golang 服务器框架（go-x.v2）**

致力于实现 1 个高可用、高易用的 Golang 服务器框架

并以插件的方式，来丰富框架内容

## 编译

- [安装 golang 1.12+](https://golang.google.cn/dl/)
- [安装 docker](https://docs.docker.com/install/linux/docker-ce/centos/)
- 编译执行以下语句即可：

  ```shell
  ./make.sh
  ```

- 【非必须】 Windows 10 下开发，请参考[在 Win10 中 Linux 环境搭建](doc/WIKI-在Win10中Linux环境搭建.md)


## 运行

- 安装 Redis ，并修改 config/config.toml 相关配置

- All In One 例子
  ```shell
  ./make.sh start
  ./make.sh stop
  ```

- Run In WSL 例子
  ```shell
  ./wsl.sh start
  ./wsl.sh stop
  ```

   wsl 目前`监听同一个端口不报错`，详细请参考 issue ： https://github.com/Microsoft/WSL/issues/2915

   因此 wsl.sh 脚本中具体指定下 --network-port 参数



## WIKI

- [主体框架](doc/规范-代码框架.md)
- 配置模块
  - [框架层配置](doc/规范-配置文件_框架层.md)
  - [逻辑层配置](doc/规范-配置文件_逻辑层.md)
- [服务发现](doc/框架层功能-服务发现.md)
- [登陆模块](doc/框架层功能-登陆模块.md)
- [闲置连接处理](doc/框架层功能-闲置连接处理.md)
- [登出模块](doc/框架层功能-登出模块.md)
- [服务器组内互联](doc/规范-服务器架构.md)

## 测试客户端

- [pyclient](https://github.com/fananchong/go-xclient/tree/master/pyclient)


## 缺省插件

- [go-xserver-plugins](https://github.com/fananchong/go-xserver-plugins)
  - mgr
  - login
  - gateway


## v0.1

- 管理服务器
- 登陆服务器
- 网关服务器
  - 客户端消息中继
  - 服务器组内消息中继
- 大厅服务器
  - 获取角色列表（登录大厅服务）
  - 创建角色
  - 获取角色详细信息（进入游戏）
  - 登出游戏
  - 角色聊天（世界聊天、私聊）
  
  
## v0.2

- 参考 micro/go-micro 改造框架层代码
- 服务发现重做，参考 micro/go-micro 提炼 接口，并默认支持 mdns 


## 已知 ISSUE

- [插件工程独立建库问题](doc/ISSUE-插件工程独立建库问题.md)

## 将要实现的功能

- 框架层功能
    - 灰度更新
    - 服务器健康监测
- 逻辑层功能
    - 匹配服务
    - 房间服务
    - 压测工具
