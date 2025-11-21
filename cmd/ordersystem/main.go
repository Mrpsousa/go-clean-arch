package main

import (
	"database/sql"
	"fmt"

	"project/clean-arch/configs"
	"project/clean-arch/internal/event"
	"project/clean-arch/internal/event/handler"
	"project/clean-arch/internal/infra/database"
	"project/clean-arch/internal/infra/web"
	"project/clean-arch/internal/infra/web/webserver"
	"project/clean-arch/pkg/events"

	"github.com/streadway/amqp"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	// createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	// webOrderHandler := NewWebOrderHandler(db, eventDispatcher) // aqui
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreated)

	webserver.AddHandler("/order", webOrderHandler.Create)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start()

	// grpcServer := grpc.NewServer()
	// createOrderService := service.NewOrderService(*createOrderUseCase)
	// pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	// reflection.Register(grpcServer)

	// fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	// if err != nil {
	// 	panic(err)
	// }
	// go grpcServer.Serve(lis)

	// srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
	// 	CreateOrderUseCase: *createOrderUseCase,
	// }}))
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	// fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	// http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
