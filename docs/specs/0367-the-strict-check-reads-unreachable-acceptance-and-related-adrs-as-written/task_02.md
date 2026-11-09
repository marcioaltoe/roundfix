---
task: task_02
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
status: completed
type: backend
complexity: medium
---

# Task 02: Only an ADR the Spec applies draws in the ADRs that cite it

## Overview

`SC-ADR-RELATED` seeds its candidates with every number in the PRD's Active
ADR obligations row, so an ADR a Spec lists only to record that it does not
apply draws in the ADRs that cite it, round after round. This Task answers
the Fiscus finding of 2026-11-07, "SC-ADR-RELATED cresce de forma transitiva
quando a Spec classifica ADRs como não aplicáveis"
([`references/2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md`](references/2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md),
reported by fiscus): as ADR-0377 decides, a listed ADR whose TechSpec
`## ADR Alternatives` line begins with `not applicable` stays listed but
seeds nothing. It also writes the amendment notices in ADR-0093 and ADR-0168
and the glossary sentence.

## Requirements

1. MUST add to `internal/speccheck/authoring_proof.go`, beside
   `adrAlternativePattern`,
   ``notApplicableRelationPattern = regexp.MustCompile(`(?i)^not applicable\b`)``,
   matched against the relation group of `adrAlternativePattern`.
2. MUST add to `internal/speccheck/citations.go`
   `appliedADRNumbers(listed map[string]bool, techSpecPath string) (map[string]bool, error)`:
   a copy of `listed` without each number whose line under the TechSpec's
   `## ADR Alternatives` (read with `markdownSectionLines` and
   `adrAlternativePattern`) has a relation matching
   `notApplicableRelationPattern`. A missing TechSpec returns every listed
   number; any other read error is returned wrapped with the path.
3. MUST make `detectADRConsistency` pass those seeds, not `listed`, to
   `firstListedCitation`, while it keeps skipping every number in `listed`,
   so a marked ADR is never itself reported. The summary, the two locations,
   the severity, the horizon and `SC-ADR-UNLISTED` stay as they are.
4. MUST change the `SC-ADR-RELATED` fix to exactly
   `Account for ADR-<N> in the Active ADR obligations row in <prd display path> by listing it or recording why it does not apply; when it does not apply, begin its ADR Alternatives line in the TechSpec with "not applicable", so the ADRs that cite only it are not drawn in.`
   Update every test that pins the old fix on the default branch when the
   Task runs, if any; at `2e59aebe` none does.
5. MUST create `internal/speccheck/adr_related_seed_test.go`, whose fixtures
   are written into `t.TempDir()` with accepted fixture decisions numbered
   0001 to 0004 (written ADR-<A> to ADR-<D> here): ADR-<A>, ADR-<B> citing
   ADR-<A> and ADR-<C> citing ADR-<B>, and a PRD whose Active ADR
   obligations row lists ADR-<A> and ADR-<B>, with:
   - `TestANotApplicableADRDrawsInNoCitingADR`: with the TechSpec line
     `- ADR-<B>: not applicable; …`, and again with `Not applicable — …`,
     `Check` reports no `SC-ADR-RELATED`;
   - `TestAnAppliedADRStillDrawsInItsCitingADR`: with
     `- ADR-<B>: applied; …`, and again with a relation that only contains
     the words later (`- ADR-<B>: kept as recorded, though its alternatives are not applicable here.`),
     `Check` reports exactly one `SC-ADR-RELATED`, for ADR-<C>, whose fix
     holds `begin its ADR Alternatives line in the TechSpec with "not applicable"`;
   - `TestWithoutATechSpecEveryListedADRSeeds`: without a TechSpec, ADR-<C>
     is reported;
   - `TestAnADRCitingAnAppliedADRIsStillRelated`: an ADR-<D> that cites
     ADR-<A> and ADR-<B> is reported while ADR-<B> is marked
     `not applicable`.

   Every new top-level test calls `t.Parallel()` as its first statement.
6. MUST, in `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`,
   set `updated_at: 2026-11-09T12:00:00Z` and add, as a new paragraph right
   after the paragraph that ends "That is a report an author can read, not
   noise.", this paragraph, exactly as the block holds it (one line in the
   file):

   ```markdown
   **Amended by [ADR-0377](0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md) (2026-11-09).** Only a listed ADR the Spec applies or changes is a seed of the candidates: an ADR whose line under the TechSpec's ADR Alternatives begins with `not applicable` stays listed but draws in no ADR that cites it. The closure stays one citation edge deep, and the check still reads only what the Spec wrote.
   ```

7. MUST, in `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`,
   set `updated_at: 2026-11-09T12:00:00Z` and add, as a new paragraph right
   after the first paragraph, which ends "a decision that did not yet exist
   when it was written.", this paragraph, exactly as the block holds it (one
   line in the file):

   ```markdown
   **Amended by [ADR-0377](0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md) (2026-11-09).** The listed ADR the citing ADR names must also be one the Spec applies or changes: a listed ADR whose ADR Alternatives line begins with `not applicable` opens no gap for the ADRs that cite it. The horizon below is unchanged.
   ```

   No other byte of either ADR changes. Neither ADR's Considered Options is
   taken by ADR-0377, so no option note is added.
