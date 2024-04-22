// Package dbstorage - db storage implements the repository interface
package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/postgresql"
)

type notFoundMetricError struct {
	tp   string
	name string
}

func (e notFoundMetricError) Error() string {
	return fmt.Sprintf("not found metric type:\"%v\" with name:\"%v\"", e.tp, e.name)
}

// DBStorage - storage
type DBStorage struct {
	dbConnections *sql.DB
}

// NewStorage - method creates a new DBStorage
func NewStorage(ctx context.Context, dsn string) (*DBStorage, error) {
	dbCon, err := postgresql.NewConn(dsn)
	if err != nil {
		logger.Log.Error("Get db connection failed", zap.Error(err))
		return nil, err
	}
	s := &DBStorage{dbConnections: dbCon}
	const createTableSQL = `create table IF NOT EXISTS metrics (ID text PRIMARY KEY,MType text, Delta bigint, Value double precision)`
	_, err = s.dbConnections.ExecContext(ctx, createTableSQL)
	if err != nil {
		logger.Log.Error("Create table failed", zap.Error(err))
		return nil, err
	}
	const createIndexSQL = `CREATE index IF NOT EXISTS ix_id_mtype ON metrics (id,MType)`
	_, err = s.dbConnections.ExecContext(ctx, createIndexSQL)
	if err != nil {
		logger.Log.Error("Create table failed", zap.Error(err))
		return nil, err
	}
	return s, nil
}

// SetGauge - method writes values to storage
func (ms *DBStorage) SetGauge(ctx context.Context, key string, val float64) error {
	const upsertSQL = `INSERT INTO metrics (id, MType, Value) VALUES ($1, 'gauge', $2) ON CONFLICT(id) DO UPDATE set Value = EXCLUDED.Value`
	_, err := ms.dbConnections.ExecContext(ctx, upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return fmt.Errorf("upsert: %w", err)
	}
	return nil
}

// SetCounter - method writes values to storage
func (ms *DBStorage) SetCounter(ctx context.Context, key string, val int64) error {
	const upsertSQL = `INSERT INTO metrics AS a (id, MType, Delta) VALUES ($1, 'counter', $2) ON CONFLICT(id) DO UPDATE set Delta = a.Delta + EXCLUDED.Delta`
	_, err := ms.dbConnections.ExecContext(ctx, upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return fmt.Errorf("upsert: %w", err)
	}
	return nil
}

// GetAll - method returns all values from storage
func (ms *DBStorage) GetAll(ctx context.Context) (map[string]map[string]string, error) {
	ret := make(map[string]map[string]string, 2)
	retGauge := make(map[string]string)
	retCounter := make(map[string]string)

	const selectSQL = `select id, MType, coalesce(Delta,-1), coalesce(Value,-1) from metrics`
	rows, err := ms.dbConnections.QueryContext(ctx, selectSQL)
	if err != nil {
		logger.Log.Error("Select all failed", zap.Error(err))
		return ret, fmt.Errorf("select: %w", err)
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error("Select all failed", zap.Error(err))
		return ret, fmt.Errorf("select rows: %w", err)
	}
	defer rows.Close()
	var (
		id, mType string
		delta     int64
		value     float64
	)
	for rows.Next() {
		err = rows.Scan(&id, &mType, &delta, &value)
		if err != nil {
			logger.Log.Error("Select rows failed", zap.Error(err))
			return ret, nil
		}
		if mType == "counter" {
			retCounter[id] = fmt.Sprintf("%d", delta)
		}
		if mType == "gauge" {
			retGauge[id] = fmt.Sprintf("%g", value)
		}
	}
	ret["counter"] = retCounter
	ret["gauge"] = retGauge
	return ret, nil
}

// GetSingle - method returns single value from storage
func (ms *DBStorage) GetSingle(ctx context.Context, tp string, name string) (string, error) {

	const selectSQL = `select id, MType, coalesce(Delta,-1), coalesce(Value,-1) from metrics where MType = $1 and id = $2`
	row := ms.dbConnections.QueryRowContext(ctx, selectSQL, tp, name)
	var (
		id, mType string
		delta     int64
		value     float64
	)
	err := row.Scan(&id, &mType, &delta, &value)
	if err != nil {
		logger.Log.Error("Select single failed", zap.String("tp", tp), zap.String("name", name), zap.Error(err))
		return "", fmt.Errorf("select single: %w", err)
	}
	switch mType {
	case "gauge":
		return fmt.Sprintf("%g", value), nil
	case "counter":
		return fmt.Sprintf("%d", delta), nil
	}
	return "", notFoundMetricError{tp, name}
}

// LoadFromFile - mock for LoadFromFile method (needed only for file storage)
func (ms *DBStorage) LoadFromFile(ctx context.Context) error {
	return nil
}

// SaveToFile - mock for SaveToFile method (needed only for file storage)
func (ms *DBStorage) SaveToFile(ctx context.Context) error {
	return nil
}

// Ping - method tests the connection to the database
func (ms *DBStorage) Ping(ctx context.Context) error {
	return ms.dbConnections.PingContext(ctx)
}

// SetBatch - method writes array values to storage
func (ms *DBStorage) SetBatch(ctx context.Context, metrics []metrics.Metrics) error {
	tx, err := ms.dbConnections.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("open transaction: %w", err)
	}
	defer tx.Rollback()

	const upsertSQL = `INSERT INTO metrics AS a (id, MType, Value, Delta) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE set Delta = a.Delta + EXCLUDED.Delta,Value = EXCLUDED.Value`
	stmt, err := tx.PrepareContext(ctx, upsertSQL)
	if err != nil {
		return fmt.Errorf("prepare sql: %w", err)
	}
	defer stmt.Close()

	for _, v := range metrics {
		_, err := stmt.ExecContext(ctx, v.ID, v.MType, v.Value, v.Delta)
		if err != nil {
			return fmt.Errorf("exec sql: %w", err)
		}
	}
	return tx.Commit()
}

// IsRetryable - determines the type of error (whether it is suitable for re-execution)
func (ms *DBStorage) IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
