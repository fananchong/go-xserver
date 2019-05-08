//go:generate sh -c "if [ ! -d default_plugins ]; then git clone https://github.com/fananchong/go-xserver-plugins.git && mv go-xserver-plugins default_plugins; fi"
//go:generate sh -c "if [ -d default_plugins ]; then cd default_plugins && git fetch --all && git reset --hard origin/master && git pull; fi"
//go:generate sh -c "cd default_plugins && sed -i 's#go-xserver-plugins#go-xserver/default_plugins#g' `grep 'go-xserver-plugins' -rl . --include *.go`"
//go:generate sh -c "cd default_plugins && rm -rf go.*"
package main

import "github.com/fananchong/go-xserver/internal"

func main() {
	app := internal.NewApp()
	app.Run()
}
