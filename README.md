# go-xserver

**go-xserver 是一个 Golang 服务器框架（go-x.v2）**

致力于实现 1 个高可用、高易用的 Golang 服务器框架

并以插件的方式，来丰富框架内容

## 编译

- [安装 golang 1.11+](https://golang.google.cn/dl/)
- [安装 docker](https://docs.docker.com/install/linux/docker-ce/centos/)
- 编译执行以下语句即可：

  ```shell
  ./make.sh
  ```

- 【非必须】 Windows 10 下开发，请参考[在 Win10 中 Linux 环境搭建](doc/编译-在Win10中Linux环境搭建.md)


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



## 已完成功能

- [主体框架](doc/规范-代码框架.md)
- [配置模块](doc/规范-配置文件.md)
- [服务发现](doc/框架层功能-服务发现.md)
- [登陆模块](doc/框架层功能-登陆模块.md)
- [闲置连接处理](doc/框架层功能-闲置连接处理.md)
- [登出模块](doc/框架层功能-登出模块.md)
- [服务器组内互联](doc/规范-服务器架构.md)

## 已完成服务器

- 管理服务器
- 登陆服务器
- 网关服务器
  - 客户端消息中继
- 大厅服务器
  - 获取角色列表（登录大厅服务）
  - 创建角色
  - 获取角色详细信息（进入游戏）
  - 登出游戏
  - 角色聊天（世界聊天、私聊）

## 测试客户端

- [pyclient](https://github.com/fananchong/go-xclient/tree/master/pyclient)

## 正在制作

- 服务器组内通信
  - 服务节点到服务节点间通信（ 制作匹配服时实现 ）
- 匹配服务器
  - 匹配功能
- 大厅服务器
  - 匹配功能
- 房间服务器
  - 匹配完毕进入房间
- 网关服务器
  - 消息加解密
- 代码整理，加强代码可读性
- 代码优化，参考 goim ，优化代码


## 已知 BUG

已知待修正的 BUG ，记录之，空闲时修正

- 服务发现
  - 依次 1. MgrServer 失效； 2. 某节点失效重启； 3. MgrServer 重启；这种情况，其他节点会自动重连该节点，但是NodeID是旧的！

## 将要实现的功能

- 框架层功能
    - 灰度更新
    - 服务器健康监测
- 逻辑层功能
    - 入口服务
    - 匹配服务
    - 房间服务
    - 压测工具
