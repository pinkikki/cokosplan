package handler

import (
	"go.uber.org/zap"
)

type Context struct {
	Logger      *zap.Logger
	SugaredLogger *zap.SugaredLogger
}
