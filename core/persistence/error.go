package persistence

type PersistenceErrorType int

const (
	NotFound PersistenceErrorType = 1 + iota
	OtherError
)

type PersistenceError interface {
	Error() string
	ErrorType() PersistenceErrorType
}

type simplePersistenceError struct {
	msg     string
	errType PersistenceErrorType
}

func (err simplePersistenceError) Error() string {
	return err.msg
}

func (err simplePersistenceError) ErrorType() PersistenceErrorType {
	return err.errType
}

func newPersistenceError(msg string, errType PersistenceErrorType) PersistenceError {
	return simplePersistenceError{msg: msg, errType: errType}
}

func wrapError(err error) PersistenceError {
	if err == nil {
		return nil
	}
	persistenceErr, isPersistenceErr := err.(PersistenceError)
	if isPersistenceErr {
		return persistenceErr
	} else {
		return newPersistenceError(err.Error(), OtherError)
	}
}
