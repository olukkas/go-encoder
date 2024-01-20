package services

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/olukkas/go-encoder/application/repositories"
	"github.com/olukkas/go-encoder/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var storagePath = os.Getenv("LOCAL_STORAGE_PATH")

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	client.Close()

	r, err := client.Bucket(bucketName).Object(v.Video.FilePath).NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	err = v.writeInDisk(r)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoService) Fragment() error {
	path := storagePath + "/" + v.Video.ID

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	source := path + ".mp4"
	target := path + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

//goland:noinspection GoDeprecation
func (v *VideoService) writeInDisk(reader *storage.Reader) error {
	body, err := ioutil.ReadAll(reader)

	f, err := os.Create(storagePath + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	log.Printf("vídeo %s was created", v.Video.ID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=======> Output: %s\n", string(out))
	}
}
