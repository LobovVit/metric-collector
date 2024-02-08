package dbstorage

import (
	"database/sql"
	"fmt"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/postgresql"
	"go.uber.org/zap"
)

type notFoundMetricError struct {
	tp   string
	name string
}

func (e notFoundMetricError) Error() string {
	return fmt.Sprintf("not found metric type:\"%v\" with name:\"%v\"", e.tp, e.name)
}

type DBStorage struct {
	dbConnections *sql.DB
}

func NewStorage(dsn string) *DBStorage {
	dbCon, err := postgresql.NweConn(dsn)
	if err != nil {
		logger.Log.Error("Get db connection failed", zap.Error(err))
	}
	s := &DBStorage{dbConnections: dbCon}
	createSQL := `create table IF NOT EXISTS metrics (ID text PRIMARY KEY,MType text, Delta bigint, Value double precision)`
	_, err = s.dbConnections.Exec(createSQL)
	if err != nil {
		logger.Log.Error("Create table failed", zap.Error(err))
	}
	return s
}

func (ms *DBStorage) SetGauge(key string, val float64) error {
	upsertSQL := `INSERT INTO metrics (id, MType, Value) VALUES ($1, 'gauge', $2) ON CONFLICT(id) DO UPDATE set Value = EXCLUDED.Value`
	_, err := ms.dbConnections.Exec(upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return fmt.Errorf("upsert failed: %w", err)
	}
	return nil
}

func (ms *DBStorage) SetCounter(key string, val int64) error {
	upsertSQL := `INSERT INTO metrics AS a (id, MType, Delta) VALUES ($1, 'counter', $2) ON CONFLICT(id) DO UPDATE set Delta = a.Delta + EXCLUDED.Delta`
	_, err := ms.dbConnections.Exec(upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return fmt.Errorf("upsert failed: %w", err)
	}
	return nil
}

func (ms *DBStorage) GetAll() map[string]map[string]string {
	ret := make(map[string]map[string]string, 2)
	retGauge := make(map[string]string)
	retCounter := make(map[string]string)

	selectSQL := `select id, MType, coalesce(Delta,-1), coalesce(Value,-1) from metrics`
	rows, err := ms.dbConnections.Query(selectSQL)
	if err != nil {
		logger.Log.Error("Select all failed", zap.Error(err))
		return ret
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error("Select all failed", zap.Error(err))
		return ret
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
		}
		if mType == "counter" {
			retCounter[id] = fmt.Sprintf("%d", delta)
		}
		if mType == "gauge" {
			retGauge[id] = fmt.Sprintf("%g", value)
		}
	}
	logger.Log.Info("cnt", zap.Int("counter", len(retCounter)))
	logger.Log.Info("cnt", zap.Int("gauge", len(retGauge)))
	ret["counter"] = retCounter
	ret["gauge"] = retGauge
	return ret
}

func (ms *DBStorage) GetSingle(tp string, name string) (string, error) {

	selectSQL := `select id, MType, coalesce(Delta,-1), coalesce(Value,-1) from metrics where MType = $1 and id = $2`
	row := ms.dbConnections.QueryRow(selectSQL, tp, name)
	var (
		id, mType string
		delta     int64
		value     float64
	)
	err := row.Scan(&id, &mType, &delta, &value)
	if err != nil {
		logger.Log.Error("Select single failed", zap.Error(err))
		return "", fmt.Errorf("select single failed: %w", err)
	}
	switch mType {
	case "gauge":
		return fmt.Sprintf("%g", value), nil
	case "counter":
		return fmt.Sprintf("%d", delta), nil
	}
	return "", notFoundMetricError{tp, name}
}

func (ms *DBStorage) LoadFromFile() error {
	return nil
}

func (ms *DBStorage) SaveToFile() error {
	return nil
}

func (ms *DBStorage) Ping() error {
	return ms.dbConnections.Ping()
}
func (ms *DBStorage) SetBatch(metrics []metrics.Metrics) error {
	tx, err := ms.dbConnections.Begin()
	if err != nil {
		return fmt.Errorf("open transaction failed: %w", err)
	}
	defer tx.Rollback()

	upsertSQL := `INSERT INTO metrics AS a (id, MType, Value, Delta) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE set Delta = a.Delta + EXCLUDED.Delta,Value = EXCLUDED.Value`
	stmt, err := tx.Prepare(upsertSQL)
	if err != nil {
		return fmt.Errorf("prepare sql failed: %w", err)
	}
	defer stmt.Close()

	for _, v := range metrics {
		_, err := stmt.Exec(v.ID, v.MType, v.Value, v.Delta)
		if err != nil {
			return fmt.Errorf("exec sql failed: %w", err)
		}
	}
	return tx.Commit()
}
