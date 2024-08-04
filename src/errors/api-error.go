package app_errors

type ApiError struct {
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) Status() int {
	return 400
}
