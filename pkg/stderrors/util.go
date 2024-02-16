package stderrors

func Populate(we *Error, err error) *Error {
	we.InternalMessage = err.Error()
	return we
}
