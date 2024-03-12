package main

import (
	"context"
	"database/sql"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fabioods/fc-ms-wallet/internal/database"
	"github.com/fabioods/fc-ms-wallet/internal/event"
	"github.com/fabioods/fc-ms-wallet/internal/event/handler"
	"github.com/fabioods/fc-ms-wallet/internal/usecase/create_account"
	"github.com/fabioods/fc-ms-wallet/internal/usecase/create_client"
	"github.com/fabioods/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/fabioods/fc-ms-wallet/internal/web"
	"github.com/fabioods/fc-ms-wallet/internal/web/webserver"
	"github.com/fabioods/fc-ms-wallet/pkg/events"
	"github.com/fabioods/fc-ms-wallet/pkg/kafka"
	"github.com/fabioods/fc-ms-wallet/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@tcp(mysql:3306)/wallet?parseTime=true"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	// Verifica a conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)
	fmt.Println("Connected to kafka")

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("transaction_created", handler.NewTransactionCreatedKafka(kafkaProducer))
	eventDispatcher.Register("balance_updated", handler.NewBalanceUpdatedKafka(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	fmt.Println("Events registered")

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	unitOfWork := uow.NewUow(ctx, db)

	unitOfWork.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	unitOfWork.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(unitOfWork, transactionCreatedEvent, balanceUpdatedEvent, eventDispatcher)

	fmt.Println("Use cases created")

	webServer := webserver.NewWebServer("0.0.0.0:8080")

	fmt.Println("Web server created")

	clientHandler := web.NewClientHandlerWeb(*createClientUseCase)
	accountHandler := web.NewAccountHandlerWeb(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandlerWeb(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClientHandlerWeb)
	webServer.AddHandler("/accounts", accountHandler.CreateAccountHandlerWeb)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransactionHandlerWeb)
	webServer.AddHandler("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	fmt.Println("Handlers added")

	fmt.Println("Server started on port 8080")
	err = webServer.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Server stopped")
	// Create table transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date)
	// Create table accounts (id varchar(255), client_id varchar(255), balance int, created_at date)
	// Create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)

}
