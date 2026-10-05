package jevrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const DefaultMinCreditUSD = 15.0

type Credits struct{ TotalCredits, TotalUsage float64 }

func (c Credits) Balance() float64 { return c.TotalCredits - c.TotalUsage }

func MinCredit(configuredUSD float64) float64 {
	if configuredUSD > 0 {
		return configuredUSD
	}
	return DefaultMinCreditUSD
}

func CreditLeft(c Credits, key KeyStatus) float64 {
	left := c.Balance()
	if key.LimitRemaining != nil {
		left = min(left, *key.LimitRemaining)
	}
	return left
}

// ReadCredits makes one bounded request for account credit.
func ReadCredits(ctx context.Context, client *http.Client, endpoint, key string) (Credits, error) {
	failure := func(err error) (Credits, error) {
		return Credits{}, fmt.Errorf("read account credits: %w", keyError{cause: err, key: key})
	}
	if key == "" {
		return failure(errors.New("ROUNDFIX_OPENROUTER_API_KEY is not set"))
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(endpoint, "/")+"/credits", nil)
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
			TotalCredits *float64 `json:"total_credits"`
			TotalUsage   *float64 `json:"total_usage"`
		} `json:"data"`
	}
	if len(body) > 1<<20 || json.Unmarshal(body, &answer) != nil {
		return failure(errors.New("invalid account credits body"))
	}
	for _, field := range []struct {
		name  string
		value *float64
	}{{"total_credits", answer.Data.TotalCredits}, {"total_usage", answer.Data.TotalUsage}} {
		if field.value == nil || *field.value < 0 || math.IsNaN(*field.value) || math.IsInf(*field.value, 0) {
			return failure(fmt.Errorf("invalid data.%s", field.name))
		}
	}
	return Credits{TotalCredits: *answer.Data.TotalCredits, TotalUsage: *answer.Data.TotalUsage}, nil
}
