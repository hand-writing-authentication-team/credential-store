package queue

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/streadway/amqp"
)

type Queue struct {
	conn *amqp.Connection
}

func NewQueueInstance(host, port, username, password string) (*Queue, error) {
	// In here the queue, exchange and binding will not be defined
	// for avoiding orchestration problems
	amqpStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	conn, err := amqp.Dial(amqpStr)
	if err != nil {
		glog.Fatalf("error %s occurred when dialling rabbitmq", err)
		return nil, err
	}
	queueClient := &Queue{}
	queueClient.conn = conn
	return queueClient, nil
}

func (q *Queue) DestroyQueueInstance() error {
	return q.conn.Close()
}

func (q *Queue) Publish(exchange, routingKey string, jsonBody map[string]interface{}) error {
	ch, err := q.conn.Channel()
	if err != nil {
		glog.Errorf("channel creation failed because %v", err)
		return err
	}
	defer ch.Close()

	jsonByte, err := json.Marshal(jsonBody)
	if err != nil {
		glog.Errorf("json marshal failed because %v", err)
		return err
	}
	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		true,       // mandatory
		true,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonByte),
		})
	if err != nil {
		glog.Errorf("rabbitmq channel publish because %v", err)
		return err
	}
	return nil
}

func (q *Queue) Consume(queue string) (<-chan amqp.Delivery, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		glog.Errorf("channel creation failed because %v", err)
		return nil, err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		glog.Errorf("channel consume failed because %v", err)
		return nil, err
	}
	return msgs, nil
}
