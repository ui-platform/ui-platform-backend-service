package rabbit_mq

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

// Producer структура, содержащая соединение с RabbitMQ.
type Producer struct {
	conn *amqp.Connection
	log  zerolog.Logger
}

// NewProducer создает новый экземпляр Producer с подключением к RabbitMQ.
func NewProducer(url string, log zerolog.Logger) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Producer{
		conn: conn,
		log:  log,
	}, nil
}

// SendMessage отправляет сообщение в указанную очередь.
// Message может быть любым объектом, который будет сериализован в JSON.
func (p *Producer) SendMessage(queueName string, message interface{}) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Объявляем очередь, чтобы убедиться, что она существует.
	q, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Сериализация сообщения в JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Публикация сообщения в очередь.
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	p.log.Info().Msgf("Message sent to queue: %s", queueName)
	return nil
}
