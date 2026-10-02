package jevrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// keyError preserves the cause for matching without exposing credentials.
type keyError struct {
	cause error
	key   string
}

func (e keyError) Error() string {
	if e.key == "" {
		return e.cause.Error()
	}
	return strings.ReplaceAll(e.cause.Error(), e.key, "[redacted]")
}
func (e keyError) Unwrap() error { return e.cause }

// KeyUsage makes one bounded request for the key's monthly usage.
func KeyUsage(ctx context.Context, client *http.Client, endpoint, key string) (float64, error) {
	failure := func(err error) (float64, error) {
		return 0, fmt.Errorf("read key usage: %w", keyError{cause: err, key: key})
	}
	if key == "" {
		return failure(errors.New("ROUNDFIX_OPENROUTER_API_KEY is not set"))
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(endpoint, "/")+"/key", nil)
	if err != nil {
		return failure(err)
	}
	request.Header.Set("Authorization", "Bearer "+key)
	request.Header["User-Agent"] = nil // Suppress net/http's implicit header.
	httpClient := http.Client{}
	if client != nil {
		httpClient = *client
	}
	httpClient.Jar = nil
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	transport := httpClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	if standard, ok := transport.(*http.Transport); ok {
		copyTransport := standard.Clone()
		copyTransport.DisableCompression = true
		defer copyTransport.CloseIdleConnections()
		transport = copyTransport
	}
	httpClient.Transport = transport
	response, err := httpClient.Do(request)
	if err != nil {
		return failure(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return failure(fmt.Errorf("HTTP %d", response.StatusCode))
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	if err != nil {
		return failure(err)
	}
	var answer struct {
		Data struct {
			UsageMonthly *float64 `json:"usage_monthly"`
		} `json:"data"`
	}
	if len(body) > 1<<20 || json.Unmarshal(body, &answer) != nil || answer.Data.UsageMonthly == nil || *answer.Data.UsageMonthly < 0 {
		return failure(errors.New("invalid data.usage_monthly"))
	}
	return *answer.Data.UsageMonthly, nil
}
