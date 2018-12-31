## 安装 WSL

  - 控制面板->程序和功能->启用或关闭Windows功能->勾选 适用于Linux的Windows子系统

  - 打开应用商城搜索“WSL”，可根据自己需求选择安装一个或多个Linux系统（推荐：Ubuntu 18）

  - Shift + 鼠标右键，点击`在此处打开 Linux Shell(L)` ,即可打开 WSL

## 安装 VSCode

主要步骤：

  - 安装 VSCode
  - 设置代理
  - 安装 go 插件
  - WSL 嵌入 VSCode

详细请参考：https://blog.csdn.net/u013272009/article/details/84971807


## 安装 Docker for Windows（非必须）

  - 按照 https://www.docker.com/products/docker-desktop 中提示安装之

  - 打开 Docker for Windows - General ，勾选`Expose daemon on tcp://localhost:2375 without TLS`

  - 打开 Docker for Windows - Shared Drives ，勾选所有盘符

  - 打开 WSL ，安装 docker.io
    - 请参考 https://blog.csdn.net/u013272009/article/details/81221661


  - 修正挂接目录问题
    - 请参考 https://blog.csdn.net/u013272009/article/details/81222689


  - Docker 容器开启失败
    - 请参考 https://blog.csdn.net/u013272009/article/details/85002613
