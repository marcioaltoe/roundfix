// Suite: repository skills lock version readiness
// Invariant: versions 1 and 2 read alike and unsupported versions retain an actionable typed refusal.
// Boundary IN: temporary repositories with installed skills and local lock bytes.
// Boundary OUT: remote skills installation and Baseline lock mutation.
package skills

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/skillhash"
)

func TestCheckRepositoryReadsASkillsLockOfVersionTwo(t *testing.T) {
	t.Parallel()
	root := writeReadyRepositoryFixture(t)
	before := snapshotRepositoryFixture(t, root)
	baseline, err := CheckRepository(t.Context(), root)
	if err != nil {
		t.Fatal(err)
	}
	lock := readRepositoryLockFixture(t, root)
	lock.Version = 2
	writeRepositoryLockFixture(t, root, lock)
	versionTwo := snapshotRepositoryFixture(t, root)
	got, err := CheckRepository(t.Context(), root)
	if err != nil || !got.Ready() || !reflect.DeepEqual(got, baseline) {
		t.Fatalf("version 2 readiness = %+v, %v; want %+v", got, err, baseline)
	}
	if reflect.DeepEqual(before, versionTwo) {
		t.Fatal("fixture version was not changed")
	}
	if after := snapshotRepositoryFixture(t, root); !reflect.DeepEqual(after, versionTwo) {
		t.Fatal("readiness changed the repository")
	}
}

func TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, lock, message string
		unsupported         bool
	}{
		{"version 3", `{"version":3,"skills":{}}`, "skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2", true},
		{"fractional version", `{"version":1.5,"skills":{}}`, "skills-lock.json version 1.5 is not supported; this Roundfix reads versions 1 and 2", true},
		{"text version", `{"version":"2","skills":{}}`, "skills-lock.json must hold a version number and a skills object", false},
		{"null version", `{"version":null,"skills":{}}`, "skills-lock.json must hold a version number and a skills object", false},
		{"boolean version", `{"version":true,"skills":{}}`, "skills-lock.json must hold a version number and a skills object", false},
		{"missing version", `{"skills":{}}`, "skills-lock.json must hold a version number and a skills object", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			filename := filepath.Join(root, "skills-lock.json")
			if err := os.WriteFile(filename, []byte(test.lock), 0o644); err != nil {
				t.Fatal(err)
			}
			before := snapshotRepositoryFixture(t, root)
			_, err := CheckRepository(t.Context(), root)
			want := fmt.Sprintf("validate skills lock version %q: %s", filename, test.message)
			var readinessErr *RepositoryReadinessError
			if !errors.As(err, &readinessErr) || readinessErr.Ownership != RepositoryOwnershipExternal || err.Error() != want || errors.Is(err, skillhash.ErrUnsupportedLockVersion) != test.unsupported {
				t.Fatalf("refusal = %v; want %q", err, want)
			}
			if after := snapshotRepositoryFixture(t, root); !reflect.DeepEqual(after, before) {
				t.Fatal("refusal changed repository")
			}
		})
	}
}
