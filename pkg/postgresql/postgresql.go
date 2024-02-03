package postgresql

import (
	"context"
	"database/sql"

	"github.com/LobovVit/metric-collector/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func NweConn(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("DB open failed", zap.Error(err))
		return nil, err
	}
	return conn, nil
}
