package postgresql

import (
	"context"
	"time"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	logger.Log.Info("YA DSN", zap.String("dsn", dsn))
	dbctx, chancel := context.WithTimeout(ctx, time.Second*5)
	defer chancel()
	conConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		logger.Log.Error("Parse config failed", zap.Error(err))
		return nil, err
	}
	conn, err := pgx.ConnectConfig(dbctx, conConfig)
	if err != nil {
		logger.Log.Error("Ошибка подключения к DB", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	return conn, nil
}
