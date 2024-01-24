package repositories

import (
	"database/sql"
	"github.com/olukkas/go-encoder/domain"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
	DB *sql.DB
}

func NewVideoRepositoryDb(DB *sql.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{DB: DB}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (v *VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	query := "insert into videos (id, resource_id, file_path, created_at) values ($1, $2, $3, $4);"
	stmt, err := v.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(video.ID, video.ResourceID, video.FilePath, video.CreatedAt)
	if err != nil {
		return nil, err
	}

	return video, nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (v *VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video

	videoQuery := `select id, resource_id, file_path, created_at from videos where id = ?`
	err := v.DB.QueryRow(videoQuery, id).Scan(&video.ID, &video.ResourceID, &video.FilePath, &video.CreatedAt)
	if err != nil {
		return nil, err
	}

	video.Jobs, err = v.getVideoJobs(&video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

//goland:noinspection SqlNoDataSourceInspection,GoConvertStringLiterals,SqlResolve
func (v *VideoRepositoryDb) getVideoJobs(video *domain.Video) ([]*domain.Job, error) {
	var jobs []*domain.Job

	jobsQuery := `
	select id, output_bucket_path, status, video_id, error, created_at, updated_at 
	from jobs where video_id = ?
	`

	rows, err := v.DB.Query(jobsQuery, video.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var job domain.Job
		err = rows.Scan(&job.ID, &job.OutputBucketPath, &job.Status, &job.Error, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}

		job.Video = video
		jobs = append(jobs, &job)
	}

	return jobs, err
}
