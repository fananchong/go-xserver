# 插件工程独立建库问题

目前 golang 尚不支持主程序、插件程序分别建库

详细请参考关注官方 issue ： https://github.com/golang/go/issues/27751

go-xserver 已经把缺省插件移至 https://github.com/fananchong/go-xserver-plugins.git

因此如果分别编译 go-xserver 、 go-xserver-plugins 出来的 go-xserver 程序是无法加载如 login.so 的


## 解决方法 1 （推荐）

1. 下载 go 源码 https://github.com/golang/go.git
2. 切最新稳定版
3. 注释 [go/src/runtime/plugin.go](https://github.com/golang/go/blob/50bd1c4d4eb4fac8ddeb5f063c099daccfb71b26/src/runtime/plugin.go) 文件中：
    ```go
	for _, pkghash := range md.pkghashes {
		if pkghash.linktimehash != *pkghash.runtimehash {
			md.bad = true
			return "", nil, "plugin was built with a different version of package " + pkghash.modulename
		}
	}
    ```
4. 编译 go 代码
5. 使用编译出来的 go 程序编译 go-xserver


## 解决方法 2

拷贝 [go-xserver-plugins](https://github.com/fananchong/go-xserver-plugins.git) 到 go-xserver 工程中编译

目前 [make.sh](make.sh) 、 [main.go](main.go) 中正是干了这件事

如果是直接 git clone go-xserver 的，是不需要做额外工作，即可完整编译的


## 解决方法 3

等待官方解决 [issue#27751](https://github.com/golang/go/issues/27751)
