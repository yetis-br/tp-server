package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yetis-br/tp-server/models"
	"github.com/yetis-br/tp-server/mq"
	"github.com/yetis-br/tp-server/util"

	db "github.com/dancannon/gorethink"
)

func init() {
	util.AppendConfigFile("config.ini")
}

var session *db.Session

func main() {

	dbAddress := util.GetKeyValue("RethinkDB", "address")

	var err error
	session, err = db.Connect(db.ConnectOpts{
		Address:  dbAddress,
		Database: "travelPlanning",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected to db on: " + dbAddress)

	tasks := mq.NewMQ()
	tasks.NewQueue("TripWorkerQueue", "Trip")
	msgs := tasks.GetMessages("TripWorkerQueue")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var message mq.Message
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
	case "GET":
		getTrip(message)
	case "POST":
		postTrip(message)
	}
}

func getAllTrips(message *mq.Message) {
	trips, err := db.Table("Trips").Run(session)
	util.LogOnError(err, "[Trip Worker] Erro on getAllTrips")

	var rows []interface{}
	err = trips.All(&rows)
	util.LogOnError(err, "[Trip Worker] Erro on All function")

	message.ResponseCode = http.StatusOK
	message.Response = rows

	trips.Close()
}

func getTrip(message *mq.Message) {
	var trip models.Trip
	trip.ID = message.Request.(string)

	result, err := db.Table("Trips").Get(trip.ID).Run(session)
	util.LogOnError(err, "[Trip Worker] Erro on getTrip")

	if result == nil {
		message.ResponseCode = http.StatusNotFound
		message.Response = nil
	} else {
		err = result.One(&trip)
		util.LogOnError(err, "[Trip Worker] Erro on One function")

		message.ResponseCode = http.StatusOK
		message.Response = trip
	}
	result.Close()
}

func postTrip(message *mq.Message) {
	var trip models.Trip
	trip.Initialize()
	trip.LoadJSON(message.Request.(string))

	if trip.Validate() {
		resp, err := db.Table("Trips").Insert(trip).RunWrite(session)
		util.LogOnError(err, "[Trip Worker] Erro inserting new Trip")

		message.ResponseCode = http.StatusOK
		message.Response = resp
	} else {
		message.ResponseCode = http.StatusBadRequest
		message.Response = nil
	}
}
