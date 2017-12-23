package model

type ScmType int

const (
	GIT ScmType = 1 + iota
	SVN
)

type Scm struct {
	RepoUrl string
	Type    ScmType
}

func (s Scm) getImage() string {
	switch s.Type {
	case GIT:
		return "dahuci/git"
	case SVN:
		return "dahuci/svn"
	default:
		return ""
	}
}
