package tool

type HumanReadableError interface {
	error
	Human() string
	Cause() error
}

// Human-readable Error
type HRError struct {
	human string
	error error
}

func NewHRError(human string, err error) HumanReadableError {
	return &HRError{human: human, error: err}
}

// Just to complain error interface, it should be named String() I guess
func (e *HRError) Error() string {
	return e.error.Error()
}

func (e *HRError) Human() string {
	return e.human
}

func (e *HRError) Cause() error {
	return e.error
}
