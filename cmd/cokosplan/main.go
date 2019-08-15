package main

import (
	"net/http"

	"github.com/pinkikki/cokosplan/handler"
	"github.com/pinkikki/cokosplan/logging"
)

func main() {
	logging.Setting(logging.NewMode("debug"))
	http.ListenAndServe(":8080", handler.Routes())
}
