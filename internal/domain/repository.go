package domain

type Repository interface {
	SaveFile(file File) error
}
