package azure

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// ISendEmailClient defines the interface for sending emails
type ISendEmailClient interface {
	SendEmail(ctx context.Context, req SendEmailRequest) error
}

type SendEmailClient struct {
	Endpoint  string
	AccessKey string
}

func NewSendEmailClient(endpoint, accessKey string) *SendEmailClient {
	return &SendEmailClient{
		Endpoint:  endpoint,
		AccessKey: accessKey,
	}
}

type EmailAttachment struct {
	Name            string `json:"name"`
	ContentType     string `json:"contentType"`
	ContentInBase64 string `json:"contentInBase64"`
}

type SendEmailPayload struct {
	SenderAddress string `json:"senderAddress"`
	Content       struct {
		Subject   string `json:"subject"`
		PlainText string `json:"plainText"`
		HTML      string `json:"html,omitempty"`
	} `json:"content"`
	Recipients struct {
		To []struct {
			Address string `json:"address"`
		} `json:"to"`
	} `json:"recipients"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
}

// SendEmailRequest encapsulates all parameters for sending an email
type SendEmailRequest struct {
	Sender      string
	Recipient   string
	Subject     string
	PlainText   string
	HTML        string
	Attachments []EmailAttachment
}

// SendEmail sends an email, optionally with attachments and HTML content
func (c *SendEmailClient) SendEmail(ctx context.Context, req SendEmailRequest) error {
	const maxAttachmentSize = 10 * 1024 * 1024 // 10MB
	for _, att := range req.Attachments {
		decoded, err := base64.StdEncoding.DecodeString(att.ContentInBase64)
		if err != nil {
			return errors.New("invalid base64 in attachment: " + att.Name)
		}
		if len(decoded) > maxAttachmentSize {
			return errors.New("attachment exceeds 10MB: " + att.Name)
		}
	}
	payload := SendEmailPayload{
		SenderAddress: req.Sender,
		Attachments:   req.Attachments,
	}
	payload.Content.Subject = req.Subject
	payload.Content.PlainText = req.PlainText
	payload.Content.HTML = req.HTML
	payload.Recipients.To = append(payload.Recipients.To, struct {
		Address string `json:"address"`
	}{Address: req.Recipient})

	url := c.Endpoint + "/emails:send?api-version=2023-03-31"
	b, _ := json.Marshal(payload)
	reqHttp, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", "Bearer "+c.AccessKey)
	reqHttp.Header.Set("x-ms-client-request-id", uuid.New().String())
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("email send failed: %s", resp.Status)
	}
	return nil
}
