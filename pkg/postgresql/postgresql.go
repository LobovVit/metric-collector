package postgresql

import (
	"context"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	logger.Log.Info("YA DSN", zap.String("dsn", dsn))
	//dbctx, chancel := context.WithTimeout(ctx, time.Second*5)
	//defer chancel()
	conConfig, err := pgx.ParseDSN(dsn)
	if err != nil {
		logger.Log.Error("Parse config failed", zap.Error(err))
		return nil, err
	}
	conn, err := pgx.Connect(conConfig)
	if err != nil {
		logger.Log.Error("Ошибка подключения к DB", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	return conn, nil
}
