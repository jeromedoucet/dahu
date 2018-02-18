package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// todo parse arguments or conf file ?
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444       // todo look if it is really necessary
	conf.ApiConf.Secret = "secret" // todo generate it
	apiInstance := api.InitRoute(conf)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.ApiConf.Port),
		Handler: apiInstance.Handler(),
	}

	go func() {
		log.Printf("INFO >> Listening on http://0.0.0.0:%d\n", conf.ApiConf.Port)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for kill signal
	<-stop

	// shutdown gracefully the entire system :
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	log.Println("INFO >> begin the shutdown of the system")
	apiInstance.Close() // todo pass the context

	s.Shutdown(ctx)

	log.Println("INFO >> server gracefully stopped")

}
