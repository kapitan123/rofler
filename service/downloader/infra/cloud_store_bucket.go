package infra

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CloudStorageBucket struct {
	bucket       *storage.BucketHandle
	subdirectory string
}

func NewCloudStoreBucketClient(ctx context.Context, projectId string, videoFilesBucketUrl string) *CloudStorageBucket {
	newStorageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return &CloudStorageBucket{
		bucket:       newStorageClient.Bucket(videoFilesBucketUrl),
		subdirectory: "saved/",
	}
}

func (b *CloudStorageBucket) Save(ctx context.Context, fromReader io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	newFilePath := b.subdirectory + uuid.New().String() + ".mp4"
	writer := b.bucket.Object(newFilePath).NewWriter(ctx)

	defer func() {
		writer.Close()
		logrus.Infof("finish upload %s", newFilePath)
	}()

	logrus.Infof("start upload %s", newFilePath)

	_, err := io.Copy(writer, fromReader)
	if err != nil {
		logrus.Error(err)
	}

	if _, err := io.Copy(writer, fromReader); err != nil {
		return "", errors.Errorf("io.Copy: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", errors.Errorf("Writer.Close: %v", err)
	}

	return newFilePath, nil
}
