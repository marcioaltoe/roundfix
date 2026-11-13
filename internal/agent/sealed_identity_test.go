// Suite: sealed working directory identity.
// Invariant: only a regular repository marker may accompany a sealed prompt.
// Boundary IN: temporary directory entries, marker bytes, size and permissions.
// Boundary OUT: Agent sessions, ACP adapters and session tooling.
package agent

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"empty", "marker", "other file", "marker with other file", "marker directory", "marker symbolic link"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			workDir := newPrivateEmptyDirectory(t)
			marker := filepath.Join(workDir, SealedIdentityMarker)
			switch name {
			case "marker", "marker with other file":
				if err := os.WriteFile(marker, []byte("project = 'example'\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "marker directory":
				if err := os.Mkdir(marker, 0o700); err != nil {
					t.Fatal(err)
				}
			case "marker symbolic link":
				target := filepath.Join(t.TempDir(), "marker")
				if err := os.WriteFile(target, []byte("project = 'example'\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, marker); err != nil {
					t.Fatal(err)
				}
			}
			if name == "other file" || name == "marker with other file" {
				if err := os.WriteFile(filepath.Join(workDir, "AGENTS.md"), []byte("instructions"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			err := validateSealedPromptRequest(context.Background(), SealedPromptRequest{
				Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"},
				WorkDir: workDir,
				Input:   []byte(`{"sealed":true}`),
			})
			if name == "empty" || name == "marker" {
				if err != nil {
					t.Fatalf("valid sealed directory refused: %v", err)
				}
			} else if err == nil || err.Error() != "sealed Agent working directory must be empty apart from its .ai-memory.toml" {
				t.Fatalf("refusal = %v, want exact sealed directory refusal", err)
			}
		})
	}
}

func TestCopySealedIdentityMarkerCopiesOnlyARegularMarker(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"absent", "directory", "symbolic link", "oversized", "regular", "size limit", "write failure", "inspection failure"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := t.TempDir()
			workDir := newPrivateEmptyDirectory(t)
			marker := filepath.Join(repositoryRoot, SealedIdentityMarker)
			body := []byte("workspace = 'default'\nproject = 'example'\n")
			switch name {
			case "directory":
				if err := os.Mkdir(marker, 0o700); err != nil {
					t.Fatal(err)
				}
			case "symbolic link":
				target := filepath.Join(t.TempDir(), "marker")
				if err := os.WriteFile(target, body, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, marker); err != nil {
					t.Fatal(err)
				}
			case "oversized", "regular", "size limit", "write failure", "inspection failure":
				if name == "oversized" {
					body = bytes.Repeat([]byte("x"), sealedIdentityMarkerMaxBytes+1)
				} else if name == "size limit" {
					body = bytes.Repeat([]byte("x"), sealedIdentityMarkerMaxBytes)
				}
				if err := os.WriteFile(marker, body, 0o644); err != nil {
					t.Fatal(err)
				}
				if name == "write failure" {
					if err := os.Mkdir(filepath.Join(workDir, SealedIdentityMarker), 0o700); err != nil {
						t.Fatal(err)
					}
				} else if name == "inspection failure" {
					repositoryRoot = marker // A file cannot contain the source marker.
				}
			}
			err := CopySealedIdentityMarker(repositoryRoot, workDir)
			if name == "write failure" || name == "inspection failure" {
				if err == nil {
					t.Fatal("copy failure was ignored")
				}
				return
			}
			if err != nil {
				t.Fatalf("copy marker: %v", err)
			}
			entries, err := os.ReadDir(workDir)
			if err != nil {
				t.Fatal(err)
			}
			if name != "regular" && name != "size limit" {
				if len(entries) != 0 {
					t.Fatalf("ineligible marker copied: %v", entries)
				}
				return
			}
			if len(entries) != 1 || entries[0].Name() != SealedIdentityMarker {
				t.Fatalf("copied entries = %v", entries)
			}
			copyPath := filepath.Join(workDir, SealedIdentityMarker)
			copied, err := os.ReadFile(copyPath)
			if err != nil || !bytes.Equal(copied, body) {
				t.Fatalf("marker bytes differ: %v", err)
			}
			info, err := os.Lstat(copyPath)
			if err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
				t.Fatalf("copied marker is not a regular private file: info=%v err=%v", info, err)
			}
			if err := validateSealedPromptRequest(context.Background(), SealedPromptRequest{
				Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"},
				WorkDir: workDir,
				Input:   []byte(`{"sealed":true}`),
			}); err != nil {
				t.Fatalf("copied marker refused: %v", err)
			}
		})
	}
}

func TestCopySealedIdentityMarkerReadsOnlyTheMarkerItInspected(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"replaced by a symbolic link", "grown past the limit"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := t.TempDir()
			workDir := newPrivateEmptyDirectory(t)
			marker := filepath.Join(repositoryRoot, SealedIdentityMarker)
			if err := os.WriteFile(marker, []byte("project = 'example'\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			replace := func() {
				switch name {
				case "replaced by a symbolic link":
					target := filepath.Join(t.TempDir(), "outside")
					if err := os.WriteFile(target, []byte("outside = true\n"), 0o600); err != nil {
						t.Fatal(err)
					}
					if err := os.Remove(marker); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(target, marker); err != nil {
						t.Fatal(err)
					}
				case "grown past the limit":
					grown := bytes.Repeat([]byte("x"), sealedIdentityMarkerMaxBytes+1)
					if err := os.WriteFile(marker, grown, 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			if err := copySealedIdentityMarker(repositoryRoot, workDir, replace); err != nil {
				t.Fatalf("copy marker: %v", err)
			}
			entries, err := os.ReadDir(workDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Fatalf("replaced marker copied: %v", entries)
			}
		})
	}
}
