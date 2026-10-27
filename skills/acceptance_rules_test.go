package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const networkProhibitedSettlementSection = "### Rows the Spec's authorization keeps offline\n\n" +
	"The qualifying declared `partial` above also accepts an outside-evidence row\n" +
	"the gate did not look up because the Spec's own authorization forbids that\n" +
	"lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)`\n" +
	"with the outside-evidence row named in its provenance. The subject names what\n" +
	"the authorization forbids, in its words: `network` for a network lookup, or\n" +
	"`another repository`; it is 1 to 80 characters and holds no parenthesis. Like\n" +
	"the network-denied row, it never decides a qualifying partial and needs no\n" +
	"Unreachable Acceptance declaration, a row counts toward one exempt kind at\n" +
	"most, and settlement, `roundfix qa-report accept` and `roundfix archive` read\n" +
	"it through the same policy (ADR-0276, ADR-0335)."

const networkProhibitedReferenceSentence = "So does an outside-evidence row the gate did not look up because the Spec's own authorization forbids that lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)` with the outside-evidence row named in its provenance, where the subject names what the authorization forbids in 1 to 80 characters without parentheses (ADR-0276, ADR-0335)."

func TestTheSettlementSkillsNameTheNetworkProhibitedRow(t *testing.T) {
	t.Parallel()

	for _, skill := range []string{"qa-gate", "archive-spec", "roundfix"} {
		path := filepath.Join("..", ".agents", "skills", skill, "SKILL.md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		section, ok := markdownSection(string(content), "### Rows the Spec's authorization keeps offline")
		if !ok {
			t.Errorf("%s is missing the Network-Prohibited Row section", path)
			continue
		}
		if normalizeAcceptanceText(section) != normalizeAcceptanceText(networkProhibitedSettlementSection) {
			t.Errorf("%s has unexpected Network-Prohibited Row section:\n%s", path, section)
		}
	}

	for _, reference := range []string{"archive", "settle"} {
		path := filepath.Join("..", ".agents", "skills", "roundfix", "references", reference+".md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if !strings.Contains(normalizeAcceptanceText(string(content)), normalizeAcceptanceText(networkProhibitedReferenceSentence)) {
			t.Errorf("%s is missing the Network-Prohibited Row reference sentence", path)
		}
	}
}

func TestWriteTasksComparesByteInvariantsWithTheDeliveryBase(t *testing.T) {
	t.Parallel()

	writeTasksPath := filepath.Join("..", ".agents", "skills", "write-tasks", "SKILL.md")
	content, err := os.ReadFile(writeTasksPath)
	if err != nil {
		t.Fatalf("read %s: %v", writeTasksPath, err)
	}
	wantSection := "## Byte invariants name the delivery base\n\n" +
		"A QA requirement that holds a file byte-identical compares it with the\n" +
		"Delivery Base, the merge base of the audited head with the default branch,\n" +
		"never with a named commit. A plan waits in the queue behind other Specs, and\n" +
		"each one that merges first may change the file, so the commit where the plan\n" +
		"measured stops being the right comparison. Keep that commit out of the\n" +
		"requirement; when it matters, record it in the TechSpec. `roundfix spec check`\n" +
		"reports `SC-QA-BYTES-PINNED` for a requirement of a pending QA Task that names\n" +
		"a byte invariant and a commit-shaped hexadecimal token together."
	section, ok := markdownSection(string(content), "## Byte invariants name the delivery base")
	if !ok {
		t.Fatalf("%s is missing the byte-invariant section", writeTasksPath)
	}
	if normalizeAcceptanceText(section) != normalizeAcceptanceText(wantSection) {
		t.Fatalf("%s has unexpected byte-invariant section:\n%s", writeTasksPath, section)
	}

	specReferencePath := filepath.Join("..", ".agents", "skills", "roundfix", "references", "spec.md")
	specReference, err := os.ReadFile(specReferencePath)
	if err != nil {
		t.Fatalf("read %s: %v", specReferencePath, err)
	}
	wantBullet := "`SC-QA-BYTES-PINNED` — a requirement of a pending QA Task names a byte invariant and a commit-shaped hexadecimal token together; compare with the Delivery Base instead. A gap that runs from the Tasks stage."
	if !strings.Contains(normalizeAcceptanceText(string(specReference)), normalizeAcceptanceText(wantBullet)) {
		t.Fatalf("%s is missing the SC-QA-BYTES-PINNED bullet", specReferencePath)
	}
}

func normalizeAcceptanceText(content string) string {
	return strings.Join(strings.Fields(content), " ")
}
