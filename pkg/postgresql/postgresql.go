package postgresql

import (
	"database/sql"

	"github.com/LobovVit/metric-collector/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func NweConn(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("DB open failed", zap.Error(err))
		return nil, err
	}
	return conn, nil
}
