package job

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/jeromedoucet/dahu/core/model"
)

var wsEventListeners map[string][]*websocket.Conn

// listener is use internally
// to be pass though a channel that
// the goroutine in charge of notifications
// will consume.
type listener struct {
	jobId string
	conn  *websocket.Conn
}

type event struct {
	jobId string
	e     model.Event
}

var newListeners chan listener
var newEvents chan event

func init() {
	wsEventListeners = make(map[string][]*websocket.Conn)

	newListeners = make(chan listener)
	newEvents = make(chan event, 100)
	go startNotifier()
}

func startNotifier() {
	for {
		select {
		case newListner := <-newListeners:
			addEventListener(newListner)
		case newEvent := <-newEvents:
			broadcastEvent(newEvent)
		}
	}
}

func addEventListener(newListener listener) {
	_, exist := wsEventListeners[newListener.jobId]
	if !exist {
		wsEventListeners[newListener.jobId] = []*websocket.Conn{newListener.conn}
	} else {
		wsEventListeners[newListener.jobId] = append(wsEventListeners[newListener.jobId], newListener.conn)
	}
}

func broadcastEvent(newEvent event) {
	listeners, exist := wsEventListeners[newEvent.jobId]
	if exist {
		for _, l := range listeners {
			if err := l.WriteJSON(newEvent.e); err != nil {
				log.Println(err) // TODO close and remove conn when error !
			}
		}
	}
}

func AddWsEventListener(jobId string, conn *websocket.Conn) {
	newListeners <- listener{jobId, conn}
}

func Broadcast(jobId string, e model.Event) {
	newEvents <- event{jobId, e}
}
