package repositories

import (
	"database/sql"
	"github.com/olukkas/go-encoder/domain"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	DB *sql.DB
}

func NewJobRepositoryDb(DB *sql.DB) *JobRepositoryDb {
	return &JobRepositoryDb{DB: DB}
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (j *JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	query := `
	insert into jobs 
	(id, output_bucket_path, status, video_id, error, created_at, updated_at) 
	values 
	($1, $2, $3, $4, $5, $6, $7)
	`
	stmt, err := j.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(job.ID, job.OutputBucketPath, job.Status, job.VideoId, job.Error, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return job, nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (j *JobRepositoryDb) Find(id string) (*domain.Job, error) {
	job := new(domain.Job)

	query := `select id, output_bucket_path, status, video_id, error, created_at, updated_at
	from jobs
	where id = ?
	`
	err := j.DB.QueryRow(query, id).
		Scan(&job.ID, &job.OutputBucketPath, &job.Status, &job.VideoId, &job.Error, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	job.Video, err = j.getJobVideo(job.VideoId)
	if err != nil {
		return nil, err
	}

	return job, nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (j *JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	query := `update jobs set status = $1, updated_at = $2, error = $3 where id = $4 `

	stmt, err := j.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(job.Status, job.UpdatedAt, job.Error, job.ID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (j *JobRepositoryDb) getJobVideo(videoId string) (*domain.Video, error) {
	var video domain.Video

	videoQuery := `select id, resource_id, file_path, created_at from videos where id = ?`
	err := j.DB.QueryRow(videoQuery, videoId).Scan(&video.ID, &video.ResourceID, &video.FilePath, &video.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &video, nil
}
