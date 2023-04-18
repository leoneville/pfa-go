package database

import (
	"database/sql"

	"github.com/leoneville/pfa-go/internal/order/entity"
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
	defer stmt.Close()

	if _, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	if err := r.Db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
