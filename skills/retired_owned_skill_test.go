package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCouncilIsNotARoundfixOwnedSkill(t *testing.T) {
	names := Names()
	if len(names) != 12 {
		t.Fatalf("owned skill count = %d, want 12", len(names))
	}
	for _, name := range names {
		if name == "council" {
			t.Fatal("council remains a Roundfix-owned skill")
		}
	}
	if _, err := embedded.ReadDir("council"); !os.IsNotExist(err) {
		t.Fatalf("embedded council entry error = %v, want not-exist", err)
	}
	if _, err := os.Stat(filepath.Join("..", ".agents", "skills", "council")); !os.IsNotExist(err) {
		t.Fatalf("installed council skill error = %v, want not-exist", err)
	}
}

func TestImplementSpecIsNotARoundfixOwnedSkill(t *testing.T) {
	t.Parallel()
	for _, name := range Names() {
		if name == "implement-spec" {
			t.Fatal("implement-spec remains a Roundfix-owned skill")
		}
	}
	if _, err := embedded.ReadDir("implement-spec"); !os.IsNotExist(err) {
		t.Fatalf("embedded implement-spec entry error = %v, want not-exist", err)
	}
	if _, err := os.Stat(filepath.Join("..", ".agents", "skills", "implement-spec")); !os.IsNotExist(err) {
		t.Fatalf("installed implement-spec skill error = %v, want not-exist", err)
	}
}
