package main

import (
	"fmt"
	"net/http"
)

func (login *Login) msgLogin(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello world")
}
