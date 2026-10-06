package lighttier

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SpendLine records a light prompt's incremental cost, never a credential.
type SpendLine struct {
	Schema     string    `json:"schema"`
	Time       time.Time `json:"time"`
	Repository string    `json:"repository"`
	RunID      string    `json:"run_id"`
	Spec       string    `json:"spec"`
	Task       string    `json:"task"`
	Session    string    `json:"session"`
	Model      string    `json:"model"`
	CostUSD    float64   `json:"cost_usd"`
	CostSource string    `json:"cost_source"`
}

func monthPath(home string, now time.Time) string {
	return filepath.Join(home, ".roundfix", "openrouter", "implement", now.UTC().Format("2006-01")+".jsonl")
}

// AppendSpend appends one private JSON line to the UTC month's Light Spend Log.
func AppendSpend(home string, now time.Time, line SpendLine) (returnErr error) {
	line.Schema = "roundfix/light-spend/v1"
	line.Time = now.UTC()
	data, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("encode light spend: %w", err)
	}
	path := monthPath(home, now)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create light spend directory: %w", err)
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return fmt.Errorf("set light spend directory permissions: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open light spend log: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			returnErr = errors.Join(returnErr, fmt.Errorf("close light spend log: %w", err))
		}
	}()
	if err := file.Chmod(0600); err != nil {
		return fmt.Errorf("set light spend log permissions: %w", err)
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("append light spend: %w", err)
	}
	return nil
}

// ReadMonth sums the UTC month's cost; a missing log has no spend.
func ReadMonth(home string, now time.Time) (total float64, returnErr error) {
	file, err := os.Open(monthPath(home, now))
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("open light spend log: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			total = 0
			returnErr = errors.Join(returnErr, fmt.Errorf("close light spend log: %w", err))
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		var row struct {
			SpendLine
			Cost *float64 `json:"cost_usd"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return 0, fmt.Errorf("parse light spend log line %d: %w", lineNumber, err)
		}
		if row.Cost == nil || *row.Cost < 0 {
			return 0, fmt.Errorf("invalid light spend cost on line %d", lineNumber)
		}
		total += *row.Cost
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("read light spend log: %w", err)
	}
	return total, nil
}
