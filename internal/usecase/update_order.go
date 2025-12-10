package usecase

import (
	"project/clean-arch/internal/entity"
	"project/clean-arch/pkg/events"
)

type UpdateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderUpdate    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewUpdateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderUpdate events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderUpdate:    OrderUpdate,
		EventDispatcher: EventDispatcher,
	}
}

func (c *UpdateOrderUseCase) Execute(order *entity.Order) error {
	err := c.OrderRepository.Update(order)
	if err != nil {
		return err
	}

	c.OrderUpdate.SetPayload(order)
	c.EventDispatcher.Dispatch(c.OrderUpdate)

	return nil
}
