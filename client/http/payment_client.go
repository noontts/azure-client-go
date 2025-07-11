package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type PaymentClientInterface interface {
	DoSomething(ctx context.Context) error
}

type PaymentClient struct {
	*HTTPClient
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{HTTPClient: NewHTTPClient(baseURL, slog.Default(), false, false)}
}

func (c *PaymentClient) DoSomething(ctx context.Context) error {
	url := c.BaseURL + "/payment/endpoint"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment request failed: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // restore body for decoder

	// Unmarshal response into PaymentResponse struct
	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Printf("PaymentResponse: %+v\n", paymentResp)

	return nil
}
