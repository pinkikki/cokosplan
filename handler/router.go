package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	Routes(r *mux.Router)
}

func Routes() *mux.Router {
	routers := []Router{
		NewExpenseRouter(),
	}
	router := mux.NewRouter()
	for _, r := range routers {
		r.Routes(router)
	}

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// alice/gorillaを使用しないでmiddlewareを使用
	// http.Handle("/simple", handler.Recover(handler.Logging(app)))

	return router
}
