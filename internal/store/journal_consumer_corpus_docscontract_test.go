//go:build docscontract

//verify:relevant internal/store/testdata/ internal/runevent/ internal/tui/ internal/cli/attach.go internal/cli/events.go internal/cli/gc*.go internal/cli/reconcile.go

package store

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestJournalConsumerCorpusReplaysEveryConsumer(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("consumer corpus replay runs a nested go test; run without -short")
	}
	repositoryRoot := journalConsumerCorpusRepositoryRoot(t)
	testdataDir := filepath.Join(repositoryRoot, "internal", "store", "testdata")
	fixture := filepath.Join(testdataDir, "2026-08-11-prechange-roundfix.db")

	tests := []struct {
		name        string
		packagePath string
		virtualFile string
		harness     string
		testName    string
		environment []string
	}{
		{
			name:        "events Attach reconcile and gc preserve pre-change observations",
			packagePath: "./internal/cli",
			virtualFile: filepath.Join(repositoryRoot, "internal", "cli", "task10_consumer_observation_test.go"),
			harness:     filepath.Join(testdataDir, "task10-cli-consumer-harness_test.go.txt"),
			testName:    "TestTask10PrechangeJournalCLIConsumers",
			environment: []string{
				"TASK10_CORPUS_DB=" + fixture,
				"TASK10_EXPECTATIONS=" + filepath.Join(testdataDir, "2026-08-11-prechange-consumer-expectations.json"),
				"TASK10_RECORD=0",
			},
		},
		{
			name:        "Cockpit rendering preserves the pre-change frame",
			packagePath: "./internal/tui",
			virtualFile: filepath.Join(repositoryRoot, "internal", "tui", "task10_cockpit_observation_test.go"),
			harness:     filepath.Join(testdataDir, "task10-cockpit-consumer-harness_test.go.txt"),
			testName:    "TestTask10PrechangeJournalCockpitConsumer",
			environment: []string{
				"TASK10_CORPUS_DB=" + fixture,
				"TASK10_COCKPIT_EXPECTATION=" + filepath.Join(testdataDir, "2026-08-11-prechange-cockpit.golden"),
				"TASK10_RECORD=0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			overlay := struct {
				Replace map[string]string `json:"Replace"`
			}{Replace: map[string]string{test.virtualFile: test.harness}}
			overlayData, err := json.Marshal(overlay)
			if err != nil {
				t.Fatalf("encode %s overlay: %v", test.name, err)
			}
			overlayPath := filepath.Join(t.TempDir(), "overlay.json")
			if err := os.WriteFile(overlayPath, overlayData, 0o644); err != nil {
				t.Fatalf("write %s overlay: %v", test.name, err)
			}

			command := exec.CommandContext(
				t.Context(),
				"go", "test",
				"-overlay="+overlayPath,
				"-count=1",
				test.packagePath,
				"-run", "^"+test.testName+"$",
				"-v",
			)
			command.Dir = repositoryRoot
			command.Env = append(os.Environ(), test.environment...)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("replay %s: %v\n%s", test.name, err, output)
			}
			t.Logf("%s", output)
		})
	}
}
