package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
)

type Connection struct {
	conn *amqp.Connection
}

func (conn Connection) Publish(queue string, data interface{}) error {
	channel, err := conn.conn.Channel()

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = channel.Publish(`main`, `main.`+queue, false, false, amqp.Publishing{
		Body:            jsonData,
		ContentEncoding: `utf8`,
		ContentType:     `application/json`,
	})

	if err != nil {
		return err
	}

	return nil
}

func (conn Connection) Subscribe(queue string, dataPrototype interface{}, handler func(data interface{}) error) {

}

func (conn Connection) Channel() (*amqp.Channel, error) {
	return conn.conn.Channel()
}

func (conn Connection) Close() error {
	return conn.conn.Close()
}

func Connect() Connection {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))

	if err != nil {
		panic(err)
	}

	return Connection{conn: conn}
}
