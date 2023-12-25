package actions

type StorInterface interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll()
}
