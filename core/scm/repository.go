package scm

var GitInstance ScmRepository = new(gitRepository)

type ScmRepository interface {
	CheckConnectionWithoutAuth(url string) ScmError
	CheckConnectionWithIdAndPassword(url string, id string, password string) ScmError
}
