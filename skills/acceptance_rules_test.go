package skills

import (
	"io/fs"
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

func TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle(t *testing.T) {
	t.Parallel()

	type skillContract struct {
		path    string
		phrases []string
	}

	contracts := []skillContract{
		{
			path: "write-idea/SKILL.md",
			phrases: []string{
				"The bound on grouping is a criterion, not a count:",
				"one context, where the sources and the work change the same component or contract, share a root cause, or complete one another,",
				"and Tasks that each earn their verification cycle.",
				"joins a neighbouring Task of the same context,",
				"small source of the same context joins a Spec that is open or being planned rather than becoming a one-Task Spec.",
				"Four implementation Tasks plus the QA gate is a reference, not a cap: needing more is a signal to check whether the Spec holds more than one context.",
				"When it holds more than one, split the scope and give each source one owning Spec.",
			},
		},
		{
			path: "write-prd/SKILL.md",
			phrases: []string{
				"The bound on grouping is a criterion, not a count:",
				"one context, where the sources and the work change the same component or contract, share a root cause, or complete one another,",
				"and Tasks that each earn their verification cycle.",
				"joins a neighbouring Task of the same context,",
				"small source of the same context joins a Spec that is open or being planned rather than becoming a one-Task Spec.",
				"Four implementation Tasks plus the QA gate is a reference, not a cap: needing more is a signal to check whether the Spec holds more than one context.",
				"When it holds more than one, split the scope and give each source one owning Spec.",
			},
		},
		{
			path: "write-techspec/SKILL.md",
			phrases: []string{
				"The bound on grouping is a criterion, not a count:",
				"one context, where the sources and the work change the same component or contract, share a root cause, or complete one another,",
				"and Tasks that each earn their verification cycle.",
				"joins a neighbouring Task of the same context,",
				"small source of the same context joins a Spec that is open or being planned rather than becoming a one-Task Spec.",
				"Four implementation Tasks plus the QA gate is a reference, not a cap: needing more is a signal to check whether the Spec holds more than one context.",
				"When it holds more than one, split the scope and give each source one owning Spec.",
			},
		},
		{
			path: "write-tasks/SKILL.md",
			phrases: []string{
				"**Sized for one fresh session.**",
				"**Earns its verification cycle.** Every Task pays an Agent turn, its declared Verification and the repository Verification at settlement.",
				"joins a neighbouring Task of the same context instead of becoming its own Task.",
				"It joins only a Task whose slice it belongs to, so the vertical-slice rule holds, and Tasks that declare the same file still run in series under the Wave collision rule.",
				"a graph that needs more is a signal to check whether the Spec holds more than one context, never a refusal by count.",
			},
		},
		{
			path: "roundfix/references/spec.md",
			phrases: []string{
				"The grouping bound is a criterion, not a count: one context, and Tasks that each earn their verification cycle; four implementation Tasks plus its QA gate is a reference, not a cap.",
			},
		},
	}

	for name, source := range roundfixLayoutSources() {
		for _, contract := range contracts {
			contract := contract
			t.Run(name+"/"+contract.path, func(t *testing.T) {
				t.Parallel()
				data, err := fs.ReadFile(source, contract.path)
				if err != nil {
					t.Fatal(err)
				}
				text := normalizeAcceptanceText(string(data))
				for _, absent := range []string{
					"fits four implementation Tasks plus its QA gate",
					"The grouping bound is four implementation Tasks",
				} {
					if strings.Contains(text, normalizeAcceptanceText(absent)) {
						t.Errorf("contains absent phrase %q", absent)
					}
				}
				for _, phrase := range contract.phrases {
					want := normalizeAcceptanceText(phrase)
					if got := strings.Count(text, want); got != 1 {
						t.Errorf("phrase %q count = %d, want 1", phrase, got)
					}
				}
			})
		}
	}
}
