package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetAll() ([]Order, error)
	GetOne(id string) (*Order, error)
	Update(order *Order) error
}
