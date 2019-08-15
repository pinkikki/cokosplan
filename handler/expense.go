package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pinkikki/cokosplan/port"
	"go.uber.org/zap"
)

type ExpenseRouter struct {
	context *Context
}

func NewExpenseRouter() *ExpenseRouter {
	return &ExpenseRouter{context: &Context{SugaredLogger: zap.L().Sugar().Named("ExpenseRouter")}}
}

func (c *ExpenseRouter) Routes(router *mux.Router) {
	r := router.PathPrefix("/expenses").Subrouter()
	r.Methods("GET").Handler(All.Then(list(c.context)))
	r.Methods("POST").Handler(All.Then(save(c.context)))
}

func list(ctx *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event := r.URL.Query().Get("event")
		fromPaymentDate := r.URL.Query().Get("fromPaymentDate")
		toPaymentDate := r.URL.Query().Get("toPaymentDate")

		result, err := port.List(port.ExpenseCriteria{Event: event, FromPaymentDate: fromPaymentDate, ToPaymentDate: toPaymentDate})
		if err != nil {
			ctx.SugaredLogger.Errorf("Failed to get expenses. err=%v.", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		data, err := json.Marshal(result)
		if err != nil {
			ctx.SugaredLogger.Errorf("Failed to encode as json. err=%v. result=%v", err, result)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Write([]byte(data))
	})
}

func save(ctx *Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var e port.Expense
		json.Unmarshal(buf.Bytes(), &e)

		port.Save(e)

		w.WriteHeader(http.StatusCreated)
	})
}
