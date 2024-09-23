package domain

type Logger interface {
	Log(level string, message string)
}
