package postgrestore

import (
	"database/sql"
	"net/http"

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
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
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
func (r *Repository) Ping(c *gin.Context) {
	r.Storage.CheckDBConnection(c)
}
