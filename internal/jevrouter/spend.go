package jevrouter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"roundfix/internal/judge"
)

// Deps supplies an explicit environment, Home and key endpoint. A zero ceiling
// loads the judge's questions; production overrides also come from judge.Load.
type Deps struct {
	Env          []string
	HomeDir      string
	Client       *http.Client
	Endpoint     string
	Ceiling      float64
	MinCreditUSD float64
}

type Spend struct {
	TypeSafeLogged, OpenRouterLogged, KeyUsageMonthly, Total, Ceiling float64
	Key                                                               KeyStatus
}

func environmentKey(env []string) string {
	var key string
	for _, entry := range env {
		if value, ok := strings.CutPrefix(entry, "ROUNDFIX_OPENROUTER_API_KEY="); ok {
			key = value
		}
	}
	return key
}

// MonthSpend combines the UTC month's logged spend with the key's usage floor.
func MonthSpend(ctx context.Context, deps Deps, now time.Time) (Spend, error) {
	spend := Spend{Ceiling: deps.Ceiling}
	if spend.Ceiling == 0 {
		questions, err := judge.Load()
		if err != nil {
			return Spend{}, fmt.Errorf("load Jev ceiling: %w", err)
		}
		spend.Ceiling = questions.MonthlyCeilingUSD
	}
	if deps.HomeDir == "" {
		return Spend{}, errors.New("read Jev spend: Home is empty")
	}
	rows, err := judge.ReadMonth(ctx, deps.HomeDir, now)
	if err != nil {
		return Spend{}, fmt.Errorf("read Jev spend: %w", err)
	}
	for _, row := range rows {
		switch row.Transport {
		case "typesafe":
			spend.TypeSafeLogged += row.CostUSD
		case "openrouter":
			spend.OpenRouterLogged += row.CostUSD
		}
	}
	endpoint := deps.Endpoint
	if endpoint == "" {
		endpoint = "https://openrouter.ai/api/v1"
	}
	spend.Key, err = ReadKey(ctx, deps.Client, endpoint, environmentKey(deps.Env))
	if err != nil {
		return Spend{}, fmt.Errorf("read Jev spend: %w", err)
	}
	spend.KeyUsageMonthly = spend.Key.UsageMonthly
	spend.Total = spend.TypeSafeLogged + max(spend.OpenRouterLogged, spend.KeyUsageMonthly)
	return spend, nil
}

// CheckKeyLimit requires OpenRouter to cap a running prompt's monthly spend.
func (spend Spend) CheckKeyLimit() error {
	if spend.Key.LimitReset == nil || *spend.Key.LimitReset != "monthly" || spend.Key.Limit == nil || *spend.Key.Limit < 0 || *spend.Key.Limit > spend.Ceiling {
		return fmt.Errorf("jev_router_key_unbounded: set a monthly credit limit of at most US$%.4f on the key at OpenRouter", spend.Ceiling)
	}
	if spend.Key.LimitRemaining != nil && *spend.Key.LimitRemaining <= 0 {
		return errors.New("jev_ceiling_reached: OpenRouter key credit limit is exhausted")
	}
	return nil
}
