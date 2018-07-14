package scm

type ScmErrorType int

const (
	BadCredentials ScmErrorType = 1 + iota
	RepositoryNotFound
	SshKeyReadingError
	OtherError
)

type ScmError interface {
	Error() string
	ErrorType() ScmErrorType
}

type simpleScmError struct {
	msg     string
	errType ScmErrorType
}

func (err simpleScmError) Error() string {
	return err.msg
}

func (err simpleScmError) ErrorType() ScmErrorType {
	return err.errType
}

func newScmError(msg string, errType ScmErrorType) ScmError {
	return simpleScmError{msg: msg, errType: errType}
}
