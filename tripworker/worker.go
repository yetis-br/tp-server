package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yetis-br/tp-server/models"
	"github.com/yetis-br/tp-server/mq"
	"github.com/yetis-br/tp-server/util"
)

func init() {
	util.AppendConfigFile("config.ini")
}

func main() {
	log.Println("Connected to db on: " + util.GetKeyValue("RethinkDB", "address"))

	tasks := mq.NewMQ()
	tasks.NewQueue("TripWorkerQueue", "Trip")
	msgs := tasks.GetMessages("TripWorkerQueue")

	var message mq.Message

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			json.Unmarshal(d.Body, &message)
			processMessage(&message)
			tasks.PublishMessage(message, d.ReplyTo, d.CorrelationId, "")
		}
	}()

	log.Printf(" [*] Awaiting for requests")
	<-forever

}

func processMessage(message *mq.Message) {
	switch message.RequestAction {
	case "GET_ALL":
		getAllTrips(message)
	}
}

func getAllTrips(message *mq.Message) {
	trip := new(models.Trip)
	trip.ID = "00000001"
	trip.Title = "Testando Europa 2016"

	message.ResponseCode = http.StatusOK
	message.Response = trip
}
