package baseline

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"strings"
	"testing"

	"roundfix/internal/spec"
)

// These literals lock the authored contract, independently of the module assets.
var retirementClauses = []struct{ id, guidance, module, rule, guide string }{
	{"clause.context.docs-one-job-per-directory", "Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers. Each of these directories holds live work only. The single history root `docs/history/` keeps only what a later reader needs, family by family (`docs/history/adr/`, `docs/history/backlog/`, `docs/history/findings/`, `docs/history/specs/`), and Git keeps every full text; `docs/history/handoffs/` and `docs/history/reviews/` hold only what an earlier layout retired there, which a history sanitize removes. Handoffs retire only on the user's explicit confirmation, and when confirmed, every handoff is deleted together and never enters history — the active `docs/handoffs/` directory is the capture door, never a shelf. Disposition is recorded in the front matter of the reduced entry for Findings and Backlog Entries, lifecycle front matter for ADRs, and the Archive Record or archive stamp for Specs. Rejected, deprecated and superseded ADRs retire whole to `docs/history/adr/`, while a `proposed` ADR stays in `docs/adr/`. A finished orphan Review Artifact retires by deletion and never enters history; a new Round never writes into history. Findings, Backlog Entries and Rollups with a terminal lifecycle status leave their active directory for the matching history family as reduced entries; preserve valid provenance and update dependent references as part of that transition. Spec-owned adopted references retain their Spec ownership and archive with the Spec; this rule clears active family directories without dismantling self-contained Specs. Spec-contained authorization records live at `<spec-root>/<slug>/_authorization.md`. The preserved legacy location `docs/workflow/authorizations/` remains readable for historical grants; its old paths are historical references and are not required to exist as worktree files. When the runtime's archive leaves an Archive Record, the record names each Spec-owned adopted reference, and the reference leaves the tree with the Spec folder. When the runtime offers a history sanitize, such as `roundfix history sanitize --apply --batch <n>`, it replaces each Spec folder an earlier archive left under the history root with its Archive Record, reduces each retired Finding and Backlog Entry to its front matter, title, first paragraph and the revision that holds its full text, and removes retired Review Artifacts and handoffs, one reviewed batch at a time after a tag marks the full history; retired ADRs stay whole.", "context-workflow", "rule.context.docs-layout", "docs-layout.md"},
	{"clause.context.backlog-01-operational-contract", "Name backlog entries `YYYY-MM-DD-<kebab-slug>.md`. Use this complete copyable Backlog Operational Contract:\n\n```markdown\n---\ntype: feat # feat | fix | perf | refactor\nstatus: open # open | promoted | declined | deferred | done | deprecated | superseded | closed | cancelled\ncreated: YYYY-MM-DD\nspec: null # Spec slug when status: promoted\nreason: null # required for terminal closure without a consuming Spec\n---\n```\n\nFor `type: feat`, use:\n\n```markdown\n# <Title — the intent in one line>\n\n## Opportunity\n\n<What could exist and for whom.>\n\n## Value\n\n<Why it would matter; the hypothesis.>\n\n## Shape\n\n<The rough form of a solution, explicitly non-binding.>\n```\n\nFor `type: fix`, use:\n\n```markdown\n# <Title — the defect in one line>\n\n## Symptom\n\n<What misbehaves, as a user or operator sees it.>\n\n## Where\n\n<The surface, command, or package, as known.>\n\n## Expected\n\n<The behavior that should replace it.>\n\n## Evidence\n\n<A finding link when one exists; `none yet` is honest.>\n```\n\nFor `type: perf`, use:\n\n```markdown\n# <Title — the cost in one line>\n\n## Slow\n\n<What is slow, for whom, and in which operation.>\n\n## Measured\n\n<The number that says so and how it was measured.>\n\n## Target\n\n<The number that would settle it.>\n```\n\nFor `type: refactor`, use:\n\n```markdown\n# <Title — the tangle in one line>\n\n## Tangled\n\n<What resists change, and where it is duplicated or coupled.>\n\n## Cost\n\n<What it makes slow, risky, or wrong to touch.>\n\n## Shape\n\n<The structure that would replace it, explicitly non-binding.>\n```\n\nKeep `open` entries in `docs/backlog/`. When a Spec adopts an entry, set `status: promoted` and `spec` to that Spec's slug, then move the entry to `docs/specs/<slug>/references/`; git history remains the discovery trail. Set `status: declined` or `status: deferred` only with a non-null `reason`. The current type set is open: a new type must be a Conventional Commits type that expresses intent. Adding a type is a contract change that requires a corpus re-record, never an informal addition. Use `refactor` as the canonical token, never an abbreviation.\n\n`docs/backlog/` holds unresolved `open` intent only. A `promoted` entry is adopted source material and moves to its consuming Spec's `references/` under the adoption contract; it is not a terminal implementation verdict. Every terminal Backlog Entry (`declined`, `deferred`, `done`, `deprecated`, `superseded`, `closed`, or `cancelled`) moves to `docs/history/backlog/` as a reduced entry in the same operation that records its terminal status. Its front matter carries the true disposition, any consuming Spec and the reason for closure, and the reduced entry keeps that front matter, the title, the first paragraph and the commit that holds the original intent, in the form the Findings archive rule gives. No absorption license is required for a Backlog Entry. An unknown status is a declaration error, never an inferred terminal disposition.", "context-workflow", "rule.context.docs-layout", "docs-layout.md"},
	{"clause.context.findings-09-archive", "The active `docs/findings/` directory holds unresolved Findings and unresolved Rollups only. Every terminal Finding or Rollup still in an active family directory (`done`, `deferred`, `deprecated`, `superseded`, `closed`, or `cancelled`) moves to `docs/history/findings/` in the same operation that records its terminal status; do not wait for the consuming Spec's implementation, QA, merge, or release. Record the true disposition in the front matter, and write the history copy as a reduced entry: the front matter, the title, the first paragraph, and a final line naming a commit that already holds the full text at its active path, normally the commit before the retirement:\n\n```markdown\nFull text in Git at `<40-hex commit>`: `<active path>`.\n```\n\nAn addendum written at retirement is committed in the active file first, so the named commit holds it; the original observation stays in Git and is never rewritten. An unknown status is a declaration error, never an inferred terminal disposition. A Spec-owned adopted reference retains its Spec ownership and travels with that Spec at archive; do not extract it from `references/` merely because its preserved source status is terminal. When absorption actually occurred, the archived Finding or Rollup retains a valid `absorbed_by:` pointer to the active or archived Spec that absorbs its content, or to an unresolved active Rollup:\n\n```yaml\nabsorbed_by: <active-rollup-basename-or-active-or-archived-spec-slug>\n```\n\nWhen any terminal Finding has no consuming Spec or Rollup, record a non-empty `closure_reason` and `closure_evidence` source locator instead of inventing an absorber. Its original observation and dated disposition stay in Git; absence of a Spec never prevents a legitimate terminal record from entering history. An invalid existing `absorbed_by` must be repaired, not hidden by these closure fields.\n\nBefore retiring a Rollup, transfer each member's absorption pointer to its true remaining owner and record the retiring Rollup's own absorption or evidenced terminal closure. Record the prior routing in a dated addendum committed before the Rollup is reduced, and validate all member and absorption references before completing the move. Referenced history is a migration obligation, not an exemption that keeps a terminal Rollup active indefinitely. Do not invent an owner or rewrite original evidence merely to satisfy the archive check.", "context-workflow", "rule.context.docs-layout", "docs-layout.md"},
}

