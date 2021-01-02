package main

import (
	"fmt"
	"log"
	"os"

	queuemanager "github.com/rodrwan/bank/pkg/queueManager"
	"github.com/rodrwan/bank/pkg/services/accounts"
	"github.com/rodrwan/bank/pkg/services/users"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	fmt.Println("running denormalizer")
	qq, err := queuemanager.MakeQueueManager(os.Getenv("RABBIT_MQ_URL"))
	if err != nil {
		failOnError(err, "could not connect to queue")
	}

	forever := make(chan bool)

	createdAccountEvent := queuemanager.NewEvent(
		accounts.CreatedAccountEvent,
		"accounts-svc",
		nil,
	)
	go qq.Subscribe(createdAccountEvent, func(event queuemanager.Event) {
		fmt.Printf(
			"data from %s: %s at %s\n",
			event.From,
			string(event.Data),
			event.CreatedAt.String(),
		)
	})

	createdUserEvent := queuemanager.NewEvent(
		users.CreatedUserEvent,
		"users-svc",
		nil,
	)
	go qq.Subscribe(createdUserEvent, func(event queuemanager.Event) {
		fmt.Printf(
			"data from %s: %s at %s\n",
			event.From,
			string(event.Data),
			event.CreatedAt.String(),
		)
	})

	<-forever
}
