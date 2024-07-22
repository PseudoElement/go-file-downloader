package errors_module

type ErrorWithStatus interface {
	Error() string
	Status() int
}
