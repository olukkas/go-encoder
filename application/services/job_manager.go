package services

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/framework/queue"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
)

type JobManager struct {
	DB             *gorm.DB
	MessageChannel chan amqp.Delivery
	ReturnChannel  chan JobWorkerResult
	RabbitMQ       *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(
	DB *gorm.DB,
	rabbitMq *queue.RabbitMQ,
	returnChannel chan JobWorkerResult,
	messageChannel chan amqp.Delivery,
) *JobManager {
	return &JobManager{
		DB:             DB,
		MessageChannel: messageChannel,
		ReturnChannel:  returnChannel,
		RabbitMQ:       rabbitMq,
	}
}

func (j *JobManager) Start() {
	videoService := VideoService{}
	videoService.VideoRepository = repositories.NewVideoRepositoryDb(j.DB)

	jobService := JobService{
		JobsRepository: repositories.NewJobRepositoryDb(j.DB),
		VideoService:   videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS")
	}

	for qtdPrecess := 0; qtdPrecess < concurrency; qtdPrecess++ {
		go JobWorker(j.MessageChannel, j.ReturnChannel, jobService)
	}

	for jobResult := range j.ReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult)
		}

		if err != nil {
			_ = jobResult.Message.Reject(false)
		}

	}
}

func (j *JobManager) notifySuccess(result JobWorkerResult) error {
	jobJson, err := json.Marshal(result.Job)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = result.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) checkParseErrors(result JobWorkerResult) error {
	if result.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v",
			result.Message.DeliveryTag, result.Job.ID, result.Job.Video.ID, result.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", result.Message.DeliveryTag, result.Error)
	}

	errorMsg := JobNotificationError{
		Message: string(result.Message.Body),
		Error:   result.Error.Error(),
	}

	jobJson, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = result.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}
func (j *JobManager) notify(jobJson []byte) error {
	return j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
}
