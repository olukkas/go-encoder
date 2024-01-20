package repositories_test

import (
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVideoRepositoryDB_Insert(t *testing.T) {
	db := database.NewDataBaseTest()
	defer db.Close()

	video, err := domain.NewVideo("resource", "path")
	require.Nil(t, err)
	require.NotNil(t, video)

	repo := repositories.NewVideoRepositoryDb(db)
	_, err = repo.Insert(video)
	require.Nil(t, err)

	v, err := repo.Find(video.ID)
	require.Nil(t, err)
	require.NotEmpty(t, v.ID)
	require.Equal(t, v.ID, video.ID)
}
