---
task: task_01
spec: 0338-the-release-decision-travels-with-the-spec
status: pending
type: backend
complexity: medium
---

# Task 01: The PRD's Delivery Type reaches the Archive Record and the queue's pull request title

## Overview

A Spec's PRD may declare its Delivery Type, `feat`, `fix`, `perf` or
`refactor`, in a frontmatter `type` key. The Spec Check reports any other
value, the archive copies the type into the Archive Record, and the Delivery
Queue titles the pull request `<type>: deliver …` from that record, `feat`
when it has none, so a fix Spec squashes as a fix and the release plan can
propose a patch. This answers the adopted Backlog Entry of 2026-11-01, "The
delivery pull request type follows the Spec"
([references](references/2026-11-01-the-delivery-pull-request-type-follows-the-spec.md)),
the Argus report of 2026-10-31.

## Requirements

1. MUST create `internal/spec/delivery_type.go` with `DeliveryType`, its four
   constants, `DeliveryTypes`, `ParsePRDDeliveryType` and
   `UnsupportedDeliveryTypeError` exactly as `_techspec.md` → Interfaces and
   API Contract 1 state, the error message being
   `PRD type "<value>" is not one of feat, fix, perf, refactor`.
2. MUST add `Type` to `ArchiveRecord` in `internal/spec/archive_record.go`:
   `BuildArchiveRecord` parses the PRD's type with `ParsePRDDeliveryType` and
   refuses on its error, `archiveRecordMapping` writes `type` after
   `delivery_commit` only when `Type` is non-empty, and `ParseArchiveRecord`
   reads `type` when present (API Contract 3). A PRD without `type` must give
   the same record bytes as before.
3. MUST create `internal/speccheck/delivery_type.go` with
   `detectPRDDeliveryType`, code constant `CodePRDType = "SC-PRD-TYPE"`, and
   call it from `Check` in `internal/speccheck/constraints.go` right after
   `detectClarificationRecord`, changing no other line there (API Contract 2).
4. MUST change `pullRequestTitle` in `internal/cli/deliver_workflow.go` to
   `pullRequestTitle(specSlug string, deliveryType spec.DeliveryType, specTitle func() (string, error)) (string, error)`,
   building every prefix from `<type>: deliver ` with `feat` for an empty
   type and keeping the 91-unit budget, cut rules and lazy title read; and
   make `Publication` read the Archive Record once, pass `Record.Type`, and
   return `plan publication: read Spec type: <err>` when the record cannot be
   read (API Contract 4). The archive and record commit subjects keep their
   bytes.
5. MUST create `internal/spec/delivery_type_test.go` with
   `TestPRDDeliveryTypeReadsTheFrontmatter` and
   `TestArchiveRecordCarriesThePRDDeliveryType`, and
   `internal/speccheck/delivery_type_test.go` with
   `TestPRDTypeFindingNamesTheValue`, as `_techspec.md` → Testing Approach
   lists them; every new top-level test and subtest calls `t.Parallel()`
   first.
6. MUST, in `internal/cli/deliver_queue_subject_test.go`, pass `""` as the
   type in the existing title cases, add
   `TestQueuePullRequestTitleFollowsTheDeliveryType` (each type for a fitting
   slug and for the 100-unit slug `0296-` plus 95 `x`, within 91 units, and
   `releaseplan.ClassifyCommit` of the title plus ` (#575)` reporting that
   type), and in `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`
   set the record's `Type` to `fix`, expect
   `fix: deliver 0296: A known Archive Record title`, and change the
   missing-record prefix to `plan publication: read Spec type:` (a declared
   break, `_techspec.md` → Declared breaks). Keep every test name.
7. MUST document the rule:
   - `docs/user-guide/commands/deliver.md` and the canonical
     `.agents/skills/roundfix/references/deliver.md`: the paragraph on the
     pull request title says `<type>: deliver <slug>`, where `<type>` is the
     Archive Record's Delivery Type and `feat` when it has none, and contains
     the sentence "The pull request title starts with the Spec's Delivery
     Type".
   - `docs/user-guide/commands/spec.md` and
     `.agents/skills/roundfix/references/spec.md`: the code list gains
     `SC-PRD-TYPE` with its meaning.
   - `docs/user-guide/commands/archive.md` and
     `.agents/skills/roundfix/references/archive.md`: the Archive Record's
     optional `type` key, copied from the PRD.