8. MUST, in `CONTEXT.md`, append to the definition line of
   **Spec Consistency Check**, after "and a Spec that carries a Glossary
   Declaration is checked at any age.", a space and this sentence, exactly:

   ```markdown
   Only an ADR the Spec applies or changes draws in the ADRs that cite it: an ADR whose `## ADR Alternatives` line in the TechSpec begins with `not applicable` stays listed but opens no `SC-ADR-RELATED` gap for an ADR that cites only it ([ADR-0377](docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md)).
   ```

   Write it through the `domain-modeling` skill; no other line changes.
9. MUST NOT change ADR-0377 (already in the tree), any other ADR, any guide,
   skill or record; task_03 documents the rule in them.

## Subtasks

- [ ] Add the relation pattern and the seed function.
- [ ] Seed the related candidates with applied ADRs and change the fix.
- [ ] Add the test file.
- [ ] Write the two amendment notices and the glossary sentence.

## Acceptance Criteria

- [ ] An ADR marked `not applicable` draws in no ADR that cites only it, and stays listed.
- [ ] An applied ADR, an ADR whose relation only contains the words later, and a Spec without a TechSpec keep today's candidates.
- [ ] ADR-0093 and ADR-0168 name ADR-0377 as their amendment, and the glossary states the rule.

## Context

- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/authoring_proof.go`
- creates: `internal/speccheck/adr_related_seed_test.go`
- interface: `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`
- interface: `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md`

## Consumers

- The `SC-ADR-RELATED` finding of `speccheck.Check` — consumers: `./internal/speccheck`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestANotApplicableADRDrawsInNoCitingADR|TestAnAppliedADRStillDrawsInItsCitingADR|TestWithoutATechSpecEveryListedADRSeeds|TestAnADRCitingAnAppliedADRIsStillRelated)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestANotApplicableADRDrawsInNoCitingADR TestAnAppliedADRStillDrawsInItsCitingADR TestWithoutATechSpecEveryListedADRSeeds TestAnADRCitingAnAppliedADRIsStillRelated; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/speccheck` — expected: exit 0. Before this Task none of the four tests exists; afterwards they pass with the whole Spec Check suite, the horizon and depth-one tests included.
- `tr -s '[:space:]' ' ' < docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md | grep -qF -- 'Amended by [ADR-0377](0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md) (2026-11-09).' || exit 1; tr -s '[:space:]' ' ' < docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md | grep -qF -- 'updated_at: 2026-11-09T12:00:00Z' || exit 1; tr -s '[:space:]' ' ' < docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md | grep -qF -- 'opens no gap for the ADRs that cite it. The horizon below is unchanged.' || exit 1; tr -s '[:space:]' ' ' < docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md | grep -qF -- 'updated_at: 2026-11-09T12:00:00Z' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'Only an ADR the Spec applies or changes draws in the ADRs that cite it' || exit 1; grep -F -- '**Spec Consistency Check**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Consistency Check\*\*(:| —( |$))' || exit 1; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither ADR carries the notice, so the first check fails; afterwards both notices and the glossary sentence are in place, and the documentation contracts pass, the active Spec corpus counts included.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Core Feature 3; Success Metric 3
- `_techspec.md` → Interfaces; API Contract 2; Testing Approach; Build Order 2
- ADR-0377; ADR-0093; ADR-0168; ADR-0265

## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, domain-modeling, tech-writer, evidence-gate.

Implemented the assigned slice for Daemon Verification. The initial worktree
change was this Task's Daemon-owned `pending` to `in_progress` transition;
its status and checkboxes remain unchanged by this Agent.

- Acceptance criterion 1: `appliedADRNumbers` copies the listed set and removes
  only numbers whose ADR Alternatives relation begins with the case-insensitive
  `not applicable` token. Both marker variants in
  `TestANotApplicableADRDrawsInNoCitingADR` report no related, unlisted, or
  missing-alternatives findings. `TestAnADRCitingAnAppliedADRIsStillRelated`
  reports ADR-0004 through applied ADR-0001 even while ADR-0002 is marked.
- Acceptance criterion 2: applied relations, relations containing the words
  later, and a missing TechSpec still report ADR-0003 through ADR-0002.
  The new tests assert the exact fix, gap severity, summary seed, and two
  location paths. Existing horizon and depth-one tests passed in the focused
  selection. Other TechSpec read errors are returned with the path and `%w`;
  the original listed map and the unlisted-citation check are preserved.
- Acceptance criterion 3: ADR-0093 and ADR-0168 carry the exact ADR-0377
  amendments and required `updated_at` values. Through `domain-modeling`,
  the Spec Consistency Check definition gains exactly the required sentence.
  A read-only Python comparison against `git show HEAD:<path>`, using the
  literal paragraphs and sentence from this Task, confirmed that no other
  bytes of those three documents changed. ADR-0377, guides, skills, records,
  the manifest, and other Tasks were not edited.

Focused-check evidence:

- Before the production change,
  `rtk proxy go test ./internal/speccheck -run '^TestANotApplicableADRDrawsInNoCitingADR$' -count=1 -timeout=2m`
  exited 1: both marker variants incorrectly reported ADR-0003. The first
  sandboxed attempt could not access the Go build cache; the escalated run
  reproduced the behavioral failure without changing cache configuration.
- After the change,
  `rtk proxy go test ./internal/speccheck -run 'Test(ANotApplicableADR|AnAppliedADR|WithoutATechSpec|AnADRCitingAnAppliedADR|CheckADRClosureDepthOne|AnADRCommitted|TheHorizon|AnUncommitted|ASpecWithoutGit|ASpecInAnotherRepository|AShallowHistory)' -count=1 -timeout=3m -v`
  exited 0: all 17 selected top-level tests passed, including the four new
  regression tests and their subtests.
- `gofmt -w internal/speccheck/authoring_proof.go internal/speccheck/citations.go internal/speccheck/adr_related_seed_test.go`
  exited 0. `rtk proxy git diff --check` exited 0.
- The old fix text search under `internal`, agent guides, user guides, and
  both skill trees found only the changed implementation; no existing test
  pinned the old fix.

The declared Verification commands, full suites, and repository Verification
were not run; they remain Daemon-owned. No commit, push, or pull request was
created. No follow-up was identified within this slice.
