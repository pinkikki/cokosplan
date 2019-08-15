package handler

import (
	"net/http"

	"github.com/justinas/alice"
	"go.uber.org/zap"
)

var All = alice.New(Recover, Logging)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging := zap.L().Named("LoggingMiddleware")
		logging.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
		next.ServeHTTP(w, r)
		logging.Info("response")
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging := zap.L().Named("RecoverMiddleware")
		defer func() {
			if err := recover(); err != nil {
				logging.Error("panic.", zap.Any("err", err))
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
