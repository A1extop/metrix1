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

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file}, nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func (p *Producer) WriteEvent(metricSt storage.MetricStorage) error {

	err := metricSt.ServerSendAllMetrics(p.file)
	return err
}

func WritingToDisk(times int, fileStoragePath string, memstorage *storage.MemStorage) {
	ticker := time.NewTicker(time.Duration(times) * time.Second)
	defer ticker.Stop()
	produc, err := NewProducer(fileStoragePath)
	if err != nil {
		log.Printf("Failed to open file for writing: %v", err)
		return
	}
	defer produc.Close()

	for {
		select {
		case <-ticker.C:
			err := produc.WriteEvent(memstorage)
			if err != nil {
				log.Printf("error writing to file: %v", err)
			}
		}
	}
}

type Consumer struct {
	file *os.File
}

func NewConsumer(filename string) (*Consumer, error) {

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{file: file}, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
