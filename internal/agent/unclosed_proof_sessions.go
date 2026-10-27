package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UnclosedProofSession identifies a disposable session whose close was unproven.
type UnclosedProofSession struct {
	Session  string
	WorkDir  string
	Deadline time.Duration
	Warning  string
}

type unclosedProofSessionRecord struct {
	Schema     string               `json:"schema"`
	Session    string               `json:"session"`
	WorkDir    string               `json:"work_dir"`
	Runtime    unclosedProofRuntime `json:"runtime"`
	RecordedAt time.Time            `json:"recorded_at"`
	Cause      string               `json:"cause"`
}

type unclosedProofRuntime struct {
	ID                    string `json:"id"`
	Protocol              string `json:"protocol"`
	Command               string `json:"command"`
	OpenRouterKeyVariable string `json:"openrouter_key_variable"`
}

func (runner *ACPXRunner) proofCloseTimeout(setupTimeout time.Duration) time.Duration {
	if runner.cleanupTimeout > 0 {
		return runner.cleanupTimeout
	}
	return setupTimeout
}

func (runner *ACPXRunner) recordUnclosedProofSession(runtime RuntimeSpec, session, workDir string, deadline time.Duration, cause error) *UnclosedProofSession {
	record := unclosedProofSessionRecord{
		Schema: "roundfix/unclosed-proof-session/v1", Session: session, WorkDir: workDir,
		Runtime:    unclosedProofRuntime{ID: runtime.ID, Protocol: runtime.Protocol, Command: runtime.Command, OpenRouterKeyVariable: runtime.OpenRouterKeyVariable},
		RecordedAt: eventClock(runner.Now)().UTC(), Cause: cause.Error(),
	}
	path, err := runner.writeUnclosedProofSession(record)
	warning := fmt.Sprintf("disposable Agent Session %q was not closed within %s after its profile proof: %v", session, deadline, cause)
	if err == nil {
		warning += fmt.Sprintf("; recorded in %s for the next profile proof to close", path)
	} else {
		args, _ := acpxCloseArgs(runtime, session, workDir)
		warning += fmt.Sprintf("; it could not be recorded (%v); close it with: acpx %s", err, strings.Join(args, " "))
	}
	runner.warningf("%s", warning)
	return &UnclosedProofSession{Session: session, WorkDir: workDir, Deadline: deadline, Warning: warning}
}

func (runner *ACPXRunner) writeUnclosedProofSession(record unclosedProofSessionRecord) (string, error) {
	home := strings.TrimSpace(environmentValue(runner.baseEnv(), "HOME"))
	if home == "" {
		return "", errors.New("HOME is not set")
	}
	dir := filepath.Join(home, ".roundfix", "unclosed-proof-sessions")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	file, err := os.CreateTemp(dir, ".unclosed-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	path := filepath.Join(dir, record.Session+".json")
	if err := os.Rename(file.Name(), path); err != nil {
		return "", err
	}
	return path, nil
}
