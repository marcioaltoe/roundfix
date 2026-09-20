package cli

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
)

// Suite: supersede command
// Invariant: supersede writes one trustworthy amendment or refuses without changing the Spec.
// Boundary IN: public CLI parsing, configured Spec roots, filesystem writes, streams, and exit codes.
// Boundary OUT: archive eligibility, which is owned by archive_test.go.
func TestSupersede(t *testing.T) {
	const (
		supersededSlug  = implementTestSlug
		supersedingSlug = "0002-delivered-widget"
		reason          = "Spec 0002 delivered the widget behavior and its acceptance evidence."
	)

	t.Run("accepted path writes only the amendment", func(t *testing.T) {
		homeDir, repoDir := newImplementWorkspace(t, nil)
		archivedSpecDir := archiveTestRepositoryPath(repoDir, spec.ArchiveKindSpec, supersedingSlug)
		mustMkdir(t, archivedSpecDir)
		mustWrite(t, filepath.Join(archivedSpecDir, "_prd.md"), "---\nstatus: archived\n---\n\n# Delivered widget\n")
		specDir := filepath.Join(repoDir, "docs", "specs", supersededSlug)
		before := snapshotDirectoryFiles(t, specDir)
		headBefore := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"supersede",
			"--spec", supersededSlug,
			"--by", supersedingSlug,
			"--reason", reason,
		}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("supersede exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
		}
		if stdout.String() != "superseded "+supersededSlug+" by "+supersedingSlug+"\n" {
			t.Fatalf("supersede stdout = %q", stdout.String())
		}
		if stderr.String() != "" {
			t.Fatalf("supersede stderr = %q, want empty", stderr.String())
		}

		record, err := spec.ReadSupersession(specDir)
		if err != nil {
			t.Fatalf("read supersession: %v", err)
		}
		if record.SupersededBy != supersedingSlug || record.Reason != reason || record.Explanation != reason {
			t.Fatalf("supersession = %+v", record)
		}
		if _, err := time.Parse(time.DateOnly, record.Date); err != nil {
			t.Fatalf("supersession date %q is not YYYY-MM-DD: %v", record.Date, err)
		}

		after := snapshotDirectoryFiles(t, specDir)
		delete(after, spec.SupersessionFilename)
		if !reflect.DeepEqual(after, before) {
			t.Fatalf("pre-existing Spec files changed\nbefore: %#v\nafter:  %#v", before, after)
		}
		headAfter := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
		if headAfter != headBefore {
			t.Fatalf("supersede changed HEAD from %s to %s", headBefore, headAfter)
		}
		assertNoRunDatabase(t, homeDir)
	})

	t.Run("refusals leave the Spec root byte-identical", func(t *testing.T) {
		tests := []struct {
			name      string
			specSlug  string
			bySlug    string
			prepare   func(t *testing.T, repoDir string)
			wantError string
		}{
			{
				name:      "unknown superseded Spec",
				specSlug:  "9998-unknown",
				bySlug:    supersedingSlug,
				wantError: "superseded Spec \"9998-unknown\" is unknown",
			},
			{
				name:      "unknown superseding Spec",
				specSlug:  supersededSlug,
				bySlug:    "9999-unknown",
				wantError: "superseding Spec \"9999-unknown\" is neither active nor archived",
			},
			{
				name:      "same Spec",
				specSlug:  supersededSlug,
				bySlug:    supersededSlug,
				wantError: "cannot supersede itself",
			},
			{
				name:     "already superseded",
				specSlug: supersededSlug,
				bySlug:   supersedingSlug,
				prepare: func(t *testing.T, repoDir string) {
					t.Helper()
					specDir := filepath.Join(repoDir, "docs", "specs", supersededSlug)
					mustWrite(t, filepath.Join(specDir, spec.SupersessionFilename), "existing amendment\n")
				},
				wantError: "already carries a supersession",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				homeDir, repoDir := newImplementWorkspace(t, nil)
				activeByDir := filepath.Join(repoDir, "docs", "specs", supersedingSlug)
				mustMkdir(t, activeByDir)
				mustWrite(t, filepath.Join(activeByDir, "_prd.md"), "---\nstatus: active\n---\n\n# Delivered widget\n")
				if tt.prepare != nil {
					tt.prepare(t, repoDir)
				}
				root := filepath.Join(repoDir, "docs", "specs")
				before := snapshotDirectoryFiles(t, root)
				var stdout bytes.Buffer
				var stderr bytes.Buffer

				code := runCLIContext(t, context.Background(), []string{
					"supersede",
					"--spec", tt.specSlug,
					"--by", tt.bySlug,
					"--reason", reason,
				}, &stdout, &stderr)

				if code != exitPreflight {
					t.Fatalf("supersede exit = %d, want %d; stderr=%q", code, exitPreflight, stderr.String())
				}
				if stdout.String() != "" {
					t.Fatalf("refusal stdout = %q, want empty", stdout.String())
				}
				if !strings.Contains(stderr.String(), tt.wantError) {
					t.Fatalf("refusal stderr = %q, want condition %q", stderr.String(), tt.wantError)
				}
				after := snapshotDirectoryFiles(t, root)
				if !reflect.DeepEqual(after, before) {
					t.Fatalf("refusal changed Spec root\nbefore: %#v\nafter:  %#v", before, after)
				}
				assertNoRunDatabase(t, homeDir)
			})
		}
	})

	t.Run("unknown flag refuses without writing", func(t *testing.T) {
		homeDir, repoDir := newImplementWorkspace(t, nil)
		root := filepath.Join(repoDir, "docs", "specs")
		before := snapshotDirectoryFiles(t, root)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{
			"supersede", "--spec", supersededSlug, "--by", supersedingSlug, "--reason", reason, "--unknown",
		}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("unknown flag exit = %d, want %d", code, exitPreflight)
		}
		if stdout.String() != "" || !strings.Contains(stderr.String(), "flag provided but not defined: -unknown") {
			t.Fatalf("unknown flag output: stdout=%q stderr=%q", stdout.String(), stderr.String())
		}
		if after := snapshotDirectoryFiles(t, root); !reflect.DeepEqual(after, before) {
			t.Fatalf("unknown flag changed Spec root\nbefore: %#v\nafter:  %#v", before, after)
		}
		assertNoRunDatabase(t, homeDir)
	})
}

