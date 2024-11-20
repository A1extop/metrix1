## Пакет github.com/A1extop/metrix1/cmd/agent

## Пакет github.com/A1extop/metrix1/cmd/server

## Пакет github.com/A1extop/metrix1/config/agentconfig
package agentconfig // import "github.com/A1extop/metrix1/config/agentconfig"


TYPES

type Parameters struct {
	AddressHTTP    string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func NewParameters() *Parameters

func (p *Parameters) GetParameters()

func (p *Parameters) GetParametersEnvironmentVariables()

## Пакет github.com/A1extop/metrix1/config/serverconfig
package serverconfig // import "github.com/A1extop/metrix1/config/serverconfig"


TYPES

type Parameters struct {
	AddressHTTP     string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	AddrDB          string
	Key             string
}

func NewParameters() *Parameters

func (p *Parameters) GetParameters()

func (p *Parameters) GetParametersEnvironmentVariables()

## Пакет github.com/A1extop/metrix1/internal/agent/agentsend
package agentsend // import "github.com/A1extop/metrix1/internal/agent/agentsend"


FUNCTIONS

func SendMetric(client *http.Client, serverAddress string, metric js.Metrics, key string) error
func SendMetrics(client *http.Client, serverAddress string, metrics []js.Metrics, key string) error
## Пакет github.com/A1extop/metrix1/internal/agent/compress
package compress // import "github.com/A1extop/metrix1/internal/agent/compress"


FUNCTIONS

func CompressData(data []byte) ([]byte, error)
    CompressData accepts an array of bytes and returns its compressed version
    using gzip.

## Пакет github.com/A1extop/metrix1/internal/agent/hash
package hash // import "github.com/A1extop/metrix1/internal/agent/hash"


FUNCTIONS

func SignRequestWithSHA256(metrics []byte, key string) (string, error)
## Пакет github.com/A1extop/metrix1/internal/agent/json
package json // import "github.com/A1extop/metrix1/internal/agent/json"


TYPES

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func GetParametersJSON(c *gin.Context) (*Metrics, error)

func NewMetrics() *Metrics

## Пакет github.com/A1extop/metrix1/internal/agent/storage
package storage // import "github.com/A1extop/metrix1/internal/agent/storage"


TYPES

type MemStorage struct {
	// Has unexported fields.
}

func NewMemStorage() *MemStorage

func (m *MemStorage) Report(ctx context.Context, client *http.Client, serverAddress string, key string, rateLimit int, reportTicker *time.Ticker)

func (m *MemStorage) ReportMetrics(client *http.Client, serverAddress string, key string)

func (m *MemStorage) UpdateMetrics()

type MetricUpdater interface {
	UpdateMetrics()
	ReportMetrics(client *http.Client, serverAddress string, key string)

	Report(ctx context.Context, client *http.Client, serverAddress string, key string, rateLimit int, reportTicker *time.Ticker)
	// Has unexported methods.
}

## Пакет github.com/A1extop/metrix1/internal/agent/updatereportmetrics
package updatereportmetrics // import "github.com/A1extop/metrix1/internal/agent/updatereportmetrics"


TYPES

type Updater struct {
	// Has unexported fields.
}

func NewAction(updater storage.MetricUpdater) *Updater

func (u *Updater) Action(ctx context.Context, parameters *config.Parameters)

## Пакет github.com/A1extop/metrix1/internal/server/compress
package compress // import "github.com/A1extop/metrix1/internal/server/compress"

отвечает за чтение сжатых данных

отвечает за сжатие данных

FUNCTIONS

func CompressData() gin.HandlerFunc
func DeCompressData() gin.HandlerFunc
## Пакет github.com/A1extop/metrix1/internal/server/data
package data // import "github.com/A1extop/metrix1/internal/server/data"


FUNCTIONS

func ReadingFromDisk(fileStoragePath string, memStorage *storage.MemStorage)
func WritingToDisk(times int, fileStoragePath string, memStorage *storage.MemStorage)

TYPES

type Producer struct {
	// Has unexported fields.
}

func NewProducer(filePath string) (*Producer, error)

func (p *Producer) Close() error

func (p *Producer) WriteEvent(metricSt storage.MetricStorage) error

## Пакет github.com/A1extop/metrix1/internal/server/domain
package domain // import "github.com/A1extop/metrix1/internal/server/domain"


VARIABLES

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

FUNCTIONS

func Validate(metricsJs *js.Metrics, c *gin.Context) error

TYPES

type Metric struct {
	Name  string
	Type  MetricType
	Value interface{}
}

func NewMetric(name string, metricType MetricType, value interface{}) (*Metric, error)

func (m *Metric) ValidateValue() error

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
## Пакет github.com/A1extop/metrix1/internal/server/hash
package hash // import "github.com/A1extop/metrix1/internal/server/hash"


FUNCTIONS

func WorkingWithHash(key string) gin.HandlerFunc
    WorkingWithHash performs work with hash.

## Пакет github.com/A1extop/metrix1/internal/server/http
package http // import "github.com/A1extop/metrix1/internal/server/http"


FUNCTIONS

func GetValue(metricsJs *js.Metrics) string
func NewRouter(handler *Handler, repos *psql.Repository, key string) *gin.Engine

TYPES

type Handler struct {
	// Has unexported fields.
}

func NewHandler(storage storage.MetricStorage) *Handler

func (h *Handler) DerivationMetric(c *gin.Context)
    DerivationMetric processes the request to get a specific metric. It takes
    the metric type and the metric name from the query parameters, and then
    extracts the metric from storage and returns it in JSON format.

func (h *Handler) DerivationMetrics(c *gin.Context)
    DerivationMetrics output metrics outputs all metrics in HTML format.

func (h *Handler) GetJSON(c *gin.Context)
    GetJSON processes a request to get metrics in JSON format.

func (h *Handler) Update(c *gin.Context)
    Update processes an HTTP request to update the metric according to the
    specified parameters.

func (h *Handler) UpdateJSON(c *gin.Context)
    UpdateJSON handles updating metrics in JSON format.

func (h *Handler) UpdatePacketMetricsJSON(c *gin.Context)
    UpdatePacketMetricsJSON processes an HTTP request to update metrics accepted
    in JSON format. The function gets an array of metrics from the request body,
    validates each metric and updates their values in the repository.

## Пакет github.com/A1extop/metrix1/internal/server/json
package json // import "github.com/A1extop/metrix1/internal/server/json"


TYPES

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func GetParametersJSON(c *gin.Context) (*Metrics, error)

func GetParametersMassiveJSON(c *gin.Context) ([]Metrics, error)

func NewMetrics() *Metrics

## Пакет github.com/A1extop/metrix1/internal/server/logging
package logging // import "github.com/A1extop/metrix1/internal/server/logging"


FUNCTIONS

func LoggingGet(logger *zap.SugaredLogger) gin.HandlerFunc
func LoggingPost(logger *zap.SugaredLogger) gin.HandlerFunc
func New() *zap.SugaredLogger

TYPES

type CustomResponseWriter struct {
	gin.ResponseWriter
	// Has unexported fields.
}

func (w *CustomResponseWriter) Write(b []byte) (int, error)

## Пакет github.com/A1extop/metrix1/internal/server/storage
package storage // import "github.com/A1extop/metrix1/internal/server/storage"


FUNCTIONS

func Record[T float64 | int64](db *sql.DB, nameType string, tpName string, value T) error
    Record writes the metric value to the database.


TYPES

type MemStorage struct {
	// Has unexported fields.
}

func NewMemStorage() *MemStorage

func (m *MemStorage) GetCounter(name string) (int64, bool)
    GetCounter getting a metric by name.

func (m *MemStorage) GetGauge(name string) (float64, bool)
    GetGauge getting a metric by name.

func (m *MemStorage) ReadingMetricsFile(file *os.File) error
    ReadingMetricsFile reads metrics from a file.

func (m *MemStorage) RecordingMetricsDB(db *sql.DB) (err error)
    RecordingMetricsDB writes metrics to the database as part of a transaction.

func (m *MemStorage) ServerFindMetric(metricName string, metricType string) (interface{}, error)
    ServerFindMetric searches for a metric by name and type in MemStorage.

func (m *MemStorage) ServerSendAllMetricsHTML(c *gin.Context)

func (m *MemStorage) ServerSendAllMetricsToFile(file *os.File) error
    ServerSendAllMetricsToFile serializes all metrics from MemStorage and writes
    them to the specified file.

func (m *MemStorage) UpdateCounter(name string, value int64)
    UpdateCounter update metric by name.

func (m *MemStorage) UpdateGauge(name string, value float64)
    UpdateGauge update metric by name.

type MetricRecorder interface {
	ServerSendAllMetricsToFile(*os.File) error
	ReadingMetricsFile(*os.File) error
	RecordingMetricsDB(db *sql.DB) error
}

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)

	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)

	ServerFindMetric(metricName string, metricType string) (interface{}, error)
	ServerSendAllMetricsHTML(c *gin.Context)
	MetricRecorder
}

