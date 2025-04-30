package utils

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

func (e *CodeError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
