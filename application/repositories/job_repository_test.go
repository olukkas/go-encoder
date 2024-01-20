package repositories_test

import (
	"github.com/jinzhu/gorm"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobRepositoryDb_Insert(t *testing.T) {
	db := database.NewDataBaseTest()
	defer db.Close()

	jobRepo := repositories.NewJobRepositoryDb(db)

	job := prepareJobHelper(t, db)

	_, err := jobRepo.Insert(job)
	require.Nil(t, err)

	j, err := jobRepo.Find(job.ID)
	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.ID, job.ID)
	require.NotEmpty(t, j.VideoId)
}

func TestJobRepositoryDb_Update(t *testing.T) {
	db := database.NewDataBaseTest()
	defer db.Close()

	jobRepo := repositories.NewJobRepositoryDb(db)
	job := prepareJobHelper(t, db)

	_, err := jobRepo.Insert(job)
	require.Nil(t, err)

	job.Status = domain.JobDownloading
	updated, err := jobRepo.Update(job)
	require.Nil(t, err)
	require.Equal(t, domain.JobDownloading, updated.Status)
}

func prepareJobHelper(t *testing.T, db *gorm.DB) *domain.Job {
	video, err := domain.NewVideo("resource", "path")
	require.Nil(t, err)
	require.NotNil(t, video)

	videoRepo := repositories.NewVideoRepositoryDb(db)

	_, err = videoRepo.Insert(video)
	require.Nil(t, err)

	job, err := domain.NewJob("output_path", domain.JobPending, video)
	require.Nil(t, err)

	return job
}
