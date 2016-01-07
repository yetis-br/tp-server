package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"github.com/yetis-br/tp-server/mq"
)

//TripHandler defines the methods to return the response
type TripHandler struct{}

//Get manages get method requests
func (t TripHandler) Get(request *http.Request) (int, interface{}) {
	var message mq.Message
	message.CorrelationID = uuid.NewV4().String()
	message.RequestAction = "GET_ALL"

	Tasks.PublishMessage(message, "Trip", message.CorrelationID, "callback")

	return processMessage(&message)
}

//Post manages post method requests
func (t TripHandler) Post(request *http.Request) (int, interface{}) {
	return http.StatusOK, nil
}

//Put manages put method requests
func (t TripHandler) Put(request *http.Request) (int, interface{}) {
	return http.StatusOK, nil
}

//Delete manages delete method requests
func (t TripHandler) Delete(request *http.Request) (int, interface{}) {
	return http.StatusOK, nil
}

func processMessage(message *mq.Message) (int, interface{}) {
	timer := time.NewTimer(time.Second * 5)

	for {
		select {
		case <-timer.C:
			//Get data from cache instead of worker
			return http.StatusInternalServerError, nil
		case d := <-CallbackMessages:
			if message.CorrelationID == d.CorrelationId {
				json.Unmarshal(d.Body, &message)
				return http.StatusOK, message.Response
			}
		}
	}
}
