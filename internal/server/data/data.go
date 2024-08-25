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
	err := metricSt.ServerSendAllMetrics(p.file)
	return err
}

func WritingToDisk(times int, fileStoragePath string, memstorage *storage.MemStorage) {
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

func ReadingFromDisk(fileStoragePath string, memstorage *storage.MemStorage) {
	if fileStoragePath == "" {
		return
	}
	file, err := os.Open(fileStoragePath)
	err = memstorage.RecordingMetricsFile(file)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
