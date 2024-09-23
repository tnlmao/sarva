package domain

type Consensus interface {
	UpdateDatabase(file File) error
}
