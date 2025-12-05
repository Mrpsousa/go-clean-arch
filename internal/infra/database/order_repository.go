package database

import (
	"database/sql"

	"project/clean-arch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetAll() ([]entity.Order, error) {
	orders := make([]entity.Order, 0)
	rows, err := r.Db.Query("Select * from orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ord entity.Order
		err := rows.Scan(&ord.ID, &ord.Price, &ord.Tax, &ord.FinalPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (r *OrderRepository) GetOne(id string) (*entity.Order, error) {
	order := entity.Order{}
	err := r.Db.QueryRow("Select * from orders where id = ?", id).Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

