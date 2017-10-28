package tests

import (
	"fmt"
	"os"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/persistence"
)

// will clean the data inside persistence
// must be done after a test has run.
func CleanPersistence(conf *configuration.Conf) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CleanPersistence", r)
		}
	}()
	close(conf.Close)
	rep := persistence.GetRepository(conf)
	rep.WaitClose()
	os.Remove(conf.PersistenceConf.Name)
}
