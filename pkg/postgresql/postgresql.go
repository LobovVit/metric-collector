package postgresql

import (
	"context"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Log.Error("Ошибка подключения к DB", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	return conn, nil
}
