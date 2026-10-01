package cli

import (
	roundconfig "roundfix/internal/config"
	"testing"
)

func TestBuiltinReviewProfileStaysOnTheReviewProvidersRuntime(t *testing.T) {
	profile, err := roundconfig.ResolveProfile(roundconfig.Builtin(), roundconfig.CategoryReview, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateReviewProfileProvider("codex", profile.Profile); err != nil {
		t.Fatal(err)
	}
}
