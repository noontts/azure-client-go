package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// LoggingRoundTripper logs HTTP requests and responses using slog
// Wraps another RoundTripper (usually http.DefaultTransport)
type LoggingRoundTripper struct {
	Proxied     http.RoundTripper
	Logger      *slog.Logger
	LogRequest  bool // flag to log requests
	LogResponse bool // flag to log responses
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log request method, url, headers
	if lrt.LogRequest {
		lrt.Logger.Info("HTTP request", "method", req.Method, "url", req.URL.String())
	}

	// Log request body (if present and small)
	if lrt.LogRequest && req.Body != nil && req.ContentLength != 0 {
		var bodyBytes []byte
		bodyBytes, _ = io.ReadAll(io.LimitReader(req.Body, 2048)) // limit to 2KB
		_ = req.Body.Close()
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) // restore body
		if len(bodyBytes) > 0 {
			lrt.Logger.Info("HTTP request body", "body", string(bodyBytes))
		}
	}

	resp, err := lrt.Proxied.RoundTrip(req)
	duration := time.Since(start)
	if err != nil {
		lrt.Logger.Error("HTTP error", "method", req.Method, "url", req.URL.String(), "err", err, "duration", duration)
		return nil, err
	}

	// Log response status, headers
	if lrt.LogResponse {
		lrt.Logger.Info("HTTP response", "method", req.Method, "url", req.URL.String(), "status", resp.StatusCode, "duration", duration)
	}

	// Log response body (if present and small)
	if lrt.LogResponse && resp.Body != nil {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048)) // limit to 2KB
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // restore body correctly
		if len(bodyBytes) > 0 {
			lrt.Logger.Info("HTTP response body", "body", string(bodyBytes))
		}
	}

	return resp, nil
}

// HTTPClient is a base client for making HTTP requests
type HTTPClient struct {
	Client  *http.Client
	BaseURL string
}

func NewHTTPClient(baseURL string, logger *slog.Logger, logRequest, logResponse bool) *HTTPClient {
	if logger == nil {
		logger = slog.Default()
	}
	return &HTTPClient{
		Client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &LoggingRoundTripper{Proxied: http.DefaultTransport, Logger: logger, LogRequest: logRequest, LogResponse: logResponse},
		},
		BaseURL: baseURL,
	}
}
