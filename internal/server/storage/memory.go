package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/A1extop/metrix1/internal/server/domain"
	"github.com/gin-gonic/gin"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)

	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)

	ServerFindMetric(metricName string, metricType string) (interface{}, error)
	ServerSendAllMetricsHTML(c *gin.Context)
	MetricRecorder
}

type MetricRecorder interface {
	ServerSendAllMetricsToFile(*os.File) error
	ReadingMetricsFile(*os.File) error
	RecordingMetricsDB(db *sql.DB) error
}
type MemStorage struct {
	mv       sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
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

// Record writes the metric value to the database.
func Record[T float64 | int64](db *sql.DB, nameType string, tpName string, value T) error {
	exists, err := isTypeExists(db, nameType, tpName)
	if err != nil {
		return err
	}
	if exists {
		command := fmt.Sprintf("UPDATE %s SET Value = $1 WHERE Name = $2", tpName)
		_, err = db.Exec(command, value, nameType)
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

// RecordingMetricsDB writes metrics to the database as part of a transaction.
func (m *MemStorage) RecordingMetricsDB(db *sql.DB) error { //возможно здесь будет ошибка
	var wg sync.WaitGroup
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	TimesDuration := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	targetError := errors.New("driver: bad connection")

	recordGauge := func(nameType string, value float64) {
		for _, times := range TimesDuration {
			err = Record(db, nameType, "MetricsGauges", value)
			if err == nil {
				break
			}
			if errors.Is(err, targetError) {
				log.Printf("recover: %v", targetError)
				time.Sleep(times)
			} else {
				return
			}
		}
	}

	recordCounter := func(nameType string, value int64) {
		for _, times := range TimesDuration {
			err = Record(db, nameType, "MetricsCounters", value)
			if err == nil {
				break
			}
			if errors.Is(err, targetError) {
				log.Printf("recover: %v", targetError)
				time.Sleep(times)
			} else {
				return
			}
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		m.mv.RLock()
		defer m.mv.RUnlock()
		for nameType, value := range m.Gauges {
			recordGauge(nameType, value)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		m.mv.RLock()
		defer m.mv.RUnlock()
		for nameType, value := range m.Counters {
			recordCounter(nameType, value)
			if err != nil {
				return
			}
		}
	}()

	wg.Wait()
	return err
}

// ReadingMetricsFile reads metrics from a file.
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
	m.mv.Lock()
	m.Gauges = loadedGauges
	m.Counters = loadedCounters
	m.mv.Unlock()
	log.Println("Metrics successfully restored from file")
	return nil
}

// ServerFindMetric searches for a metric by name and type in MemStorage.
func (m *MemStorage) ServerFindMetric(metricName string, metricType string) (interface{}, error) {
	m.mv.RLock()
	defer m.mv.RUnlock()
	switch domain.MetricType(metricType) {
	case domain.Gauge:
		if value, ok := m.Gauges[metricName]; ok {
			return value, nil
		}
	case domain.Counter:
		if value, ok := m.Counters[metricName]; ok {
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

// ServerSendAllMetricsToFile serializes all metrics from MemStorage and writes them to the specified file.
func (m *MemStorage) ServerSendAllMetricsToFile(file *os.File) error {
	dataGauges, err := json.MarshalIndent(m.Gauges, "", " ")
	if err != nil {
		return fmt.Errorf("error serializing data: %v", err)
	}
	dataCounters, err := json.MarshalIndent(m.Counters, "", " ")
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
		mv:       sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

// UpdateGauge update metric by name.
func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.mv.Lock()
	defer m.mv.Unlock()
	m.Gauges[name] = value
}

// UpdateCounter update metric by name.
func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.mv.Lock()
	defer m.mv.Unlock()
	m.Counters[name] += value
}

// GetGauge getting a metric by name.
func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.mv.RLock()
	defer m.mv.RUnlock()
	value, exists := m.Gauges[name]
	return value, exists
}

// GetCounter getting a metric by name.
func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.mv.RLock()
	defer m.mv.RUnlock()
	value, exists := m.Counters[name]
	return value, exists
}
