package infra

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type CloudStorageBucket struct {
	bucket *storage.BucketHandle
}

func NewCloudStoreBucketClient(ctx context.Context, projectId string, videoFilesBucketUrl string) *CloudStorageBucket {
	newStorageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return &CloudStorageBucket{
		bucket: newStorageClient.Bucket(videoFilesBucketUrl),
	}
}

func (b *CloudStorageBucket) Read(ctx context.Context, addr string, r io.Reader) error {
	// AK TODO implement
	return nil
}
