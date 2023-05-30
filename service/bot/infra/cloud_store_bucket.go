package infra

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
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

func (b *CloudStorageBucket) Read(ctx context.Context, addr string) (io.Reader, error) {
	reader, err := b.bucket.Object(addr).NewReader(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "unable to read file")
	}

	return reader, nil
}
