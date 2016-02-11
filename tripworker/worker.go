package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	case "POST":
		postTrip(message)
	}
}

func getAllTrips(message *mq.Message) {
	trips, err := db.Table("Trips").Run(session)
	if err != nil {
		log.Print(err)
	}
	var rows []interface{}
	err = trips.All(&rows)
	if err != nil {
		log.Print(err)
	}

	message.ResponseCode = http.StatusOK
	message.Response = rows

	trips.Close()
}

func postTrip(message *mq.Message) {
	var trip models.Trip
	jsonTrip := message.Request.(string)
	json.Unmarshal([]byte(jsonTrip), &trip)

	trip.CreatedDate = time.Now()
	trip.UpdatedDate = time.Now()

	resp, err := db.Table("Trips").Insert(trip).RunWrite(session)
	if err != nil {
		log.Print(err)
	}

	message.ResponseCode = http.StatusOK
	message.Response = resp
}
