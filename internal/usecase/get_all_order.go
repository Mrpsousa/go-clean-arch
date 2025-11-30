package usecase

import (
	"project/clean-arch/internal/entity"
	"project/clean-arch/pkg/events"
)

type OrderGetAllOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type GetAllOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderGetAll    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewGetAlleOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderGetAll events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *GetAllOrderUseCase {
	return &GetAllOrderUseCase{
		OrderRepository: OrderRepository,
		OrderGetAll:    OrderGetAll,
		EventDispatcher: EventDispatcher,
	}
}

func (c *GetAllOrderUseCase) Execute() ([]entity.Order, error) {
	orders ,err := c.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	c.OrderGetAll.SetPayload(orders)
	c.EventDispatcher.Dispatch(c.OrderGetAll)

	return orders, nil
}
