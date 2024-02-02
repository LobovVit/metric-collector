package postgresql

import (
	"context"
	"time"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	dbctx, chancel := context.WithTimeout(ctx, time.Second*5)
	defer chancel()
	conn, err := pgx.Connect(dbctx, dsn)
	if err != nil {
		logger.Log.Error("Ошибка подключения к DB", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	return conn, nil
}
