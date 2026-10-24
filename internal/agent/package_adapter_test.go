package agent

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestPackageAdaptersAreCopiesOfTheTestBinary(t *testing.T) {
	t.Parallel()
	if len(packageAdapterLinks) != 3 {
		t.Fatalf("package adapter count = %d, want 3", len(packageAdapterLinks))
	}
	binaryInfo, err := os.Stat(os.Args[0])
	if err != nil {
		t.Fatalf("stat compiled test binary: %v", err)
	}
	for packageName, target := range packageAdapterLinks {
		t.Run(packageName, func(t *testing.T) {
			targetInfo, err := os.Stat(target)
			if err != nil {
				t.Fatalf("stat package adapter: %v", err)
			}
			if os.SameFile(binaryInfo, targetInfo) {
				t.Error("package adapter shares the compiled test binary's file")
			}
			if targetInfo.Size() != binaryInfo.Size() {
				t.Errorf("package adapter size = %d, want %d", targetInfo.Size(), binaryInfo.Size())
			}
			if runtime.GOOS != "windows" && targetInfo.Mode()&0o111 == 0 {
				t.Errorf("package adapter mode = %v, want an executable bit", targetInfo.Mode())
			}
		})
	}
}

func TestASpawnedFixtureChildLeavesNoPackageAdapter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestSpawnedFixtureDeathFailsFast$")
	cmd.Env = environmentForTest(
		spawnedFixtureDeathTestEnv+"=fixture",
		"TMPDIR="+dir,
		"TMP="+dir,
		"TEMP="+dir,
	)
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 23 {
		t.Fatalf("fixture child error = %v, want exit status 23:\n%s", err, output)
	}
	if !strings.Contains("\n"+strings.ReplaceAll(string(output), "\r\n", "\n"), "\nfixture started\n") {
		t.Fatalf("fixture child did not report fixture started:\n%s", output)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read fixture temporary directory: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "roundfix-agent-package-adapters-") {
			t.Errorf("fixture child left package adapter tree %q", entry.Name())
		}
	}
}
