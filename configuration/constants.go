package configuration

import (
	"time"
)

// persistence type supported by
// Dahu.
type PersistenceType int

const (
	InMemory PersistenceType = 1 + iota
)

// configuration of Dahu
// data persistence
type Persistence struct {
	Type    PersistenceType
	Name    string
	TimeOut time.Duration
}

// configuration of Dahu
// http API
type Api struct {
	Port            int
	ShutdownTimeOut time.Duration
	Secret          string
}

// global configuration of
// Dahu
type Conf struct {
	PersistenceConf Persistence
	ApiConf         Api
	Close           chan interface{}
}

func InitConf() (c *Conf) {
	c = new(Conf)
	c.Close = make(chan interface{})
	c.PersistenceConf.Type = InMemory
	c.PersistenceConf.TimeOut = time.Second * 5
	c.PersistenceConf.Name = "dahu"
	c.ApiConf.Port = 80
	c.ApiConf.ShutdownTimeOut = 30 * time.Second
	return
}
