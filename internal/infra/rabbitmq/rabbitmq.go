package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type QueueInfo struct {}

type QueueData struct {
	Name     string `json:"name"`
	Messages int    `json:"messages"` 
	Ready    int    `json:"messages_ready"` 
	Unacked  int    `json:"messages_unacknowledged"`
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