func TestTheRetirementClausesCarryTheirText(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("context-workflow")
	if !ok {
		t.Fatal("missing context-workflow module")
	}
	for _, obsolete := range []string{
		"a byte-identical move for Review Artifacts and handoffs",
		"A finished orphan Review Artifact retires to `docs/history/reviews/`",
		"record the true disposition and its reason in a dated addendum",
	} {
		if strings.Contains(string(asset.Data), obsolete) {
			t.Errorf("obsolete retirement text: %s", obsolete)
		}
	}

	for _, want := range retirementClauses {
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Module(want.module)
			if !ok {
				t.Fatalf("missing module %s", want.module)
			}
			var module document
			if err := json.Unmarshal(asset.Data, &module); err != nil {
				t.Fatal(err)
			}
			found := false
			for _, rule := range objectsOrEmpty(module["rules"]) {
				if rule["id"] != want.rule {
					continue
				}
				for _, clause := range objectsOrEmpty(rule["clauses"]) {
					if clause["id"] != want.id {
						continue
					}
					found = true
					if clause["enforcement"] != "mandatory" || strings.Count(clause["guidance"].(string), want.guidance) != 1 {
						t.Errorf("clause %s force/text differs: %+v", want.id, clause)
					}
					if _, exists := clause["replaces"]; exists {
						t.Error("reworded clause has replaces")
					}
				}
			}
			if !found {
				t.Fatalf("missing clause %s in %s", want.id, want.rule)
			}
		})
	}
}

