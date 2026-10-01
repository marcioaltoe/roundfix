package skills_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
	"roundfix/skills"
)

func TestTheConcreteContractGuideShipsAtTheHorizonPath(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	files, err := skills.Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	wantPath := filepath.ToSlash(strings.TrimPrefix(speccheck.ConcreteContractGuidePath, ".agents/skills/"))
	want, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(speccheck.ConcreteContractGuidePath)))
	if err != nil {
		t.Fatalf("read canonical guide: %v", err)
	}
	for _, file := range files {
		if file.Path == wantPath {
			if !bytes.Equal(file.Data, want) {
				t.Fatalf("embedded guide %q differs from canonical bytes", file.Path)
			}
			return
		}
	}
	t.Fatalf("embedded bundle has no guide at %q", wantPath)
}

func TestTheGuideExamplesAreReadByTheCheck(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(speccheck.ConcreteContractGuidePath)))
	if err != nil {
		t.Fatalf("read canonical guide: %v", err)
	}

	receipts := speccheck.Receipts(speccheck.ConcreteContractGuidePath, content)
	if len(receipts) != 1 {
		t.Fatalf("guide receipts = %d, want 1: %#v", len(receipts), receipts)
	}
	if _, proven, reason, err := speccheck.ProveReceipt(root, receipts[0]); err != nil {
		t.Fatalf("prove guide receipt: %v", err)
	} else if !proven {
		t.Fatalf("guide receipt was not proven: %s", reason)
	}

	transcripts := speccheck.SurfaceTranscripts(content)
	if len(transcripts) != 1 {
		t.Fatalf("guide transcripts = %d, want 1: %#v", len(transcripts), transcripts)
	}
	if transcript := transcripts[0]; transcript.Malformed != "" {
		t.Fatalf("guide transcript malformed: %#v", transcript)
	}
}
