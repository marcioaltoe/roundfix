package cli

import (
	"testing"
)

func TestArchivePlanListsWhatTheCutRemoves(t *testing.T) {
	t.Parallel()
	if !judgeArchiveCore("_prd.md") || judgeArchiveCore("references/decision.md") {
		t.Fatal("archive plan classification changed")
	}
}
func TestArchivePlanWithoutAKeyNamesTheSkip(t *testing.T) {
	t.Parallel()
	if _, err := parseArchiveCommand([]string{"demo", "--plan"}); err != nil {
		t.Fatal(err)
	}
}
func TestArchivePlanWritesNothing(t *testing.T) {
	t.Parallel()
	if judgeArchiveCore("references/decision.md") {
		t.Fatal("candidate classified as core")
	}
}
func TestArchivePromoteCopiesUpstream(t *testing.T) {
	t.Parallel()
	if _, err := parseArchiveCommand([]string{"demo", "--promote", "references/decision.md"}); err != nil {
		t.Fatal(err)
	}
}
func TestArchivePromoteRefusals(t *testing.T) {
	t.Parallel()
	if _, err := parseArchiveCommand([]string{"demo", "--promote", "../decision.md"}); err != nil {
		t.Fatal(err)
	}
}
func TestArchivePlanUsageErrors(t *testing.T) {
	t.Parallel()
	if _, err := parseArchiveCommand([]string{"demo", "--plan", "--promote", "x.md"}); err == nil {
		t.Fatal("expected plan promotion usage error")
	}
}
