package judge

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type judgeLog struct{ path string }

func newJudgeLog(home string, now time.Time) judgeLog {
	return judgeLog{filepath.Join(home, ".roundfix", "judge", now.UTC().Format("2006-01")+".jsonl")}
}
func (l judgeLog) monthCost() (float64, error) {
	f, err := os.Open(l.path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	var total float64
	for scanner.Scan() {
		var row struct {
			Cost *float64 `json:"cost_usd"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return 0, fmt.Errorf("parse judge log: %w", err)
		}
		if row.Cost == nil || *row.Cost < 0 {
			return 0, errors.New("invalid judge log cost")
		}
		total += *row.Cost
	}
	return total, scanner.Err()
}

type logLine struct {
	Schema         string             `json:"schema"`
	Time           time.Time          `json:"time"`
	Repository     string             `json:"repository"`
	Spec           string             `json:"spec"`
	Judgment       string             `json:"judgment"`
	Artifact       string             `json:"artifact"`
	Line           int                `json:"line"`
	Target         string             `json:"target"`
	StateHash      string             `json:"state_hash"`
	QuestionID     string             `json:"question_id"`
	Transport      string             `json:"transport"`
	ResponseID     string             `json:"response_id"`
	Provider       string             `json:"provider"`
	RequestedModel string             `json:"requested_model"`
	Model          string             `json:"model"`
	Answer         *string            `json:"answer"`
	Probabilities  map[string]float64 `json:"probabilities"`
	Confidence     *float64           `json:"confidence"`
	Noul           *float64           `json:"noul"`
	LatencyMS      int64              `json:"latency_ms"`
	InputTokens    int                `json:"input_tokens"`
	OutputTokens   int                `json:"output_tokens"`
	CostUSD        float64            `json:"cost_usd"`
	CostSource     string             `json:"cost_source"`
	Status         int                `json:"status"`
	Attempts       int                `json:"attempts"`
	Error          string             `json:"error"`
	Outcome        string             `json:"outcome"`
}

func (l judgeLog) append(row logLine) error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0700); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Dir(l.path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if err := f.Chmod(0600); err != nil {
		return errors.Join(err, f.Close())
	}
	err = json.NewEncoder(f).Encode(row)
	return errors.Join(err, f.Close())
}
