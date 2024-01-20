package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type JobStatus string

const (
	JobDownloading JobStatus = "DOWNLOADING"
)

type Job struct {
	ID               string    `valid:"uuid"`
	OutputBucketPath string    `valid:"notnull"`
	Status           JobStatus `valid:"notnull"`
	Video            *Video    `valid:"-"`
	VideoId          string    `valid:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"-"`
	UpdatedAt        time.Time `valid:"-"`
}

func NewJob(output string, status JobStatus, video *Video) (*Job, error) {
	job := &Job{
		ID:               uuid.NewV4().String(),
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	return job, nil
}

func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)
	if err != nil {
		return err
	}

	return nil
}
