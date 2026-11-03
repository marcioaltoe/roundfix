---
task: task_01
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
status: completed
type: backend
complexity: low
---

# Task 01: The inbox source and the delivery plan skip hidden files and folders

## Overview

The one walk of `docs/_inbox/`, shared by `roundfix deliver plan` and the
follow-up's `repository-inbox` source, records every regular file, so a
folder marker such as `.gitkeep` and the contents of a hidden folder become
`note` items and inbox rows. After this Task the walk skips every file and
folder whose name starts with a dot. This answers the third request of the
adopted Backlog Entry of 2026-11-02, "The follow-up names what waits on
someone else and where each item went"
([references](references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md)),
reported by fluxus.

## Requirements

1. MUST, before changing code, characterize the starting tree: run the
   follow-up and `deliver plan` tests that exist today over a tree holding
   `docs/_inbox/.gitkeep` and record in the Result that both list it.
2. MUST change the `docs/_inbox/` walk in `deliveryPlanIntent`
   (`internal/cli/deliver_plan.go`) exactly as `_techspec.md` → API Contract 1
   states: an entry other than the walk root whose name starts with `.` is
   skipped, a folder with `fs.SkipDir`. The backlog and finding globs, the
   sort and every other line keep their bytes.
3. MUST add `TestDeliverPlanIgnoresHiddenInboxEntries` to
   `internal/cli/deliver_plan_test.go` and
   `TestFollowUpIgnoresHiddenInboxEntries` to `internal/cli/follow_up_test.go`,
   each over a repository holding `docs/_inbox/.gitkeep`,
   `docs/_inbox/.obsidian/app.json` and `docs/_inbox/note.md`: the plan prints
   exactly one inbox row, `docs/_inbox/note.md`, and the follow-up reports
   `repository-inbox` count 1 with exactly one `note` item, the note. Each new
   top-level test calls `t.Parallel()` first. Keep every existing test
   unchanged.
4. MUST document the rule with the phrase "whose name starts with a dot":
   - `docs/user-guide/commands/follow-up.md` and the canonical
     `.agents/skills/roundfix/references/follow-up.md`: `repository-inbox`
     reads every file in `docs/_inbox/` except a file or folder whose name
     starts with a dot, such as `.gitkeep`, and what such a folder holds;
   - `docs/user-guide/commands/deliver.md` and the canonical
     `.agents/skills/roundfix/references/deliver.md`: the same rule for inbox
     rows, as one sentence in the `deliver plan` paragraph that lists the
     intent rows, changing no other line.
5. MUST revise **Follow-up Item** in `CONTEXT.md` through `domain-modeling`:
   the note it lists is a file in `docs/_inbox/` outside any file or folder
   whose name starts with a dot.
6. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`. Never hand-edit a mirror or a version
   field.
7. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of any guide, with `go test -count=1 -tags docscontract
   ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
   -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
8. MUST NOT touch `internal/cli/follow_up.go`, `internal/cli/follow_up_inventory.go`,
   the `modes` reference, any Baseline asset, `docs/agents/`, any other skill,
   any `### QA settlement` section, the Makefile, CI workflows, `go.mod` or
   `go.sum`.

## Subtasks

