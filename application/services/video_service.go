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
	defer client.Close()

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
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID

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

func (v *VideoService) Encode() error {
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID
	cmdArgs := []string{
		path + ".frag",
		"--use-segment-timeline",
		"-o",
		path,
		"-f",
		"--exec-dir",
		"/opt/bento4/bin/",
	}

	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID

	err := os.Remove(path + ".mp4")
	if err != nil {
		log.Println("error removing mp4", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(path + ".frag")
	if err != nil {
		log.Println("error removing frag", v.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		log.Println("error removing directory", v.Video.ID)
		return err
	}

	log.Println("files has been removed: ", v.Video.ID)

	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	return err
}

//goland:noinspection GoDeprecation
func (v *VideoService) writeInDisk(reader *storage.Reader) error {
	storagePath := os.Getenv("LOCAL_STORAGE_PATH")

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
