package main

import (
	"github.com/joho/godotenv"
	"github.com/olukkas/go-encoder/application/services"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/olukkas/go-encoder/framework/queue"
	"github.com/olukkas/go-encoder/framework/utils"
	"github.com/streadway/amqp"
	"os"
)

var db database.DataBase

func init() {
	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	err = configDb()
	utils.FailOnError(err, "erro configuring Database")
}

func configDb() error {
	db.Dsn = os.Getenv("DSN")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")

	return nil
}

//goland:noinspection GoUnhandledErrorResult
func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConn, err := db.Connect()
	utils.FailOnError(err, "fail connecting to database")

	defer db.Connect()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	manager := services.NewJobManager(dbConn, rabbitMQ, jobReturnChannel, messageChannel)
	manager.Start()
}
