package fileprocessor

import (
	"errors"
	"fmt"
	"os"
	"sarva/internal/domain"
)

type DummyFileProcessor struct{}

func NewFileProcessor() *DummyFileProcessor {
	return &DummyFileProcessor{}
}
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
	fmt.Println(fileInfo.Size())
	if fileInfo.Size() == 0 {
		return nil, errors.New("file is empty")
	}
	return &domain.File{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
	}, nil
}
