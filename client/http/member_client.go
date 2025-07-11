package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type MemberClientInterface interface {
	DoSomething(ctx context.Context) error
}

type MemberClient struct {
	*HTTPClient
}

func NewMemberClient(baseURL string) *MemberClient {
	return &MemberClient{HTTPClient: NewHTTPClient(baseURL, slog.Default(), false, false)}
}

func (c *MemberClient) DoSomething(ctx context.Context) error {
	url := c.BaseURL + "/member/endpoint"
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
		return fmt.Errorf("member request failed: %s", resp.Status)
	}
	return nil
}
