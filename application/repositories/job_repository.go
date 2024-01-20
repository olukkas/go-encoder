package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olukkas/go-encoder/domain"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	DB *gorm.DB
}

func NewJobRepositoryDb(DB *gorm.DB) *JobRepositoryDb {
	return &JobRepositoryDb{DB: DB}
}

func (j *JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	if err := j.DB.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (j *JobRepositoryDb) Find(id string) (*domain.Job, error) {
	job := new(domain.Job)
	j.DB.Preload("Video").First(job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("job with id %s does not exists \n", id)
	}

	return job, nil
}

func (j *JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := j.DB.Save(&job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}