- [ ] Characterize the listing of `.gitkeep` on the starting tree.
- [ ] Skip dot-named entries in the inbox walk.
- [ ] Add the two tests.
- [ ] Update the guides, references and `CONTEXT.md`; sync, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] `.gitkeep` and a file inside `.obsidian/` are listed neither by the follow-up nor by `deliver plan`, and a visible note still is.
- [ ] The follow-up and deliver guides and references and `CONTEXT.md` state the rule, the mirrors equal their canonical copies, and the skill version and Behavior Surface records pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/cli/deliver_plan.go`
- interface: `internal/cli/deliver_plan_test.go`
- interface: `internal/cli/follow_up_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `deliveryPlanIntent`'s inbox intents — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliverPlanIgnoresHiddenInboxEntries|TestFollowUpIgnoresHiddenInboxEntries|TestFollowUpListsRepositoryNotesAndMessageTitleFallback|TestFollowUpListsEveryPendingItemOnce)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverPlanIgnoresHiddenInboxEntries TestFollowUpIgnoresHiddenInboxEntries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards a hidden file or folder under the inbox is skipped by both readers and the existing note and inventory tests still pass.
- `for f in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md CONTEXT.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'whose name starts with a dot' || { printf 'missing hidden-entry rule in %s\n' "$f" >&2; exit 1; }; done; grep -F -- '**Follow-up Item**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Follow-up Item\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; cmp .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide, reference or glossary entry names the rule, so the first check fails; afterwards they do, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Sources 1; Goals; Core Features 1 and 7; Success Metric 1; Success Metric 8; Glossary
- `_techspec.md` → Current behavior 1; API Contract 1; API Contract 9; Invariant 5; Declared breaks; Testing Approach; Build Order 1


## Result

Skills loaded: implement-task, coding-guidelines, domain-modeling, golang-cli, golang-testing, testing-boss, no-workarounds, systematic-debugging, tech-writer, evidence-gate.

### Starting-tree characterization

Before changing repository code, ran the existing intent and repository-note tests
with temporary Go overlays that add `docs/_inbox/.gitkeep` to their fixtures:

`go test -overlay /private/tmp/roundfix-0343-characterization-bx043dvm/overlay.json -count=1 -timeout 5m -v ./internal/cli -run '^(TestDeliverPlanListsIntentThatIsNotApprovedToRun|TestFollowUpListsRepositoryNotesAndMessageTitleFallback)$'`

Both readers listed it. The plan printed `inbox\tdocs/_inbox/.gitkeep\t-`.
The follow-up listed `Kind:note Ref:docs/_inbox/.gitkeep` and its
`repository-inbox` count increased from 2 to 3; its existing count assertion
failed as expected (combined command exit 1). The plan test passed. The
overlays did not alter existing test files in the repository.

### Implementation and focused evidence

- Acceptance criterion 1: added only the dot-name guard to
  `deliveryPlanIntent`, excluding the walk root and returning `fs.SkipDir`
  for a hidden folder. The backlog/finding globs, sort and all other
  production bytes are unchanged. Added the two required public-CLI tests,
  each starting with `t.Parallel()` and holding `.gitkeep`,
  `.obsidian/app.json` and `note.md`. Existing tests remain byte-identical.
- Before the guard, `go test -count=1 -timeout 5m -v ./internal/cli -run
  'HiddenInboxEntries$'` exited 1: the plan listed all three files and
  follow-up counted 3 notes. After the guard,
  `go test -count=1 -timeout 5m -v ./internal/cli -run
  '^(TestDeliverPlanIgnoresHiddenInboxEntries|TestFollowUpIgnoresHiddenInboxEntries|TestDeliverPlanListsIntentThatIsNotApprovedToRun|TestFollowUpListsRepositoryNotesAndMessageTitleFallback|TestFollowUpListsEveryPendingItemOnce)$'`
  exited 0 with all five tests passing. The new tests require exactly one
  inbox row and exactly one note item, both `docs/_inbox/note.md`, and
  `repository-inbox` state `checked`, count 1.
- Acceptance criterion 2: the follow-up guides/references describe the file,
  folder and descendant exclusions; the deliver guide and canonical
  reference each gain exactly one sentence in their intent-row paragraph.
  Domain-modeling revised only **Follow-up Item** in `CONTEXT.md` to define
  the note as a file outside dot-named entries. A focused Python inspection
  confirmed the required phrase in all seven paths, both references and
  `SKILL.md` equal their canonical copies, and all existing test bytes are
  unchanged after removing the two added functions.
- `make skills-sync` exited 0. Then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`
  exited 0, raising both version fields in both copies from `0.1.172` to
  `0.1.173` and appending the generated digest record. No mirror or version
  field was edited by hand.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^TestGlossary'`
  exited 0. `git diff --check` reported no whitespace errors.
- As the final implementation regeneration, after every guide edit and the
  version recorder,
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 and regenerated `docs/references/behavior-surfaces.json`.

### Execution boundary

The starting worktree differed only in this Task's Daemon-owned
`status: in_progress`; it is preserved. The sandbox rejected its symlinked
writable root, so commands used the physical checkout with scoped escalation.
No declared Verification command, repository gate, full package suite, commit,
push or pull request was run. Task settlement and the declared Verification
remain Daemon-owned. No other Task, manifest, Baseline asset or prohibited
implementation path was changed. No follow-up work was needed for this slice.
