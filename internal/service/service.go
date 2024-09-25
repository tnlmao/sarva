package service

import (
	"sarva/internal/domain"
)

type FileService struct {
	processor  domain.FileProcessor
	repository domain.Repository
	consensus  domain.Consensus
	logger     domain.Logger
}

func NewFileService(p domain.FileProcessor, r domain.Repository, c domain.Consensus, l domain.Logger) *FileService {
	return &FileService{
		processor:  p,
		repository: r,
		consensus:  c,
		logger:     l,
	}
}

func (s *FileService) UploadFile(filePath string) error {
	file, err := s.processor.Process(filePath)
	if err != nil {
		s.logger.Log("ERROR", "Failed to process file")
		return err
	}

	if err := s.consensus.UpdateDatabase(*file); err != nil {
		s.logger.Log("ERROR", "Consensus failed")
		return err
	}

	// if err := s.repository.SaveFile(*file); err != nil {
	// 	s.logger.Log("ERROR", "Failed to save file to database")
	// 	return err
	// }

	s.logger.Log("INFO", "File uploaded successfully")
	return nil
}
