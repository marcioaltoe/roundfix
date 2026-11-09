package daemon

// Suite: Dirty Discard record persistence.
// Invariant: invalid dispositions never write and valid ones name the discarded set.
// Boundary IN: temporary artifact and worktree paths.
// Boundary OUT: Git removal and the Run Database.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestWriteDirtyDiscardRecordNamesTheDiscardedPaths(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	disposition := BranchDisposition{
		RunID: "ended-run", Worktree: filepath.Join(root, "worktree"),
		DiscardedPaths: []string{"unrelated.txt"}, DiscardDigest: "confirmed-digest",
	}
	path := filepath.Join(root, "artifacts", "branch-disposition.json")
	if err := WriteDirtyDiscardRecord(path, disposition, time.Now()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record branchDispositionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != BranchDispositionRecordSchemaVersion || record.Superseded || !slices.Equal(record.DiscardedPaths, disposition.DiscardedPaths) || record.DiscardDigest != disposition.DiscardDigest {
		t.Fatalf("record = %+v", record)
	}
	for _, name := range []string{"superseded", "no paths", "no digest", "inside worktree", "relative path"} {
		t.Run(name, func(t *testing.T) {
			invalid := disposition
			recordPath := filepath.Join(t.TempDir(), "refused.json")
			switch name {
			case "superseded":
				invalid.Superseded = true
			case "no paths":
				invalid.DiscardedPaths = nil
			case "no digest":
				invalid.DiscardDigest = ""
			case "inside worktree":
				recordPath = filepath.Join(disposition.Worktree, "record.json")
			case "relative path":
				recordPath = "relative-refused.json"
			}
			if err := WriteDirtyDiscardRecord(recordPath, invalid, time.Now()); err == nil {
				t.Fatal("accepted invalid record")
			}
			if _, err := os.Stat(recordPath); !os.IsNotExist(err) {
				t.Fatalf("record path exists: %v", err)
			}
		})
	}
}
