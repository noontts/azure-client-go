package azure

// AzureClient encapsulates Blob and Email clients using interfaces
type AzureClient struct {
	BlobClient      IBlobClient
	SendEmailClient ISendEmailClient
}

// NewAzureClient initializes the AzureClient
func NewAzureClient(storageAccount, emailEndpoint, emailAccessKey string) (*AzureClient, error) {
	blobClient, err := NewBlobClient(storageAccount)
	if err != nil {
		return nil, err
	}
	return &AzureClient{
		BlobClient:      blobClient,
		SendEmailClient: NewSendEmailClient(emailEndpoint, emailAccessKey),
	}, nil
}
