package services

import (
	"encoding/json"
	"errors"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"os"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService) {
	for message := range messageChannel {
		job, err := prepare(message, jobService)
		if err != nil {
			returnChan <- invalidJobResult(message, err)
			continue
		}

		jobService.Job = job

		err = jobService.Start()
		if err != nil {
			returnChan <- invalidJobResult(message, err)
			continue
		}

		returnChan <- JobWorkerResult{
			Job:     *job,
			Message: &message,
			Error:   nil,
		}
	}
}

func prepare(message amqp.Delivery, jobService JobService) (*domain.Job, error) {
	video := jobService.VideoService.Video

	if !utils.IsJson(string(message.Body)) {
		return nil, errors.New("message is not a json")
	}

	err := json.Unmarshal(message.Body, video)
	if err != nil {
		return nil, err
	}
	video.ID = uuid.NewV4().String()

	err = video.Validate()
	if err != nil {
		return nil, err
	}

	err = jobService.VideoService.InsertVideo()
	if err != nil {
		return nil, err
	}

	outputBucket := os.Getenv("OUTPUT_BUCKET_NAME")
	job, err := domain.NewJob(outputBucket, domain.JobStarting, jobService.VideoService.Video)
	if err != nil {
		return nil, err
	}

	_, err = jobService.JobsRepository.Insert(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func invalidJobResult(message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job:     domain.Job{},
		Message: &message,
		Error:   err,
	}
}
