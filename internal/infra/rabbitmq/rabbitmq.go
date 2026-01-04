package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

type QueueInfo struct {}

type QueueData struct {
	Ready    int    `json:"messages_ready"` 
}

type RabbitMq struct {
	Conn *amqp.Connection
}

type RabbitMsg struct {
	CreatedAt    time.Time
	ExameName    string
	PacienteName string
	DocImagePath string
}

func(q *QueueInfo) getQueueMessages(queueName string) (int, error) {
	url := fmt.Sprintf("http://localhost:15672/api/queues/%%2F/%s", queueName)

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("guest", "guest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro HTTP: %d", resp.StatusCode)
	}

	var queueData QueueData
	if err := json.NewDecoder(resp.Body).Decode(&queueData); err != nil {
		return 0, err
	}

	return queueData.Ready, nil 
}

func(q *QueueInfo) Run() (int, error) {
	count, err := q.getQueueMessages("imagem-queue")
	if err != nil {
		fmt.Println("Erro:", err)
		return 0, err
	}
	return count, nil
}


func NewRabbitMq(conn *amqp.Connection) *RabbitMq {
	return &RabbitMq{Conn: conn}
}

func (r *RabbitMq) Receiver(routingKey, queueName, exchangeName string) (*RabbitMsg, error) {
	rabbitMsg := &RabbitMsg{}

	defer r.Conn.Close()

	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// Declara a exchange
	err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Declara a fila
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Faz o bind da fila na exchange com a routing key "X"
	err = ch.QueueBind(q.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	// Configura QoS para consumir 1 mensagem por vez
	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	// Consome (autoAck = false para confirmar manualmente)
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Recebe APENAS UMA mensagem
	msg, ok := <-msgs
	if !ok {
		return nil, fmt.Errorf("some queue error, can't receive message 'no ok'")	
	}

	err = json.Unmarshal(msg.Body, rabbitMsg)
	if err != nil {
		return nil, err
	}

	// Confirma o recebimento (ack)
	msg.Ack(false)
	
	return rabbitMsg, nil
}