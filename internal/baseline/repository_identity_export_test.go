package baseline

import (
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"testing"
)

func TestInspectRepositoryIdentityIsTheSameInTwoClones(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	gittest.Run(t, repo, "init", "-b", "main")
	gittest.Run(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-m", "root")
	clone := filepath.Join(t.TempDir(), "clone")
	gittest.Run(t, repo, "clone", "--no-hardlinks", repo, clone)
	first, err := InspectRepositoryIdentity(t.Context(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := InspectRepositoryIdentity(t.Context(), clone, nil)
	if err != nil {
		t.Fatal(err)
	}
	inspection, err := InspectRepository(t.Context(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) || !reflect.DeepEqual(first, inspection.Identity) {
		t.Fatalf("identities differ: %+v %+v %+v", first, second, inspection.Identity)
	}
}
