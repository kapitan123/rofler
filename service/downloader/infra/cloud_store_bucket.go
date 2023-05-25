package infra

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		subdirectory: "/saved/",
	}
}

func (b *CloudStorageBucket) Save(ctx context.Context, fromReader io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	newFilePath := b.subdirectory + uuid.New().String() + ".mp4"
	writer := b.bucket.Object(newFilePath).NewWriter(ctx)

	defer writer.Close()

	if _, err := io.Copy(writer, fromReader); err != nil {
		return "", errors.Wrap(err, "unable to copy data to bucket object writer")
	}

	if err := writer.Close(); err != nil {
		return "", errors.Wrap(err, "unable to upload data to storage")
	}

	return newFilePath, nil
}

func (b *CloudStorageBucket) Read(ctx context.Context, id uuid.UUID, r io.Reader) error {
	// AK TODO implement
	return nil
}
