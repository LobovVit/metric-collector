package postgresql

import (
	"context"
	"database/sql"

	"github.com/LobovVit/metric-collector/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*sql.DB, error) {
	logger.Log.Info("YA=DSN", zap.String("dsn", dsn))
	//dbctx, chancel := context.WithTimeout(ctx, time.Second*5)
	//defer chancel()
	conn, err := sql.Open("postgres", dsn)
	//conConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		logger.Log.Error("Parse config failed", zap.Error(err))
		return nil, err
	}
	//conn, err := pgx.ConnectConfig(dbctx, conConfig)
	//if err != nil {
	//	logger.Log.Error("Ошибка подключения к DB", zap.String("dsn", dsn), zap.Error(err))
	//	return nil, err
	//}
	return conn, nil
}
