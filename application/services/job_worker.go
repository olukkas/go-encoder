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
		if !utils.IsJson(string(message.Body)) {
			returnChan <- returnJobResult(domain.Job{}, message, errors.New("message is not a json"))
			continue
		}

		err := json.Unmarshal(message.Body, &jobService.VideoService.Video)
		jobService.VideoService.Video.ID = uuid.NewV4().String()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		bucket := os.Getenv("OUTPUT_BUCKET_NAME")
		job, err := domain.NewJob(bucket, domain.JobStarting, jobService.VideoService.Video)
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		_, err = jobService.JobsRepository.Insert(job)

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = job
		err = jobService.Start()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(*job, message, nil)

	}

}
func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}

	return result
}
