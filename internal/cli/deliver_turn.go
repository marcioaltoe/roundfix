package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

const deliveryTurnSchema = "roundfix/delivery-turn/v1"

type deliveryTurnRecord struct {
	Schema        string `json:"schema"`
	OwnerPID      int    `json:"owner_pid"`
	OwnerIdentity string `json:"owner_identity"`
	delivery.Turn
}

func turnRecordPath(artifactDir, gitRoot string, pid int) string {
	return filepath.Join(artifactDir, "delivery", deliveryOwnerID(gitRoot, pid), "turn.json")
}

type ownerTurnRecorder struct {
	loaded         roundconfig.Loaded
	once           sync.Once
	pid            int
	identity, path string
	err            error
}

func (recorder *ownerTurnRecorder) RecordTurn(ctx context.Context, turn delivery.Turn) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	recorder.once.Do(func() {
		recorder.pid = os.Getpid()
		recorder.identity, recorder.err = store.OwnerProcessIdentity(ctx, recorder.pid)
		if recorder.err != nil {
			return
		}
		var artifactDir string
		artifactDir, recorder.err = roundconfig.ResolveArtifactDirectory(recorder.loaded.Config.Defaults.ArtifactDir, recorder.loaded.GitRoot, recorder.loaded.HomeDir)
		if recorder.err == nil {
			recorder.path = turnRecordPath(artifactDir, recorder.loaded.GitRoot, recorder.pid)
		}
	})
	if recorder.err != nil {
		return recorder.err
	}
	dir := filepath.Dir(recorder.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create Turn Record directory: %w", err)
	}
	file, err := os.CreateTemp(dir, ".turn-*")
	if err != nil {
		return fmt.Errorf("create temporary Turn Record: %w", err)
	}
	defer func() { _ = os.Remove(file.Name()) }()
	record := deliveryTurnRecord{Schema: deliveryTurnSchema, OwnerPID: recorder.pid, OwnerIdentity: recorder.identity, Turn: turn}
	encodeErr := json.NewEncoder(file).Encode(record)
	if err := errors.Join(encodeErr, file.Close()); err != nil {
		return fmt.Errorf("write Turn Record: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.Rename(file.Name(), recorder.path); err != nil {
		return fmt.Errorf("replace Turn Record: %w", err)
	}
	return nil
}

func printDeliveryTurns(loaded roundconfig.Loaded, queue store.DeliveryQueue, stdout, stderr io.Writer) {
	artifactDir, err := roundconfig.ResolveArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
	path := turnRecordPath(artifactDir, loaded.GitRoot, queue.OwnerPID)
	var record deliveryTurnRecord
	if err == nil {
		var data []byte
		data, err = os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		if err == nil {
			err = json.Unmarshal(data, &record)
		}
		if err == nil && record.Schema != deliveryTurnSchema {
			err = fmt.Errorf("unsupported schema %q", record.Schema)
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: Delivery Queue turn record %s could not be read: %v\n", path, err)
		return
	}
	if record.OwnerPID != queue.OwnerPID || record.OwnerIdentity != queue.OwnerIdentity {
		return
	}
	for _, line := range delivery.TurnLines(queue, record.Turn) {
		fmt.Fprintln(stdout, line)
	}
}
