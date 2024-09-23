package fileprocessor

import (
	"os"
	"sarva/internal/domain"
)

type DummyFileProcessor struct{}

func (p *DummyFileProcessor) Process(filePath string) (*domain.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &domain.File{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
	}, nil
}
