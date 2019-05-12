set CUR_DIR=%~dp0
set GOBIN=%GOPATH%\bin
mkdir %GOBIN%
cd /d c:
cd %GOBIN%
go get -u github.com/gpmgo/gopm
gopm bin -u -v github.com/ramya-rao-a/go-outline
gopm bin -u -v github.com/acroca/go-symbols
gopm bin -u -v github.com/mdempsky/gocode
gopm bin -u -v github.com/zmb3/gogetdoc
gopm bin -u -v github.com/fatih/gomodifytags
gopm bin -u -v golang.org/x/tools/cmd/gorename
gopm bin -u -v golang.org/x/tools/cmd/goimports
gopm bin -u -v golang.org/x/tools/cmd/guru
gopm bin -u -v github.com/josharian/impl
gopm bin -u -v github.com/haya14busa/goplay/cmd/goplay
gopm bin -u -v github.com/uudashr/gopkgs/cmd/gopkgs
gopm bin -u -v github.com/davidrjenni/reftools/cmd/fillstruct
gopm bin -u -v github.com/cweill/gotests/gotests
gopm bin -u -v golang.org/x/tools/cmd/gopls
gopm bin -u -v github.com/sqs/goreturns
gopm bin -u -v golang.org/x/lint/golint
go get -u -v github.com/alecthomas/gometalinter
gometalinter --install
go get -u -v github.com/rogpeppe/godef
gopm bin -u -v golang.org/x/tools/cmd/godoc

cd %CUR_DIR%