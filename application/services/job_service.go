package services

import (
	"errors"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"os"
	"strconv"
)

type JobService struct {
	Job            *domain.Job
	JobsRepository repositories.JobRepository
	VideoService   VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStatus(domain.JobDownloading)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("INPUT_BUCKET_NAME"))
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobFragmenting)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobEncoding)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobFinishing)
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobCompleted)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus(domain.JobUploading)
	if err != nil {
		return j.failJob(err)
	}

	vu := NewVideoUpload()
	vu.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
	vu.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	//goland:noinspection GoUnhandledErrorResult
	go vu.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStatus(status domain.JobStatus) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobsRepository.Update(j.Job)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = domain.JobFailed
	j.Job.Error = error.Error()

	_, err := j.JobsRepository.Update(j.Job)
	if error != nil {
		return err
	}

	return error
}
