package queuemanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/twinj/uuid"
)

const (
	defaultTimeout = time.Second
	defaultQOS     = 0
)

// Event ...
type Event struct {
	ID        string    `json:"id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Data      []byte    `json:"data,omitempty"`
	From      string    `json:"from,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// NewEvent ...
func NewEvent(topic, from string, data []byte) Event {
	return Event{
		ID:        uuid.NewV1().String(),
		Type:      topic,
		Data:      data,
		CreatedAt: time.Now(),
		From:      from,
	}
}

// QueueManager represent operation for a queue manager.
type QueueManager interface {
	Subscribe(event Event, cb func(event Event))
	Publish(event Event) error
}

type mqttQueueManager struct {
	client *amqp.Connection
	sync.RWMutex
}

// MakeQueueManager instantiate a new connection to rabbitmq broker.
func MakeQueueManager(dsn string) (QueueManager, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	return &mqttQueueManager{
		client: conn,
	}, nil
}

func (q *mqttQueueManager) Subscribe(event Event, cb func(event Event)) {
	ch, err := q.client.Channel()
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to open a channel"))
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"topics_exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to declare an exchange"))
		return
	}

	qq, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to declare a queue"))
		return
	}

	if err := ch.QueueBind(
		qq.Name,           // queue name
		event.Type,        // routing key
		"topics_exchange", // exchange
		false,
		nil); err != nil {
		fmt.Println(errors.Wrap(err, "Failed to bind a queue"))
		return
	}

	msgs, err := ch.Consume(
		qq.Name, // queue
		"",      // consumer
		true,    // auto ack
		false,   // exclusive
		false,   // no local
		false,   // no wait
		nil,     // args
	)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to register a consumer"))
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			event.Data = d.Body
			cb(event)
		}
	}()

	<-forever
}

func (q *mqttQueueManager) Publish(event Event) error {
	q.RLock()
	defer q.RUnlock()

	if q.client != nil {
		ch, err := q.client.Channel()
		if err != nil {
			return errors.Wrap(err, "Failed to open a channel")
		}
		defer ch.Close()

		if err := ch.ExchangeDeclare(
			"topics_exchange", // name
			"topic",           // type
			true,              // durable
			false,             // auto-deleted
			false,             // internal
			false,             // no-wait
			nil,               // arguments
		); err != nil {
			return errors.Wrap(err, "Failed to declare an exchange")
		}

		if err := ch.Publish(
			"topics_exchange", // exchange
			event.Type,        // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        event.Data,
			}); err != nil {
			return errors.Wrap(err, "Failed to declare an exchange")

		}
		return nil
	}

	return fmt.Errorf("not connected")
}
