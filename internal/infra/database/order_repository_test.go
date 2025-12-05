package database

import (
	"database/sql"
	"fmt"
	"testing"

	"project/clean-arch/internal/entity"

	"github.com/stretchr/testify/assert"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sql.DB

func init () {
	db, err := sql.Open("sqlite3", ":memory:")
	globalDB = db
	if err != nil {
		panic(err)
	}
	
	globalDB.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
}

func TestSaveAndGetOrder(t *testing.T) {
	tests := []struct {
		id         string
		price      float64
		tax        float64
		finalPrice float64
		quantity   int
	}{
		{"1", 100.0, 10.0, 110.0, 1},
		{"2", 200.0, 20.0, 220.0, 2},
		{"3", 300.0, 30.0, 330.0, 3},
	}

	for i, tst := range tests {
		t.Run(fmt.Sprintf("test: %s", tst.id), func(t *testing.T) {		
			order, err := entity.NewOrder(tst.id, tst.price, tst.tax)
			assert.NoError(t, err)
			assert.NoError(t, order.CalculateFinalPrice())

			repo := NewOrderRepository(globalDB)
			err = repo.Save(order)
			assert.NoError(t, err)

			var orderResult []entity.Order
			orderResult, err = repo.GetAll()
			assert.NoError(t, err)
			assert.Equal(t, tst.quantity, len(orderResult))
			assert.Equal(t, tst.finalPrice, orderResult[i].FinalPrice)
			assert.Equal(t, tst.id, orderResult[i].ID)
			assert.Equal(t, tst.price, orderResult[i].Price)
			assert.Equal(t, tst.tax, orderResult[i].Tax)
		})
	}
	defer globalDB.Close()
}