func TestSupersedeRejectsANonActiveDeliverer(t *testing.T) {
	const (
		supersededSlug  = implementTestSlug
		supersedingSlug = "0002-delivered-widget"
		reason          = "Spec 0002 delivered the widget behavior and its acceptance evidence."
	)
	tests := []struct {
		name        string
		prd         string
		wantError   string
		omitPRD     bool
		omitSpecDir bool
	}{
		{
			name:      "draft",
			prd:       "---\nstatus: draft\n---\n\n# Draft deliverer\n",
			wantError: `frontmatter status is "draft"; expected "active"`,
		},
		{
			name:      "archived in active root",
			prd:       "---\nstatus: archived\n---\n\n# Misplaced archived deliverer\n",
			wantError: `frontmatter status is "archived"; expected "active"`,
		},
		{
			name:      "malformed frontmatter",
			prd:       "---\nstatus: [active\n---\n\n# Malformed deliverer\n",
			wantError: "malformed _prd.md",
		},
		{
			name:      "missing status",
			prd:       "---\nspec: 0002-delivered-widget\n---\n\n# Statusless deliverer\n",
			wantError: "frontmatter has no status",
		},
		{
			name:        "absent",
			wantError:   `superseding Spec "0002-delivered-widget" is neither active nor archived`,
			omitSpecDir: true,
		},
		{
			name:      "missing PRD",
			wantError: `superseding Spec "0002-delivered-widget" is neither active nor archived`,
			omitPRD:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, nil)
			root := filepath.Join(repoDir, "docs", "specs")
			if !tt.omitSpecDir {
				delivererDir := filepath.Join(root, supersedingSlug)
				mustMkdir(t, delivererDir)
				if !tt.omitPRD {
					mustWrite(t, filepath.Join(delivererDir, "_prd.md"), tt.prd)
				}
			}
			before := snapshotDirectoryFiles(t, root)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(t, context.Background(), []string{
				"supersede",
				"--spec", supersededSlug,
				"--by", supersedingSlug,
				"--reason", reason,
			}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("supersede exit = %d, want %d; stderr=%q", code, exitPreflight, stderr.String())
			}
			if stdout.String() != "" {
				t.Fatalf("refusal stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.wantError) {
				t.Fatalf("refusal stderr = %q, want condition %q", stderr.String(), tt.wantError)
			}
			if after := snapshotDirectoryFiles(t, root); !reflect.DeepEqual(after, before) {
				t.Fatalf("refusal changed Spec root\nbefore: %#v\nafter:  %#v", before, after)
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestSupersedeAcceptsAnActiveDeliverer(t *testing.T) {
	const (
		supersededSlug  = implementTestSlug
		supersedingSlug = "0002-delivered-widget"
	)
	homeDir, repoDir := newImplementWorkspace(t, nil)
	delivererDir := filepath.Join(repoDir, "docs", "specs", supersedingSlug)
	mustMkdir(t, delivererDir)
	mustWrite(t, filepath.Join(delivererDir, "_prd.md"), "---\nstatus: active\n---\n\n# Delivered widget\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"supersede",
		"--spec", supersededSlug,
		"--by", supersedingSlug,
		"--reason", "Spec 0002 delivered the widget behavior.",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("supersede exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("supersede stderr = %q, want empty", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func snapshotDirectoryFiles(t *testing.T, root string) map[string][]byte {
	t.Helper()
	snapshot := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relative)] = content
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return snapshot
}
