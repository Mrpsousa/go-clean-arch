package usecase

import (
	"project/clean-arch/internal/entity"
	rbmq "project/clean-arch/internal/infra/rabbitmq"
)

type GetRabbitMsgUseCase struct {
	RabbitMq *rbmq.RabbitMq
}

func NewGetRabbitMsgUseCase(rabbit *rbmq.RabbitMq) *GetRabbitMsgUseCase {
	return &GetRabbitMsgUseCase{
		RabbitMq: rabbit,
	}
}

func (c *GetRabbitMsgUseCase) Execute(routingKey, queueName, exchange string) (*entity.RabbitMsg, error) {
	msg, err := c.RabbitMq.Receiver(routingKey, queueName, exchange)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
