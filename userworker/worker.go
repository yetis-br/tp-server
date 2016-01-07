package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yetis-br/tp-server/mq"
	"github.com/yetis-br/tp-server/util"
)

func init() {
	util.AppendConfigFile("config.ini")
}

func main() {
	log.Println("Connected to db on: " + util.GetKeyValue("RethinkDB", "address"))

	tasks := mq.NewMQ()
	tasks.NewQueue("UserWorkerQueue", "User")
	msgs := tasks.GetMessages("UserWorkerQueue")

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
	case "LOGIN":
		loginUser(message)
	}
}

func loginUser(message *mq.Message) {
	message.ResponseCode = http.StatusOK
	message.Response = nil
}
