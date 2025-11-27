package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatingNewOrderCreated(t *testing.T) {
	order := NewOrderCreated()
	assert.NotNil(t, order)
	assert.Equal(t, "OrderCreated", order.GetName())

}

func TestOrderCreated_SetAndGetPayload(t *testing.T) {
	order := NewOrderCreated()
	payload := map[string]interface{}{
		"id":    "123",
		"price": 10.0,
		"tax":   2.0,
	}
	order.SetPayload(payload)
	assert.Equal(t, payload, order.GetPayload())
}