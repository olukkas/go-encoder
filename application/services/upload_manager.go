package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)

	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	return vu.walkDir(vu.VideoPath)
}

func (vu *VideoUpload) walkDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			// Recursively walk through subdirectories
			if err := vu.walkDir(path); err != nil {
				return err
			}
		} else {
			// Add file path to Paths
			vu.Paths = append(vu.Paths, path)
		}
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, done chan string) error {
	in := make(chan string, runtime.NumCPU())
	returnChan := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	client, ctx, err := getClient()
	if err != nil {
		return err
	}

	for i := 0; i < concurrency; i++ {
		go vu.uploadWorker(in, returnChan, client, ctx)

	}

	go func() {
		for _, path := range vu.Paths {
			in <- path
		}
		close(in)
	}()

	for r := range returnChan {
		if r != "" {
			done <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(
	in chan string,
	returnChan chan string,
	uploadClient *storage.Client,
	ctx context.Context,
) {
	for path := range in {
		err := vu.UploadObject(path, uploadClient, ctx)
		if err != nil {
			vu.Errors = append(vu.Errors, path)
			log.Printf("error during upload %s. Error: %s", path, err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func getClient() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
