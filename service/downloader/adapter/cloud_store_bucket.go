package adapter

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type CloudStorageBucket struct {
	bucket       *storage.BucketHandle
	subdirectory string
}

func NewCloudStoreBucketClient(ctx context.Context, projectId string, videoFilesBucketUrl string) CloudStorageBucket {
	newStorageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return CloudStorageBucket{
		bucket:       newStorageClient.Bucket(videoFilesBucketUrl),
		subdirectory: "/downloaded",
	}
}

func (b *CloudStorageBucket) Save(id string, fromReader io.Reader) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	writer := b.bucket.Object(b.subdirectory + id).NewWriter(ctx)

	defer writer.Close()

	if _, err := io.Copy(writer, fromReader); err != nil {
		return errors.Wrap(err, "unable to copy data to bucket object writer")
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "unable to upload data to storage")
	}

	return nil
}

func (b *CloudStorageBucket) Read(ctx context.Context, addr string, r io.Reader) error {
	// AK TODO implement
	return nil
}
