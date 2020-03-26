package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Storage is a struct for storage config
type Storage struct {
	StorageBucketName   string
	StorageFolderName   string
	FirebaseCredentials string
}

// Upload uploads a file to cloud storage
func (s *Storage) Upload(fileName string) (string, error) {

	baseURL := "https://storage.cloud.google.com/" + s.StorageBucketName + "/" + s.StorageFolderName + "/"
	sa := option.WithCredentialsFile(s.FirebaseCredentials)
	bucketName := s.StorageBucketName
	ctx := context.Background()
	client, err := storage.NewClient(ctx, sa)
	if err != nil {
		return "", err
	}
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, file := filepath.Split(fileName)
	obj := client.Bucket(bucketName).Object(s.StorageFolderName + "/" + file)
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return baseURL + file, nil
}