## Пакет github.com/A1extop/metrix1/internal/server/store/postgrestore
package postgrestore // import "github.com/A1extop/metrix1/internal/server/store/postgrestore"


FUNCTIONS

func ConnectDB(connectionToBD string) (*sql.DB, error)
func CreateOrConnectTable(db *sql.DB)
func WritingToBD(repos *Repository, times int, DBStoragePath string, memStorage *storage.MemStorage)

TYPES

type Repository struct {
	Storage Storage
}

func NewRepository(s Storage) *Repository

func (r *Repository) Ping(c *gin.Context)

type Storage interface {
	Create() error
	Update() error
	CheckDBConnection(c *gin.Context)
	Writing(MetricStorage storage.MetricStorage)
}

type Store struct {
	// Has unexported fields.
}

func NewStore(db *sql.DB) *Store

func (s *Store) CheckDBConnection(c *gin.Context)

func (s *Store) Create() error

func (s *Store) Update() error

func (s *Store) Writing(MetricStorage storage.MetricStorage)

## Пакет github.com/A1extop/metrix1/internal/server/usecase
package usecase // import "github.com/A1extop/metrix1/internal/server/usecase"


FUNCTIONS

func UpdateMetric(storage storage.MetricStorage, metricType, metricValue, metricName string) error
    UpdateMetric updates the metric in the repository based on the data
    provided.

## Пакет github.com/A1extop/metrix1/pkg/validator
package validator // import "github.com/A1extop/metrix1/pkg/validator"


FUNCTIONS

func ValidateRequest(c *gin.Context, expectedContentType, metricName string) bool
