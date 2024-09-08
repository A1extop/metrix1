package postgrestore

import (
	"database/sql"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/A1extop/metrix1/internal/server/storage"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	Storage Storage
}

func NewRepository(s Storage) *Repository {
	return &Repository{Storage: s}
}

type Store struct {
	db *sql.DB
}

type Storage interface {
	Create() error
	Update() error
	CheckDBConnection(c *gin.Context)
	Writing(MetricStorage storage.MetricStorage)
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s');", tableName)
	err := db.QueryRow(query).Scan(&exists)
	return exists, err
}
func CreateOrConnectTable(db *sql.DB) {

	exists, err := tableExists(db, "MetricsGauges")
	if err != nil {
		log.Printf("error in checking for database presence: %v", err)
		return
	}
	if !exists {
		_, err = db.Exec("CREATE TABLE MetricsGauges (Name VARCHAR(255), VALUE DOUBLE PRECISION)")
		if err != nil {
			log.Printf("database creation error: %v", err)
		}
	}

	exists, err = tableExists(db, "MetricsCounters")
	if err != nil {
		log.Printf("error in checking for database presence: %v", err)
		return
	}
	if !exists {
		_, err = db.Exec("CREATE TABLE MetricsCounters (Name VARCHAR(255), VALUE INTEGER)")
		log.Printf("database creation error: %v", err)
	}
}

func ConnectDB(connectionToBD string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionToBD)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *Store) Create() error {
	return nil
}

func (s *Store) Update() error {
	return nil
}
func (s *Store) CheckDBConnection(c *gin.Context) {
	if s.db == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err := s.db.Ping()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
func (s *Store) Writing(MetricStorage storage.MetricStorage) {
	err := MetricStorage.RecordingMetricsDB(s.db)
	if err != nil {
		log.Printf("error writing to database: %v", err)
	}
}

func (r *Repository) Ping(c *gin.Context) {
	r.Storage.CheckDBConnection(c)
}

func WritingToBD(repos *Repository, times int, DBStoragePath string, memStorage *storage.MemStorage) {
	if DBStoragePath == "" {
		log.Println("DBStoragePath empty")
	}

	ticker := time.NewTicker(time.Duration(times) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Writing metrics to database...")
		repos.Storage.Writing(memStorage)
	}
}
