1 - adjust rabbit (infra) module
    - adjust return to struct insted of "int"
    type QueueData struct {
	Name     string `json:"name"`
	Messages int    `json:"messages"` 
	Ready    int    `json:"messages_ready"` 
	Unacked  int    `json:"messages_unacknowledged"`
}
2 - rabbitmq need to be a different module