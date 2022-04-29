package utils

type Err struct {
	message string
}

// custom error
func (err *Err) Error() string {
	return err.message
}
