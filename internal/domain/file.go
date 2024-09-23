package domain

type File struct {
	Name string
	Size int64
}

type FileProcessor interface {
	Process(filePath string) (*File, error)
}
