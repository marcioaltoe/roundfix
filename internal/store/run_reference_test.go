// Suite: Run reference resolution.
// Invariant: a suffix prefix names exactly one Run or refuses every match.
// Boundary IN: disposable SQLite Run Database with deterministic Run ids.
// Boundary OUT: command rendering and Run mutation.
package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveRunReferenceNamesOneRunByItsShortID(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	home := t.TempDir()
	writer := openTestStore(t, ctx, home)
	ids := []string{
		"run_20260706T120000Z_93cee2a36ff33269",
		"run_20260706T120100Z_93cee2a3aaaaaaaa",
		"run_20260706T120200Z_0123456789abcdef",
	}
	raw, err := sql.Open("sqlite", "file:"+DatabasePath(home)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(0)")
	if err != nil {
		t.Fatal(err)
	}
	for i, id := range ids {
		req := sampleCreateRunRequest()
		req.GitRoot = t.TempDir()
		req.ArtifactDir = filepath.Join(req.GitRoot, ".roundfix")
		req.HeadRepository = "owner/repo" + string(rune('a'+i))
		run, err := writer.CreateRun(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		result, err := raw.ExecContext(ctx, `UPDATE runs SET id = ? WHERE id = ?`, id, run.ID)
		if err != nil {
			t.Fatal(err)
		}
		if n, err := result.RowsAffected(); err != nil || n != 1 {
			t.Fatalf("rename Run: affected=%d err=%v", n, err)
		}
	}
	if err := raw.Close(); err != nil {
		t.Fatal(err)
	}
	closeStore(t, writer)
	reader, err := OpenReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	defer closeStore(t, reader)

	for _, tc := range []struct{ name, reference, want string }{
		{"exact", ids[0], ids[0]},
		{"unique prefix", "0123456789abcdef", ids[2]},
		{"unique prefix", "01234567", ids[2]},
		{"longer prefix", "93cee2a36", ids[0]},
		{"trimmed", " 01234567 ", ids[2]},
		{"too short", "0123456", ""},
		{"not hex", "0123456z", ""},
		{"not hex", "0123456789ABCDEF", ""},
		{"too long", "0123456789abcdef0", ""},
		{"no match", "ffffffff", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run, found, err := reader.ResolveRunReference(ctx, tc.reference)
			if err != nil || found != (tc.want != "") || run.ID != tc.want {
				t.Fatalf("reference=%q id=%q found=%v err=%v; want id=%q", tc.reference, run.ID, found, err, tc.want)
			}
		})
	}
	t.Run("ambiguous", func(t *testing.T) {
		run, found, err := reader.ResolveRunReference(ctx, "93cee2a3")
		var ambiguous *AmbiguousRunReferenceError
		if found || run.ID != "" || !errors.As(err, &ambiguous) {
			t.Fatalf("id=%q found=%v err=%v; want ambiguous reference", run.ID, found, err)
		}
		want := `Run id prefix "93cee2a3" matches 2 Runs: run_20260706T120100Z_93cee2a3aaaaaaaa, run_20260706T120000Z_93cee2a36ff33269; pass more digits or the full Run id`
		if ambiguous.Reference != "93cee2a3" || !reflect.DeepEqual(ambiguous.RunIDs, []string{ids[1], ids[0]}) || err.Error() != want {
			t.Fatalf("ambiguity=%+v text=%q; want %q", ambiguous, err.Error(), want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		for _, reference := range []string{"", "  "} {
			run, found, err := reader.ResolveRunReference(ctx, reference)
			if found || run.ID != "" || err == nil || err.Error() != "Run ID is required" {
				t.Fatalf("reference=%q id=%q found=%v err=%v", reference, run.ID, found, err)
			}
		}
	})
	t.Run("Run stays exact", func(t *testing.T) {
		if run, found, err := reader.Run(ctx, "01234567"); found || err != nil || run.ID != "" {
			t.Fatalf("exact lookup resolved a prefix: run=%+v found=%v err=%v", run, found, err)
		}
	})
}
