package usecase

import (
	"project/clean-arch/internal/entity"
	"project/clean-arch/pkg/events"
)

type GetOneOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderGetOne    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewGetOneOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderGetOne events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *GetOneOrderUseCase {
	return &GetOneOrderUseCase{
		OrderRepository: OrderRepository,
		OrderGetOne:    OrderGetOne,
		EventDispatcher: EventDispatcher,
	}
}

func (c *GetOneOrderUseCase) Execute(id string) (*entity.Order, error) {
	order ,err := c.OrderRepository.GetOne(id)
	if err != nil {
		return nil, err
	}

	c.OrderGetOne.SetPayload(order)
	c.EventDispatcher.Dispatch(c.OrderGetOne)

	return order, nil
}
