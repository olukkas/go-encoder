package repositories_test

import (
	"database/sql"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVideoRepositoryDB_Insert(t *testing.T) {
	db := database.NewDataBaseTest()
	defer db.Close()

	err := createVideosTable(db)
	require.Nil(t, err)

	err = createJobsTable(db)
	require.Nil(t, err)

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

//goland:noinspection SqlNoDataSourceInspection
func createVideosTable(db *sql.DB) error {
	ddl := `
	create table videos (
	    id text,
		resource_id text,
		file_path text,
		created_at date
	)
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return err
	}

	return nil
}

//goland:noinspection SqlNoDataSourceInspection
func createJobsTable(db *sql.DB) error {
	ddl := `
	create table jobs (
	    id text,
		output_bucket_path text,
		status text,
		video_id text,
		error text,
		created_at date,
		updated_at date
	)
	`
	_, err := db.Exec(ddl)
	if err != nil {
		return err
	}

	return nil
}
