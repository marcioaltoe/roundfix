package agent

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Reaping owns only sessions recorded by earlier proofs. Its outcome is a
// warning, independent of the setup verdict of the proof about to start.
func (runner *ACPXRunner) reapUnclosedProofSessions(ctx context.Context, deadline time.Duration) {
	runner.proofReapOnce.Do(func() {
		home := strings.TrimSpace(environmentValue(runner.baseEnv(), "HOME"))
		if home == "" || ctx.Err() != nil {
			return
		}
		stateDir := filepath.Join(home, ".roundfix")
		// A file at the state root cannot contain any recorded sessions.
		// If this proof cannot record its own close, the writer reports that
		// failure with the recovery command rather than a second reap warning.
		if info, err := os.Stat(stateDir); err == nil && !info.IsDir() {
			return
		}
		dir := filepath.Join(stateDir, "unclosed-proof-sessions")
		entries, err := os.ReadDir(dir)
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		if err != nil {
			runner.warnUnreadableProofRecord(dir, err)
			return
		}
		type storedRecord struct {
			path   string
			record unclosedProofSessionRecord
		}
		var records []storedRecord
		for _, entry := range entries {
			if ctx.Err() != nil {
				return
			}
			if !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			data, err := os.ReadFile(path)
			var record unclosedProofSessionRecord
			if err == nil {
				err = json.Unmarshal(data, &record)
			}
			if err == nil && (record.Schema != "roundfix/unclosed-proof-session/v1" || record.Session == "" || record.WorkDir == "" || record.Runtime.ID == "" || record.RecordedAt.IsZero() || entry.Name() != record.Session+".json") {
				err = errors.New("invalid record metadata")
			}
			if err != nil {
				runner.warnUnreadableProofRecord(path, err)
				continue
			}
			records = append(records, storedRecord{path, record})
		}
		sort.Slice(records, func(i, j int) bool {
			a, b := records[i].record, records[j].record
			if a.RecordedAt.Equal(b.RecordedAt) {
				return a.Session < b.Session
			}
			return a.RecordedAt.Before(b.RecordedAt)
		})
		for _, stored := range records {
			if ctx.Err() != nil {
				return
			}
			record := stored.record
			runtime := RuntimeSpec{ID: record.Runtime.ID, Protocol: record.Runtime.Protocol, Command: record.Runtime.Command, OpenRouterKeyVariable: record.Runtime.OpenRouterKeyVariable}
			closeCtx, cancel := context.WithTimeout(ctx, deadline)
			err := runner.closeDisposableSession(closeCtx, runtime, record.Session, record.WorkDir)
			cancel()
			if err != nil {
				runner.warningf("recorded disposable Agent Session %q (%s) was not closed within %s: %v; the next profile proof retries it", record.Session, stored.path, deadline, err)
				return
			}
			if err := os.Remove(stored.path); err != nil && !errors.Is(err, os.ErrNotExist) {
				runner.warningf("remove closed Unclosed Proof Session record %s: %v", stored.path, err)
				return
			}
		}
	})
}

func (runner *ACPXRunner) warnUnreadableProofRecord(path string, err error) {
	runner.warningf("unreadable Unclosed Proof Session record %s: %v; remove it after closing its session by hand", path, err)
}
