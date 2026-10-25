---
task: task_03
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: completed
type: backend
complexity: medium
---

# Task 03: The follow-up command reads no Secondbrain inbox and its keys are deprecated

## Overview

`roundfix follow-up` reads pending Inbox Entries from `secondbrain.root`. With
the inbox retired, this Task removes that source, takes the report's project
name from the `origin` remote or the checkout directory, and turns
`secondbrain.root` and `secondbrain.inbox_namespace` into deprecated keys
that load with one warning each, so the maintainer's User Config keeps
working. The command guide, the configuration guide, the Roundfix skill
reference and the glossary's **Follow-up Item** follow (`_prd.md` → Core
Feature 3; ADR-0330).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`](references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md)
   of 2026-10-24, its expected outcome 4.
2. MUST characterize before changing: run the follow-up and configuration
   tests on the starting tree and record which assert the Secondbrain
   source, the two keys or the configured namespace.
3. MUST implement `_techspec.md` → Data Model 7 and API Contracts 5 and 6 in
   `internal/config/config.go` and `internal/cli/follow_up_inventory.go`:
   matching Surface Transcript 1: remove the `Secondbrain` configuration type, field, overlay, validation,
   home expansion and Project Config warning, the two key constants and
   `followUpInbox`; add the two deprecated entries and the empty-mapping
   removal; keep `items` an empty array when nothing is found. No other
   source, field, exit code or output line changes.
4. MUST rewrite `internal/config/secondbrain_test.go` and update
   `internal/cli/follow_up_test.go` and `internal/cli/follow_up_status_test.go`
   as `_techspec.md` → Testing Approach (Build Order 3) describes, adding
   `TestTheRetiredSecondbrainKeysLoadWithAWarning`,
   `TestAnUnknownSecondbrainKeyStillFailsStrictValidation` and
   `TestFollowUpReadsNoSecondbrainInbox`. Every assertion an Inbox Entry fed
   keeps its strength with a message or Backlog Entry input; every new
   top-level test and subtest calls `t.Parallel()` first. No test starts the
   test binary as `roundfix`, and every exploratory `go test` carries
   `-timeout`.
5. MUST update `docs/user-guide/commands/follow-up.md` (the source list
   without `secondbrain-inbox`; the read-only paragraph says the command
   never reads or writes the Secondbrain and never calls ai-memory) and
   `docs/user-guide/configuration.md` (the `## Secondbrain inbox` section and
   the two table rows removed; both keys added to the deprecated-keys list
   with the warning of API Contract 6), naming no command, flag or key the
   binary lacks.
6. MUST update `.agents/skills/roundfix/references/follow-up.md`: the source
   list without `secondbrain-inbox`, and its last paragraph replaced by one
   saying the Secondbrain inbox source was retired, the two keys are
   deprecated and ignored with a warning, and the project name comes from
   the `origin` remote or the checkout directory. Raise both version fields
   of `.agents/skills/roundfix/SKILL.md` by one patch over the value the tree
   holds when this Task starts, then run, in order, `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which rewrites `skills/testdata/owned-skill-versions.json`. No mirror or
   record is hand-edited; never touch a `### QA settlement` section.
7. MUST change the **Follow-up Item** definition in `CONTEXT.md` to "One
   pending unit the Follow-up Command lists: an ai-memory message the caller
   hands in, an open Backlog Entry, an unresolved Finding, a note in
   `docs/_inbox/` or an active Spec, with its kind, reference, type, date,
   title and the project that sent it." and delete the Secondbrain Root
   entry, following the `domain-modeling` skill.
8. MUST, as its LAST step and after the last edit of any user guide,
   re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   which rewrites `docs/references/behavior-surfaces.json`; re-record again if
   a guide changes after it.
9. MUST NOT change any Baseline asset, any other skill, any other user guide
   page, the Makefile, `.roundfixrc.yml` or `go.mod`/`go.sum`.

## Subtasks

- [ ] Characterize the follow-up and configuration tests.
- [ ] Remove the source and the configuration type; add the deprecated keys.
- [ ] Update the tests and add the three new ones.
- [ ] Update the two guides, the Roundfix reference and its version, and the glossary; re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `roundfix follow-up` lists no `secondbrain-inbox` source and reads no Secondbrain directory.
- [ ] A User or Project Config naming either key loads, with one warning per key, and an unknown key under `secondbrain` still fails strict validation.
- [ ] The guides, the Roundfix reference and the glossary describe the command without the inbox.

## Consumers

- the follow-up report's sources and the configuration loader — consumers: `./internal/cli`, `./internal/config`, `./internal/docscontract`, `./skills`

## Context

