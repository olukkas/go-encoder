package services_test

import (
	"github.com/joho/godotenv"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/application/services"
	"github.com/olukkas/go-encoder/domain"
	"github.com/olukkas/go-encoder/framework/database"
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

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDataBaseTest()
	defer db.Close()

	return &domain.Video{
		ID:        uuid.NewV4().String(),
		FilePath:  "file.mp4", // make sure the file exists
		CreatedAt: time.Now(),
	}, repositories.NewVideoRepositoryDb(db)
}

func TestVideoService_Download(t *testing.T) {
	video, repo := prepare()

	videoService := services.VideoService{}
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("go-encoder-bucket")
	require.Nil(t, err)
}
