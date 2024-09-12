package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/A1extop/metrix1/internal/server/domain"
	"github.com/gin-gonic/gin"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)

	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)

	ServerSendMetric(metricName string, metricType string) (interface{}, error)
	ServerSendAllMetricsHTML(c *gin.Context)
	MetricRecorder
}

type MetricRecorder interface {
	ServerSendAllMetricsToFile(*os.File) error
	ReadingMetricsFile(*os.File) error
	RecordingMetricsDB(db *sql.DB) error
}
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func isTypeExists(db *sql.DB, typeName string, tablesName string) (bool, error) {
	if tablesName != "MetricsGauges" && tablesName != "MetricsCounters" {
		return false, fmt.Errorf("invalid table name: %s", tablesName)
	}
	var exists bool

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE Name = $1)", tablesName)

	err := db.QueryRow(query, typeName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking type existence: %w", err)
	}
	return exists, nil
}
func Record[T float64 | int64](db *sql.DB, nameType string, tpName string, value T) error {
	exists, err := isTypeExists(db, nameType, tpName)
	if err != nil {
		return err
	}
	if exists {
		command := fmt.Sprintf("UPDATE %s SET Value = $1 WHERE Name = $2", tpName)
		_, err := db.Exec(command, value, nameType)
		if err != nil {
			return err
		}
	} else {
		command := fmt.Sprintf("INSERT INTO %s (Name, Value) VALUES($1, $2)", tpName)
		_, err = db.Exec(command, nameType, value)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *MemStorage) RecordingMetricsDB(db *sql.DB) error { // надо будет доделать
	tx, err := db.Begin()
	if err != nil {
		return nil
	}

	for nameType, value := range m.gauges {
		err := Record(db, nameType, "MetricsGauges", value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	for nameType, value := range m.counters {
		err := Record(db, nameType, "MetricsCounters", value)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}
	return nil
}

func (m *MemStorage) ReadingMetricsFile(file *os.File) error {
	var loadedGauges map[string]float64
	var loadedCounters map[string]int64

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return fmt.Errorf("file is empty, no metrics to restore:%v", err)
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedGauges); err != nil {
		return fmt.Errorf("error deserializing gauges: %v", err)
	}
	if err := decoder.Decode(&loadedCounters); err != nil {
		return fmt.Errorf("error deserializing counters: %v", err)
	}
	m.gauges = loadedGauges
	m.counters = loadedCounters
	log.Println("Metrics successfully restored from file")
	return nil
}
func (m *MemStorage) ServerSendMetric(metricName string, metricType string) (interface{}, error) {
	switch domain.MetricType(metricType) {
	case domain.Gauge:
		if value, ok := m.gauges[metricName]; ok {
			return value, nil
		}
	case domain.Counter:
		if value, ok := m.counters[metricName]; ok {
			return value, nil
		}
	}
	return "", fmt.Errorf("metric not found")
}

func (m *MemStorage) ServerSendAllMetricsHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	metricsTemplate, err := template.ParseFiles("templates/metrics.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: %v", err)
		return
	}

	err = metricsTemplate.Execute(c.Writer, m)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		return
	}
}
func (m *MemStorage) ServerSendAllMetricsToFile(file *os.File) error {
	dataGauges, err := json.MarshalIndent(m.gauges, "", " ")
	if err != nil {
		return fmt.Errorf("error serializing data: %v", err)
	}
	dataCounters, err := json.MarshalIndent(m.counters, "", " ")
	if err != nil {
		return fmt.Errorf("error serializing data: %v", err)
	}
	if _, err := file.Write(dataGauges); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	if _, err := file.Write(dataCounters); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	log.Println("Metrics successfully written to file ")
	return nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.counters[name] += value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := m.gauges[name]
	return value, exists
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := m.counters[name]
	return value, exists
}
