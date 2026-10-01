package baseline

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAssetSyncFollowsASkillRenamedUpstream(t *testing.T) {
	t.Parallel()
	repo, root := newAssetsSyncTarget(t)
	source, _ := newAssetsSyncSource(t, root)
	renameAssetsSyncTargetSkill(t, root)
	renameAssetsSyncSourceSkill(t, source)
	payload, err := syncAssets(context.Background(), AssetsSyncRequest{SourceDir: source}, assetsSyncDependencies{repository: repo, assetRoot: root})
	if err != nil || !payload.OK || payload.Summary.Info == 0 {
		t.Fatalf("refresh = %+v, %v", payload, err)
	}
	_, err = LoadCatalog(os.DirFS(root))
	if err != nil {
		t.Fatalf("refreshed catalog: %v", err)
	}
	snapshots, err := loadAssetsSyncSnapshots(root)
	if err != nil {
		t.Fatal(err)
	}
	for id, snapshot := range snapshots {
		found := false
		for _, skill := range snapshot.Skills {
			if skill.Name == "handoff" {
				t.Fatalf("%s retains handoff", id)
			}
			if skill.Name == "handoff-next" {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s lacks handoff-next", id)
		}
	}
}

func TestAssetSyncCheckReportsDriftForARenameTheRefreshRepairs(t *testing.T) {
	t.Parallel()
	repo, root := newAssetsSyncTarget(t)
	source, _ := newAssetsSyncSource(t, root)
	renameAssetsSyncTargetSkill(t, root)
	renameAssetsSyncSourceSkill(t, source)
	before := captureAssetsSyncTree(t, root)
	payload, err := syncAssets(context.Background(), AssetsSyncRequest{SourceDir: source, Check: true}, assetsSyncDependencies{repository: repo, assetRoot: root})
	var syncErr *AssetsSyncError
	if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncExecution || payload.OK || len(payload.Findings) == 0 {
		t.Fatalf("check = %+v, %v", payload, err)
	}
	for _, finding := range payload.Findings {
		if finding.Code != "skills.setup-snapshot.drift" || finding.Severity != "error" || !strings.Contains(finding.Message, "differs") {
			t.Fatalf("drift finding = %+v", finding)
		}
	}
	if !reflect.DeepEqual(before, captureAssetsSyncTree(t, root)) {
		t.Fatal("check wrote assets")
	}
}

func TestAssetSyncRefusesARenameNoSetupProvides(t *testing.T) {
	t.Parallel()
	for _, check := range []bool{false, true} {
		name := "refresh"
		if check {
			name = "check"
		}
		t.Run(name, func(t *testing.T) {
			repo, root := newAssetsSyncTarget(t)
			source, _ := newAssetsSyncSource(t, root)
			renameAssetsSyncTargetSkill(t, root)
			assertAssetsSyncRenameRefusal(t, repo, root, source, check, "Generated setup snapshots are incompatible with the Baseline catalog", "Fix the canonical setup source before synchronizing snapshots.")
		})
	}
}

func TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift(t *testing.T) {
	t.Parallel()
	for _, check := range []bool{false, true} {
		name := "refresh"
		if check {
			name = "check"
		}
		t.Run(name, func(t *testing.T) {
			repo, root := newAssetsSyncTarget(t)
			source, _ := newAssetsSyncSource(t, root)
			if _, err := syncAssets(context.Background(), AssetsSyncRequest{SourceDir: source}, assetsSyncDependencies{repository: repo, assetRoot: root}); err != nil {
				t.Fatal(err)
			}
			renameAssetsSyncTargetSkill(t, root)
			assertAssetsSyncRenameRefusal(t, repo, root, source, check, "Go-owned canonical Baseline assets are invalid", "Fix the canonical Baseline catalog before synchronizing.")
		})
	}
}

func renameAssetsSyncTargetSkill(t *testing.T, root string) {
	t.Helper()
	filename := filepath.Join(root, "modules", "core.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"handoff"`) {
		t.Fatal("core lacks handoff")
	}
	if err := os.WriteFile(filename, []byte(strings.ReplaceAll(string(data), `"handoff"`, `"handoff-next"`)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func renameAssetsSyncSourceSkill(t *testing.T, source string) {
	t.Helper()
	checkout := filepath.Dir(source)
	if err := os.Rename(filepath.Join(checkout, "skills", "07-evidence-delivery", "handoff"), filepath.Join(checkout, "skills", "07-evidence-delivery", "handoff-next")); err != nil {
		t.Fatal(err)
	}
	lists, err := filepath.Glob(filepath.Join(source, "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range lists {
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, []byte(strings.ReplaceAll(string(data), "/handoff\n", "/handoff-next\n")), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runAssetsSyncGit(t, checkout, "add", ".")
	runAssetsSyncGit(t, checkout, "commit", "--quiet", "-m", "rename handoff upstream")
}

func assertAssetsSyncRenameRefusal(t *testing.T, repo, root, source string, check bool, message, action string) {
	t.Helper()
	before := captureAssetsSyncTree(t, root)
	payload, err := syncAssets(context.Background(), AssetsSyncRequest{SourceDir: source, Check: check}, assetsSyncDependencies{repository: repo, assetRoot: root})
	var syncErr *AssetsSyncError
	if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncInvalid || payload.OK || payload.Summary.Errors != 1 || len(payload.Findings) != 1 {
		t.Fatalf("refusal = %+v, %v", payload, err)
	}
	finding := payload.Findings[0]
	if finding.Code != "skills.setup-snapshot.drift" || finding.Severity != "error" || !strings.HasPrefix(finding.Message, message) || finding.Action != action {
		t.Fatalf("refusal finding = %+v", finding)
	}
	if !reflect.DeepEqual(before, captureAssetsSyncTree(t, root)) {
		t.Fatal("refusal wrote assets")
	}
}
