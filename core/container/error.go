package container

type ContainerErrorType int

const (
	BadCredentials ContainerErrorType = 1 + iota
	RegistryNotFound
	OtherError
)

type ContainerError interface {
	Error() string
	ErrorType() ContainerErrorType
}

type simpleContainerError struct {
	msg     string
	errType ContainerErrorType
}

func (err simpleContainerError) Error() string {
	return err.msg
}

func (err simpleContainerError) ErrorType() ContainerErrorType {
	return err.errType
}

func newContainerError(msg string, errType ContainerErrorType) ContainerError {
	return simpleContainerError{msg: msg, errType: errType}
}