- instruction: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- instruction: `internal/cli/follow_up.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `internal/config/config.go`
- interface: `internal/config/secondbrain_test.go`
- interface: `internal/cli/follow_up_inventory.go`
- interface: `internal/cli/follow_up_test.go`
- interface: `internal/cli/follow_up_status_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for name in TestTheRetiredSecondbrainKeysLoadWithAWarning TestAnUnknownSecondbrainKeyStillFailsStrictValidation; do go test -list "^${name}\$" ./internal/config | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -list '^TestFollowUpReadsNoSecondbrainInbox$' ./internal/cli | grep -qx TestFollowUpReadsNoSecondbrainInbox || { echo 'missing test TestFollowUpReadsNoSecondbrainInbox' >&2; exit 1; }; go test -count=1 ./internal/config && go test -count=1 -timeout 20m -run 'FollowUp' ./internal/cli` — expected: exit 0. Before this Task none of the three tests exists, so the command fails at its first check. Afterwards they exist, the configuration package passes with the deprecated keys, and every follow-up test passes with the source removed.
- `if grep -qF -- 'secondbrain-inbox' docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md internal/cli/follow_up_inventory.go; then echo 'the inbox source is still named' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'config: secondbrain.inbox_namespace is deprecated and ignored; the Secondbrain inbox source was retired, so remove it' || { echo 'configuration guide lacks the deprecation' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Follow-up Item**: One pending unit the Follow-up Command lists: an ai-memory message the caller hands in, an open Backlog Entry' || { echo 'glossary not updated' >&2; exit 1; }; if grep -qF -- '**Secondbrain Root**' CONTEXT.md; then echo 'retired glossary term kept' >&2; exit 1; fi; make skills-sync-check && go test -count=1 ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide names the `secondbrain-inbox` source, so the command fails at its first check. Afterwards no surface names the source, the configuration guide states the deprecation, the glossary carries the new **Follow-up Item** and no retired term, the mirrors and version record are current, and the documentation contracts pass with the re-recorded Behavior Surface Record.

## References

- `_prd.md` → Goal 5; Core Feature 3; User Experience; Success Metrics 4 and 5; User Story 4; Glossary
- `_techspec.md` → Current behavior 9-10; Interfaces; Data Model 7; API Contracts 5-6; Surface Transcript 1; Invariants 4-5; Testing Approach; Build Order 3
- ADR-0330

## Result

Implemented only this Task's slice for Backlog expected outcome 4. Task status
and checkboxes remain Daemon-owned; no declared Verification command, selected
repository Verification, commit, push or Pull Request was run.

### Starting-tree characterization

- `go test -timeout 10m ./internal/config -run 'Secondbrain|InboxNamespace' -count=1 -v`:
  passed. `TestSecondbrainRootIsReadFromUserConfig` pinned absolute/home-expanded
  roots; `TestProjectConfigCannotSetTheSecondbrainRoot` pinned the special
  Project Config warning; `TestSecondbrainKeysRefuseInvalidValues` pinned
  rejected relative roots and namespace values; and
  `TestTheInboxNamespaceMayComeFromProjectConfig` pinned namespace precedence.
  Those four tests were replaced with tests of the intentionally retired contract.
- `go test -timeout 15m ./internal/cli -run FollowUp -count=1`: passed before
  implementation. `followUpJSON` pinned the source list; inventory, four-answer
  ordering, fallback ordering and partial/secret judge tests consumed Inbox
  Entries. The empty report pinned the inbox source line, the missing-root test
  pinned its missing state, and the namespace test pinned the configured override.
  Status's no-inventory test used malformed inbox content.
- Initial working-tree inspection showed only this Task file changed; its
  existing `in_progress` status was retained. Shell execution required the
  resolved `/Volumes/Brain-SSD/cache/.roundfix/worktrees/` checkout because
  sandbox initialization rejected the symlinked configured writable root.

### Acceptance evidence

1. No Secondbrain source or directory reads: removed `followUpInbox` and its
   inventory source, and removed the namespace config override. Project names
   now come from `origin`, else the checkout directory. Inventory starts with
   a non-nil empty slice. `TestFollowUpReadsNoSecondbrainInbox` exercises the
   public CLI entry point with both retired keys, pending and malformed inbox
   files, and a namespace different from `origin`. It asserts exit 0, exact
   warning bytes, the exact remaining source list, the origin project name,
   `nothing_to_do`, `"items":[]`, and unchanged fixture bytes.
2. Backward-compatible config: removed the Secondbrain type, field, overlay,
   validation, expansion, constants and Project Config exception. Both keys
   use the existing deprecation registry. Only an empty Secondbrain mapping is
   removed before strict decoding. The two new config tests cover User and
   Project Config, formerly invalid values, ignored non-scalar values, an empty
   mapping, warning deduplication across both layers, and unknown-key refusal
   in either layer. New tests and subtests start with `t.Parallel()`.
3. Documentation: updated only the two declared guides and the canonical
   Roundfix follow-up reference; replaced Follow-up Item's definition and
   removed Secondbrain Root. The Roundfix version rose from 0.1.135 to
   0.1.136 in both fields. Mirrors and records were generated, not hand-edited;
   no QA settlement section changed. Message fixtures retain the former inbox
   titles, dates, types, bodies and ordering/judge assertions. Status now proves
   it ignores malformed repository notes, preserving its no-inventory boundary.

### Focused checks and generated records

- `go test -timeout 10m ./internal/config -run 'RetiredSecondbrain|UnknownSecondbrain' -count=1`:
  passed after the final config edits.
- `go test -timeout 15m ./internal/cli -run FollowUp -count=1`:
  passed after the final code edits, including list/status, ordering, secret
  handling, read-only and live-WAL checks. One intervening run printed PASS
  but was rejected by suiteguard because I overlapped it with the Behavior
  Surface recorder's repository write. After that recorder finished, the
  focused CLI run above passed with no concurrent writes.
- `make skills-sync`, then
  `go test -timeout 10m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  passed; only the declared Roundfix mirrors and its version record changed.
- `go test -timeout 15m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  passed after the last guide edit. Re-run as the last implementation step
  after recording this Result. It changes only the configuration and two
  guide fingerprints in the declared Behavior Surface Record.
- `git diff --check`: passed. Exact mirror comparison, retired-source absence
  inspection and changed-path allowlist inspection passed; no other Task,
  Task Graph, Baseline asset, unrelated skill/guide, Makefile, repository config
  or dependency file changed.

Declared Verification and terminal settlement remain for the Daemon. No
follow-up outside this Task's scope was implemented.

## Carry-forward provenance

- Source Run: `run_20261025T133232Z_0a2b2319e6c91868`
- Source commit: `8ca1eabdd1074d1a9a1c1bebc9a492b14f188a77`