8. MUST add **Delivery Type** to `CONTEXT.md` through `domain-modeling`: the
   Conventional Commits type a Spec delivers as (`feat`, `fix`, `perf` or
   `refactor`), declared in the PRD's `type`, kept in the Archive Record and
   used in the Delivery Queue's pull request title, `feat` when absent, with a
   link to ADR-0309.
9. MUST append to
   `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
   an amendment beginning "Amendment (2026-11-01): the title's type is the
   Spec's Delivery Type." stating that `feat` becomes the Archive Record's
   `type`, `feat` when absent, with the same budget; set its `updated_at` to
   `2026-11-01T00:00:00Z`. Change no other line.
10. MUST copy the canonical skill with `make skills-sync`, then raise both
    version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
    -record-skill-versions -count=1`. Never hand-edit a mirror or a version
    field.
11. MUST re-record the Behavior Surface Record as the LAST step, after the
    last edit of any guide, with `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
12. MUST NOT touch `internal/releaseplan`, `internal/authorization`, any
    Baseline module, `docs/agents/`, the `write-*` skills, any existing
    Archive Record, any `### QA settlement` section, the Makefile, CI
    workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: a long-slug title and an archived PRD without `type` give today's bytes.
- [ ] Add the parser, the record key, the finding and the title type.
- [ ] Add and update the tests.
- [ ] Update the guides, references, `CONTEXT.md` and ADR-0309; sync, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] A PRD with `type: fix` archives to a record with `type: fix` and the queue titles its pull request `fix: deliver …`; without `type` both keep their bytes.
- [ ] `SC-PRD-TYPE` reports an unsupported type and nothing else.
- [ ] The guides, references, `CONTEXT.md` and ADR-0309 state the rule, the mirrors equal their canonical copies, and the skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
- instruction: `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/spec/delivery_type.go`
- creates: `internal/spec/delivery_type_test.go`
- creates: `internal/speccheck/delivery_type.go`
- creates: `internal/speccheck/delivery_type_test.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_queue_subject_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`

## Consumers

- `spec.ParsePRDDeliveryType` and the Archive Record's `type` — consumers: `./internal/spec`, `./internal/speccheck`, `./internal/cli`
- `pullRequestTitle` and `Publication` — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestQueuePullRequestTitleFollowsTheDeliveryType|TestQueuePullRequestTitleKeepsTheSlugWhileItFits|TestQueueMessagesStayReadableBySquashReaders|TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestQueuePullRequestTitleFollowsTheDeliveryType ' || { echo 'missing pass: TestQueuePullRequestTitleFollowsTheDeliveryType' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestPRDDeliveryTypeReadsTheFrontmatter|TestArchiveRecordCarriesThePRDDeliveryType)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestArchiveRecordCarriesThePRDDeliveryType ' || { echo 'missing pass: TestArchiveRecordCarriesThePRDDeliveryType' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^TestPRDTypeFindingNamesTheValue$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestPRDTypeFindingNamesTheValue ' || { echo 'missing pass: TestPRDTypeFindingNamesTheValue' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/spec ./internal/speccheck` — expected: exit 0. Before this Task the three new tests do not exist, so the first pass line is missing; afterwards the title follows the Delivery Type, the record carries it, the finding reports an unsupported value, and the Spec and Spec Check suites pass.
- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'The pull request title starts with the Spec' || { printf 'missing title rule in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do grep -qF -- 'SC-PRD-TYPE' "$f" || { printf 'missing SC-PRD-TYPE in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Delivery Type**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md | grep -qF -- 'Amendment (2026-11-01): the title' || { echo 'missing ADR-0309 amendment' >&2; exit 1; }; cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides name only `feat: deliver`, so the first check fails; afterwards the guides, references, glossary and ADR-0309 carry the rule, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Sources 1; Goals; Core Features 1, 2 and 5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 8; Glossary
- `_techspec.md` → Current behavior 1-3; Interfaces; Data Models; API Contracts 1-4 and 12; Invariants 1-2; Declared breaks; Testing Approach; Build Order 1
