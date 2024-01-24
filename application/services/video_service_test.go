package services_test

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/application/services"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
	"github.com/olukkas/go-encoder/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("could not load .env file")
	}
}

func TestVideoService(t *testing.T) {
	video, repo := prepare()

	videoService := services.VideoService{}
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("go-encoder-bucket")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDataBaseTest()
	defer db.Close()

	err := createVideosTable(db)
	utils.FailOnError(err, "error on create videos table")

	err = createJobsTable(db)
	utils.FailOnError(err, "error on create jobs table")

	return &domain.Video{
		ID:        uuid.NewV4().String(),
		FilePath:  "file.mp4", // make sure the file exists
		CreatedAt: time.Now(),
	}, repositories.NewVideoRepositoryDb(db)
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
