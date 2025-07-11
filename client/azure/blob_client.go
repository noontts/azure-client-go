package azure

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// IBlobClient defines the interface for blob operations
type IBlobClient interface {
	UploadBlob(ctx context.Context, container, blobName string, data []byte) error
	DownloadBlob(ctx context.Context, container, blobName string) ([]byte, error)
}

type BlobClient struct {
	Client *azblob.Client
}

func NewBlobClient(storageAccount string) (*BlobClient, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)
	azBlobClient, err := azblob.NewClient(blobURL, credential, nil)
	if err != nil {
		return nil, err
	}
	return &BlobClient{Client: azBlobClient}, nil
}

func (bc *BlobClient) UploadBlob(ctx context.Context, container, blobName string, data []byte) error {
	_, err := bc.Client.UploadBuffer(ctx, container, blobName, data, nil)
	return err
}

func (bc *BlobClient) DownloadBlob(ctx context.Context, container, blobName string) ([]byte, error) {
	resp, err := bc.Client.DownloadStream(ctx, container, blobName, nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	rr := resp.NewRetryReader(ctx, nil)
	defer rr.Close()
	_, err = buf.ReadFrom(rr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
