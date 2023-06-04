package infra

import (
	"context"
	"fmt"
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

func (b *CloudStorageBucket) Read(ctx context.Context, addr string, writer io.Writer) error {
	reader, err := b.bucket.Object(addr).NewReader(ctx)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to read file %s", addr))
	}

	_, err = io.Copy(writer, reader)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to copy file content %s", addr))
	}

	return nil
}
