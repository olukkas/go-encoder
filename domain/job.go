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
	JobPending     JobStatus = "PENDING"
	JobFailed      JobStatus = "FAILED"
	JobFragmenting JobStatus = "FRAGMENTING"
	JobEncoding    JobStatus = "ENCODING"
	JobFinishing   JobStatus = "FINISHING"
	JobUploading   JobStatus = "UPLOADING"
	JobCompleted   JobStatus = "COMPLETED"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           JobStatus `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoId          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `json:"-" valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
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
