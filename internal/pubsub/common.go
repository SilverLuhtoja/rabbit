package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		return errors.New("PublishJSON: couldn't marshal value")
	}

	ctx := context.Background()
	err = ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonData})
	if err != nil {
		return err
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType bool, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("DeclareAndBind: Coulnt start channel")
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	queue, err := channel.QueueDeclare(queueName, simpleQueueType, simpleQueueType, simpleQueueType, false, nil)
	if err != nil {
		log.Fatal("DeclareAndBind: Coulnt declare queue")
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	channel.QueueBind(queue.Name, key, exchange, false, nil)

	return channel, queue, nil
}
