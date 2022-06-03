package repo

import (
	"context"
	"errors"
	"notification/pkg/logger"
	"notification/pkg/postgres"
)

var l logger.Interface

func LogSQL(ctx context.Context, sql string, args ...interface{}) {
	if log := GetLogger(ctx); log != nil {
		log.Infoln(sql, args)
	}
}

func GetLogger(ctx context.Context) logger.Interface {
	if l != nil {
		return l
	}
	if i := ctx.Value("logger"); i != nil {
		if log, ok := i.(logger.Interface); ok {
			l = log
			return log
		}
	}
	return nil

}

func LogCacheError(ctx context.Context, err error) {
	if log := GetLogger(ctx); log != nil {
		if errors.Is(err, postgres.ErrNotSetCache) {
			return
		}
		log.Errorln(err)
	}
}
