package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Configurações do RabbitMQ (ajuste conforme seu ambiente)
const (
	rabbitURL   = "amqp://guest:guest@localhost:5672/" // ou sua URL real
	exchange    = "general_channel"
	routingKey  = "exame_imagem"
	intervalSec = 60
)

type RabbitMsg struct {
	CreatedAt    time.Time
	ExameName    string
	PacienteName string
	DocImagePath string
}

func main() {
	// Conexão ao RabbitMQ com reconexão automática
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Abre um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir canal: %v", err)
	}
	defer ch.Close()

	// Declara a exchange (se não existir)
	err = ch.ExchangeDeclare(
		exchange, // nome
		"topic",  // tipo
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Falha ao declarar exchange: %v", err)
	}

	log.Printf("Exchange '%s' declarada. Enviando mensagens a cada %d segundos com routing key '%s'", exchange, intervalSec, routingKey)

	// Loop infinito: envia mensagem a cada 60 segundos
	ticker := time.NewTicker(intervalSec * time.Second)
	defer ticker.Stop()

	count := 1
	for {
		select {
		case <-ticker.C:
			rabbitMsg := RabbitMsg{
				CreatedAt:    time.Now(),
				ExameName:    "Image Exame",
				PacienteName: "Pacient Test",
				DocImagePath: "/path/to/image.jpg",
			}
			bytesMsg, err := json.Marshal(rabbitMsg)
			if err != nil {
				log.Fatal(err)
			}
			err = ch.Publish(
				exchange,   // exchange
				routingKey, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(bytesMsg),
				},
			)
			if err != nil {
				log.Printf("Falha ao enviar mensagem: %v", err)
				// Tenta reconectar se necessário
				conn, err = connectRabbitMQ()
				if err != nil {
					log.Printf("Falha na reconexão: %v", err)
					continue
				}
				ch, _ = conn.Channel()
			} else {
				log.Printf("Mensagem enviada")
				count++
			}
		}
	}
}

// Função auxiliar para conectar com retry
func connectRabbitMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 5; i++ { // Tenta 5 vezes
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			return conn, nil
		}
		log.Printf("Tentativa %d falhou: %v. Tentando novamente em 5s...", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("falha ao conectar após 5 tentativas: %v", err)
}