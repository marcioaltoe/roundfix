//go:build linux

// Boundary: one executable copy in a temporary directory, held open for writing.
package cli

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestDeliveryRunRetriesAnItemBinaryStillOpenForWriting(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	source, err := os.ReadFile("/bin/true")
	if err != nil {
		t.Skipf("no /bin/true to copy: %v", err)
	}
	path := filepath.Join(t.TempDir(), "roundfix")
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(source); err != nil {
		t.Fatal(err)
	}
	if _, err := f.workflow.runRoundfix(commandContextForTest(t, t.Context()), path, f.repo); !errors.Is(err, syscall.ETXTBSY) {
		// The kernel refuses every start while the descriptor stays open, so
		// the bounded retries end with the same refusal.
		t.Fatalf("start while open for writing: err=%v, want ETXTBSY", err)
	}
	released := make(chan error, 1)
	go func() {
		time.Sleep(30 * time.Millisecond)
		released <- writer.Close()
	}()
	result, err := f.workflow.runRoundfix(commandContextForTest(t, t.Context()), path, f.repo)
	if closeErr := <-released; closeErr != nil {
		t.Fatal(closeErr)
	}
	if err != nil || result.exitCode != exitOK {
		t.Fatalf("start after the writer closed: result=%+v err=%v", result, err)
	}
}
