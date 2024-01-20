package repositories_test

import (
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobRepositoryDb_Insert(t *testing.T) {
	db := database.NewDataBaseTest()
	defer db.Close()

	video, err := domain.NewVideo("resource", "path")
	require.Nil(t, err)
	require.NotNil(t, video)

	videoRepo := repositories.NewVideoRepositoryDb(db)
	jobRepo := repositories.NewJobRepositoryDb(db)

	_, err = videoRepo.Insert(video)
	require.Nil(t, err)

	job, err := domain.NewJob("output_path", domain.JobPending, video)
	require.Nil(t, err)

	_, err = jobRepo.Insert(job)
	require.Nil(t, err)

	j, err := jobRepo.Find(job.ID)
	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoId, video.ID)
}
