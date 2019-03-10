package queue

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Queue struct {
	conn *amqp.Connection
	ch   []*amqp.Channel
}

func NewQueueInstance(host, port, username, password string) (*Queue, error) {
	// In here the queue, exchange and binding will not be defined
	// for avoiding orchestration problems
	amqpStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	conn, err := amqp.Dial(amqpStr)
	if err != nil {
		log.Infof("error %s occurred when dialling rabbitmq", err)
		log.Info("Will start retrying")
		counter := 0
		for err != nil {
			counter++
			time.Sleep(5 * time.Second)
			log.Infof("restart %s th times", counter)
			conn, err = amqp.Dial(amqpStr)
		}
	}
	queueClient := &Queue{}
	queueClient.conn = conn
	log.WithField("amqpStr", amqpStr).Info("rabbitmq connection started")
	return queueClient, nil
}

func (q *Queue) DestroyQueueInstance() error {
	log.Info("start to destroy queue instance")
	for _, ch := range q.ch {
		err := ch.Close()
		if err != nil {
			return err
		}
	}
	return q.conn.Close()
}

func (q *Queue) Publish(exchange, routingKey string, jsonBody map[string]interface{}) error {
	ch, err := q.conn.Channel()
	if err != nil {
		log.Errorf("channel creation failed because %v", err)
		return err
	}
	defer ch.Close()

	jsonByte, err := json.Marshal(jsonBody)
	if err != nil {
		log.Errorf("json marshal failed because %v", err)
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
		log.Errorf("rabbitmq channel publish because %v", err)
		return err
	}
	return nil
}

func (q *Queue) Consume(queue string) (<-chan amqp.Delivery, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		log.Errorf("channel creation failed because %v", err)
		return nil, err
	}
	q.ch = append(q.ch, ch)

	_, err = ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Errorf("queue declaration failed for queue %s", queue)
		return nil, err
	}

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
		log.Errorf("channel consume failed because %v", err)
		return nil, err
	}
	log.Infof("start to consume on queue %s", queue)
	return msgs, nil
}