func TestTheRetirementClausesRenderInTheDocsLayoutGuide(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	for _, want := range retirementClauses {
		t.Run(want.id, func(t *testing.T) {
			asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + want.guide)
			if !ok {
				t.Fatalf("missing guide %s", want.guide)
			}
			local, err := os.ReadFile("../../docs/agents/" + want.guide)
			if err != nil {
				t.Fatal(err)
			}
			for name, data := range map[string][]byte{"Standard TypeScript golden": asset.Data, "repository guide": local} {
				t.Run(name, func(t *testing.T) {
					text := strings.Join(strings.Fields(string(data)), " ")
					literal := strings.Join(strings.Fields(want.guidance), " ")
					if count := strings.Count(text, literal); count != 1 {
						t.Fatalf("guide contains literal %d times, want exactly once", count)
					}
				})
			}
		})
	}
}

func TestAFindingInTheGuideReducedFormNeedsNoSanitize(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("context-workflow")
	if !ok {
		t.Fatal("missing context-workflow module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	var guidance string
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] == "clause.context.findings-09-archive" {
				guidance = clause["guidance"].(string)
			}
		}
	}
	_, fenced, ok := strings.Cut(guidance, "```markdown\n")
	if !ok {
		t.Fatal("Findings archive clause lacks fenced provenance form")
	}
	line, _, ok := strings.Cut(fenced, "\n```")
	if !ok || strings.Contains(line, "\n") {
		t.Fatal("provenance form must be one fenced line")
	}
	const revision = "0123456789abcdef0123456789abcdef01234567"
	line = strings.NewReplacer("<40-hex commit>", revision, "<active path>", "docs/findings/2026-10-08-observation.md").Replace(line)
	const prefix = "---\nstatus: closed\nclosure_reason: resolved\nclosure_evidence: docs/references/evidence.md\n---\n\n# Observation\n\nThe original observation.\n\n"
	for _, tc := range []struct {
		name, content string
		reduced       bool
	}{
		{"guide form", prefix + line + "\n", true},
		{"missing provenance", prefix, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, spec.ArchiveDir(spec.ArchiveKindFinding))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			const basename = "2026-10-08-observation.md"
			if err := os.WriteFile(filepath.Join(dir, basename), []byte(tc.content), 0o644); err != nil {
				t.Fatal(err)
			}
			if got := spec.IsReducedHistoryEntry([]byte(tc.content)); got != tc.reduced {
				t.Fatalf("IsReducedHistoryEntry = %v, want %v; line %q", got, tc.reduced, line)
			}
			plan, err := spec.PlanHistoryKind(root, revision, spec.ArchiveKindFinding)
			if err != nil {
				t.Fatal(err)
			}
			if tc.reduced {
				if len(plan.Files) != 0 {
					t.Fatalf("guide form needs sanitize: %+v", plan)
				}
			} else {
				want := spec.ArchiveDir(spec.ArchiveKindFinding) + "/" + basename
				if len(plan.Files) != 1 || plan.Files[0] != want || plan.Action != "reduce" {
					t.Fatalf("missing provenance must plan reduction of %s: %+v", want, plan)
				}
			}
		})
	}
}

func TestAnAdopterRetainsTheRetirementClauses(t *testing.T) {
	request, catalog := newClauseReplacementAdopter(t)
	outcome, err := buildPlanWithCatalog(context.Background(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil || outcome.Result.State != "ready" {
		t.Fatalf("refresh is not ready: %+v", outcome.Result)
	}
	delta := outcome.Plan.ClauseDelta
	if delta == nil {
		t.Fatal("refresh has no clause accounting")
	}
	for id, disposition := range delta.Dispositions {
		if disposition == ClauseUnaccounted {
			t.Errorf("unaccounted clause %s", id)
		}
	}
	for _, want := range retirementClauses {
		t.Run(want.id, func(t *testing.T) {
			if delta.Dispositions[want.id] != ClauseRetained {
				t.Errorf("disposition = %s, want retained", delta.Dispositions[want.id])
			}
			found := false
			for _, evidence := range outcome.Plan.Retention {
				if evidence.FromClause == want.id {
					found = true
					if evidence.Disposition != string(ClauseRetained) {
						t.Errorf("retention evidence = %+v", evidence)
					}
				}
			}
			if !found {
				t.Error("missing retention evidence")
			}
		})
	}
}
