package data

import (
	"log"
	"os"
	"time"

	"github.com/A1extop/metrix1/internal/server/storage"
)

type Producer struct {
	file *os.File
}

func NewProducer(filePath string) (*Producer, error) {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file}, nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func (p *Producer) WriteEvent(metricSt storage.MetricStorage) error {
	p.file.Truncate(0)
	p.file.Seek(0, 0)
	err := metricSt.ServerSendAllMetricsToFile(p.file)
	return err
}

func WritingToDisk(times int, fileStoragePath string, memStorage *storage.MemStorage) {
	if fileStoragePath == "" {
		return
	}
	ticker := time.NewTicker(time.Duration(times) * time.Second)
	defer ticker.Stop()
	produc, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Printf("Failed to open file for writing: %v", err)
		return
	}
	defer produc.Close()

	for range ticker.C {
		err := produc.WriteEvent(memStorage)
		if err != nil {
			log.Printf("error writing to file: %v", err)
		}
	}
}

func ReadingFromDisk(fileStoragePath string, memStorage *storage.MemStorage) {
	if fileStoragePath == "" {
		log.Println("fileStoragePath empty")
		return
	}
	file, err := os.Open(fileStoragePath)
	if err != nil {
		log.Println(err)
		return
	}
	err = memStorage.ReadingMetricsFile(file)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
