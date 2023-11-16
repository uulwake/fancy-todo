package libs

type CustomError struct {
	HTTPCode int 
	BusinessCode int
	Message string
}

func (ce CustomError) Error() string {
	return ce.Message
}