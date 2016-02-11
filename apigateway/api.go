package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/yetis-br/tp-server/mq"
)

//Tasks define a message queue to communicate to the workers
var Tasks *mq.MessageQueue

//CallbackMessages defined to manage the api callback responses
var CallbackMessages <-chan amqp.Delivery

func main() {
	Tasks = mq.NewMQ()
	Tasks.NewQueue("APICallbackQueue", "callback")
	CallbackMessages = Tasks.GetMessages("APICallbackQueue")

	trip := new(TripHandler)

	router := mux.NewRouter()
	router.HandleFunc("/trips", requestHandler(trip)).Methods("GET", "POST")
	router.HandleFunc("/trip/{id}", requestHandler(trip)).Methods("GET", "PUT", "DELETE")

	http.Handle("/", handlers.CORS()(router))

	log.Println("[TP-Server] Listening on 3000")

	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)

	http.ListenAndServe(":3000", loggedRouter)
}

//ResourceHandler interface to manage requests
type ResourceHandler interface {
	Get(request *http.Request) (int, interface{})
	Post(request *http.Request) (int, interface{})
	Put(request *http.Request) (int, interface{})
	Delete(request *http.Request) (int, interface{})
}

func requestHandler(resource ResourceHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		method := request.Method

		switch method {
		case "GET":
			code, data = resource.Get(request)
		case "POST":
			code, data = resource.Post(request)
		case "PUT":
			code, data = resource.Put(request)
		case "DELETE":
			code, data = resource.Delete(request)
		default:
			code = http.StatusMethodNotAllowed
			data = nil
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			code = http.StatusBadRequest
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}
