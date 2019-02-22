# go-xserver

**go-xserver 是一个 Golang 服务器框架（go-x.v2）**

致力于实现 1 个高可用、高易用的 Golang 服务器框架

并以插件的方式，来丰富框架内容

## 编译

- 编译环境需要[翻墙设置](doc/编译-翻墙设置.md)

- 编译执行以下语句即可：

  ```shell
  ./make.sh
  ```

- Windows 10 下开发，请参考[在Win10中Linux环境搭建](doc/编译-在Win10中Linux环境搭建.md)


## 运行

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
- [服务器组内互联](doc/规范-服务器架构.md)

## 已完成服务器

- 管理服务器
- 登陆服务器

## 测试客户端

- [pyclient](https://github.com/fananchong/go-xclient/tree/master/pyclient)

## 正在制作

- 服务器组内通信
  - Client <-> Gateway <-> Other Node , 代码整理，完善
  - One Node <-> Gateway <-> Other Node , 待开始
- INode Send系列接口不干净
  - msg proto.Message 改为 msgdata []byte, userdata []byte 代替，方能支持使用方协议自由

## 将要实现的功能

- 框架层功能
    - *服务器组内通信*
    - 灰度更新
    - 服务器健康监测


- 逻辑层功能
    - 入口服务
    - 网关服务
    - 大厅服务
    - 匹配服务
    - 房间服务
