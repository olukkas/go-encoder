package services_test

import (
	"github.com/joho/godotenv"
	"github.com/olukkas/go-encoder/application/services"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

//goland:noinspection GoUnhandledErrorResult
func TestVideoServiceUpload(t *testing.T) {
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

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "go-encoder-bucket"
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(10, doneUpload)

	result := <-doneUpload
	require.Len(t, videoUpload.Errors, 0)
	require.Equal(t, result, "upload completed")

	err = videoService.Finish()
	require.Nil(t, err)
}
