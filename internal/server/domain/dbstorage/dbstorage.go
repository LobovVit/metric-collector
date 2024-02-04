package dbstorage

import (
	"database/sql"
	"fmt"

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

type DbStorage struct {
	dbConnections *sql.DB
}

func NewStorage(dsn string) *DbStorage {
	dbCon, err := postgresql.NweConn(dsn)
	if err != nil {
		logger.Log.Error("Get db connection failed", zap.Error(err))
	}
	s := &DbStorage{dbConnections: dbCon}
	createSQL := `create table IF NOT EXISTS metrics (ID text PRIMARY KEY,MType text, Delta bigint, Value double precision)`
	_, err = s.dbConnections.Exec(createSQL)
	if err != nil {
		logger.Log.Error("Create table failed", zap.Error(err))
	}
	return s
}

func (ms *DbStorage) SetGauge(key string, val float64) error {
	upsertSQL := `INSERT INTO metrics (id, MType, Value) VALUES ($1, 'Gauge', $2) ON CONFLICT(id) DO UPDATE set Value = EXCLUDED.Value`
	_, err := ms.dbConnections.Exec(upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return err
	}
	return nil
}

func (ms *DbStorage) SetCounter(key string, val int64) error {
	upsertSQL := `INSERT INTO metrics AS a (id, MType, Delta) VALUES ($1, 'Counter', $2) ON CONFLICT(id) DO UPDATE set Delta = a.Delta + EXCLUDED.Delta`
	_, err := ms.dbConnections.Exec(upsertSQL, key, val)
	if err != nil {
		logger.Log.Error("Upsert failed", zap.Error(err))
		return err
	}
	return nil
}

func (ms *DbStorage) GetAll() map[string]map[string]string {
	ret := make(map[string]map[string]string, 2)
	retGauge := make(map[string]string)
	retCounter := make(map[string]string)

	selectSQL := `select id, MType, coalesce(Delta,-1), coalesce(Value,-1) from metrics`
	rows, err := ms.dbConnections.Query(selectSQL)
	if err != nil {
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
		fmt.Println(id, "-", mType, "-", delta, "-", value)
		if mType == "Counter" {
			retCounter[id] = fmt.Sprintf("%d", delta)
		}
		if mType == "Gauge" {
			retGauge[id] = fmt.Sprintf("%g", value)
		}
	}
	ret["Counter"] = retCounter
	ret["Gauge"] = retGauge
	return ret
}

func (ms *DbStorage) GetSingle(tp string, name string) (string, error) {

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
	case "Gauge":
		return fmt.Sprintf("%g", value), nil
	case "Counter":
		return fmt.Sprintf("%d", delta), nil
	}
	return "", notFoundMetricError{tp, name}
}

func (ms *DbStorage) LoadFromFile() error {
	return nil
}

func (ms *DbStorage) SaveToFile() error {
	return nil
}

func (ms *DbStorage) Ping() error {
	return ms.dbConnections.Ping()
}
