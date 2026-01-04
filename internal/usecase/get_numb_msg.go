package usecase

import (
	rbmq "project/clean-arch/internal/infra/rabbitmq"
)

type GetNumbMsgInQueueUseCase struct {
	Rabbit *rbmq.RabbitMq
}

func NewNumbMsgInQueueUseCase(rabbit *rbmq.RabbitMq) *GetNumbMsgInQueueUseCase {
	return &GetNumbMsgInQueueUseCase{
		Rabbit: rabbit,
	}
}

func (c *GetNumbMsgInQueueUseCase) Execute() (int, error) {
	numb, err := c.Rabbit.Run()
	if err != nil {
		return 0, err
	}

	return numb, nil
}
