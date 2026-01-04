package web

import (
	"encoding/json"
	"net/http"

	"project/clean-arch/internal/entity"
	rbmq "project/clean-arch/internal/infra/rabbitmq"
	"project/clean-arch/internal/usecase"
	"project/clean-arch/pkg/events"

	"github.com/go-chi/chi/v5"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	RabbitInfo *rbmq.QueueInfo
	RabbitMq 	 *rbmq.RabbitMq
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
	RabbitInfo *rbmq.QueueInfo,
	RabbitMq 	 *rbmq.RabbitMq,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
		RabbitInfo:            RabbitInfo,
		RabbitMq: RabbitMq,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	getOrders := usecase.NewGetAllOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := getOrders.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	if orderID == "" {
		http.Error(w, "missing ID ", http.StatusBadRequest)
		return
	}

	getOrder := usecase.NewGetOneOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := getOrder.Execute(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) GetNumbMsgInQueue(w http.ResponseWriter, r *http.Request) {
	getNumbMsgs := usecase.NewNumbMsgInQueueUseCase(h.RabbitInfo)
	output, err := getNumbMsgs.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) GetRabbitMqMsg(w http.ResponseWriter, r *http.Request) {
	getNumbMsgs := usecase.NewGetRabbitMsgUseCase(h.RabbitMq)
	output, err := getNumbMsgs.Execute("exame_imagem", "test-queue", "general_channel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
