package model

// Error that mus be used when
// trying to update some entity
// with outdated data
type outDated struct {
	msg string
}

func (o outDated) Error() string {
	return o.msg
}

// return an error when trying to
// update some entity with outdated
// data
func NewOutDated(msg string) error {
	e := new(outDated)
	e.msg = msg
	return e
}

func IsOutDated(err error) bool {
	// accept pointer and value both
	_, res := err.(outDated)
	if !res {
		_, res = err.(*outDated)
	}
	return res
}

// Error that mus be used when
// trying to update some entity
// That doesn't exist anymore.
// this is different from a non existing
// entity in the way that it applies on
// data that should be remove after a while
// like model.JobRun
type noMorePersisted struct {
	msg string
}

func (o noMorePersisted) Error() string {
	return o.msg
}

func NewNoMorePersisted(msg string) error {
	e := new(noMorePersisted)
	e.msg = msg
	return e
}

func IsNoMorePersisted(err error) bool {
	// accept pointer and value both
	_, res := err.(noMorePersisted)
	if !res {
		_, res = err.(*noMorePersisted)
	}
	return res
}
