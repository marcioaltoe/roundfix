// Suite: Numbered File collisions.
// Invariant: only equal directory and number values collide, and the report names every owner and the next free number.
// Boundary IN: already-filtered repository-relative paths.
// Boundary OUT: configuration, Git and queue claim collection belong to the caller.
package delivery

import (
	"reflect"
	"slices"
	"testing"
)

func TestFindNumberedFileCollisions(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name         string
		defaultFiles []string
		subject      []string
		ahead        []NumberedFileClaims
		want         NumberedFileReport
	}{
		{
			name:    "ahead_item",
			subject: []string{"drizzle/0052_b.sql"},
			ahead:   []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}},
			want:    NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "drizzle/0052_b.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"}}, NextFree: []string{"drizzle/0053"}},
		},
		{
			name:         "default_branch",
			defaultFiles: []string{"drizzle/0051_x.sql", "drizzle/0052_main.sql"},
			subject:      []string{"drizzle/0052_b.sql"},
			want:         NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "drizzle/0052_b.sql", Other: "drizzle/0052_main.sql"}}, NextFree: []string{"drizzle/0053"}},
		},
		{
			name:    "same_path",
			subject: []string{"drizzle/0052_a.sql"},
			ahead:   []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}},
			want:    NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "drizzle/0052_a.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"}}, NextFree: []string{"drizzle/0053"}},
		},
		{
			name:    "by_value",
			subject: []string{"drizzle/052_b.sql"},
			ahead:   []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}},
			want:    NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "drizzle/052_b.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"}}, NextFree: []string{"drizzle/053"}},
		},
		{
			name:    "other_directory",
			subject: []string{"drizzle/meta/0052_snapshot.json"},
			ahead:   []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}},
		},
		{
			name:         "no_leading_digit",
			defaultFiles: []string{"drizzle/meta/_journal.json"},
			subject:      []string{"drizzle/meta/_journal.json"},
		},
		{
			name:         "next_free_above_every_claim",
			defaultFiles: []string{"drizzle/0051_x.sql"},
			subject:      []string{"drizzle/0052_b.sql"},
			ahead:        []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}, {Owner: "gamma", Paths: []string{"drizzle/0053_c.sql"}}},
			want:         NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "drizzle/0052_b.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"}}, NextFree: []string{"drizzle/0054"}},
		},
		{
			name:         "several",
			defaultFiles: []string{"drizzle/meta/0052_snapshot.json"},
			subject:      []string{"drizzle/meta/0052_snapshot.json", "drizzle/0052_b.sql"},
			ahead:        []NumberedFileClaims{{Owner: "alpha", Paths: []string{"drizzle/0052_a.sql"}}},
			want: NumberedFileReport{Collisions: []NumberedFileCollision{
				{Path: "drizzle/0052_b.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"},
				{Path: "drizzle/meta/0052_snapshot.json", Other: "drizzle/meta/0052_snapshot.json"},
			}, NextFree: []string{"drizzle/0053", "drizzle/meta/0053"}},
		},
		{name: "empty"},
		{
			name:         "root_and_digit_limit",
			defaultFiles: []string{"999999999999999999_a", "1000000000000000000_ignored", "0000000000000000000_ignored", "９_a"},
			subject:      []string{"999999999999999999_b", "1000000000000000000_ignored", "0000000000000000000_ignored", "９_a"},
			want:         NumberedFileReport{Collisions: []NumberedFileCollision{{Path: "999999999999999999_b", Other: "999999999999999999_a"}}, NextFree: []string{"1000000000000000000"}},
		},
		{
			name:         "sorted_distinct_subjects_and_rendered_owners",
			defaultFiles: []string{"d/9_z", "d/9_a"},
			subject:      []string{"d/009_b", "d/9_b", "d/009_b", "d/999_unclaimed"},
			ahead:        []NumberedFileClaims{{Owner: "zeta", Paths: []string{"d/9_a"}}, {Owner: "alpha", Paths: []string{"d/9_a"}}},
			want: NumberedFileReport{Collisions: []NumberedFileCollision{
				{Path: "d/009_b", Other: "d/9_a", Owner: "alpha"},
				{Path: "d/009_b", Other: "d/9_a"},
				{Path: "d/009_b", Other: "d/9_a", Owner: "zeta"},
				{Path: "d/009_b", Other: "d/9_z"},
				{Path: "d/9_b", Other: "d/9_a", Owner: "alpha"},
				{Path: "d/9_b", Other: "d/9_a"},
				{Path: "d/9_b", Other: "d/9_a", Owner: "zeta"},
				{Path: "d/9_b", Other: "d/9_z"},
			}, NextFree: []string{"d/010"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defaultBefore, subjectBefore := slices.Clone(tt.defaultFiles), slices.Clone(tt.subject)
			aheadBefore := slices.Clone(tt.ahead)
			for i := range aheadBefore {
				aheadBefore[i].Paths = slices.Clone(aheadBefore[i].Paths)
			}
			got := FindNumberedFileCollisions(tt.defaultFiles, tt.subject, tt.ahead)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("report = %#v, want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.defaultFiles, defaultBefore) || !reflect.DeepEqual(tt.subject, subjectBefore) || !reflect.DeepEqual(tt.ahead, aheadBefore) {
				t.Fatal("comparison modified its inputs")
			}
			if tt.name == "several" && got.Blocker() != "numbered-file-collision: drizzle/0052_b.sql collides with drizzle/0052_a.sql (alpha), drizzle/meta/0052_snapshot.json collides with drizzle/meta/0052_snapshot.json (default branch); next free number: drizzle/0053, drizzle/meta/0053" {
				t.Fatalf("several blocker = %q", got.Blocker())
			}
			if len(tt.want.Collisions) == 0 && (got.Detail() != "" || got.Blocker() != "") {
				t.Fatalf("empty report rendered detail=%q blocker=%q", got.Detail(), got.Blocker())
			}
		})
	}
}

func TestNumberedFileReportRendersTheBlocker(t *testing.T) {
	t.Parallel()
	report := NumberedFileReport{
		Collisions: []NumberedFileCollision{
			{Path: "drizzle/0052_b.sql", Other: "drizzle/0052_a.sql", Owner: "alpha"},
			{Path: "drizzle/meta/0052_snapshot.json", Other: "drizzle/meta/0052_snapshot.json"},
		},
		NextFree: []string{"drizzle/0053", "drizzle/meta/0053"},
	}
	const detail = "drizzle/0052_b.sql collides with drizzle/0052_a.sql (alpha), drizzle/meta/0052_snapshot.json collides with drizzle/meta/0052_snapshot.json (default branch); next free number: drizzle/0053, drizzle/meta/0053"
	if got := report.Detail(); got != detail {
		t.Fatalf("detail = %q, want %q", got, detail)
	}
	if got, want := report.Blocker(), "numbered-file-collision: "+detail; got != want {
		t.Fatalf("blocker = %q, want %q", got, want)
	}
}
