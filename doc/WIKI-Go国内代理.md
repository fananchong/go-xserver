## go get

gopm ： https://github.com/gpmgo/gopm

- 安装 gopm

  ```sh
  go get -u github.com/gpmgo/gopm
  ```

- 下载包，使用 `gopm get` 代替 `go get`

  ```sh
  gopm get -u -v [包地址]
  ```

## go module

goproxy ： https://goproxy.io/

在执行诸如 `go build` 等命令前，先执行以下语句：

- Bash (MAC/Linux)
  ```sh
  export GOPROXY=https://goproxy.io
  ```

- PowerShell (Windows)
  ```sh
  $env:GOPROXY = "https://goproxy.io"
  ```
