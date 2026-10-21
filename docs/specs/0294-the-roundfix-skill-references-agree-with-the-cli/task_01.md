---
task: task_01
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
status: pending
type: docs
complexity: low
---

# Task 01: The review, events and deliver documents and the glossary state what the CLI does

## Overview

The Roundfix skill's `review` reference gives the reviewer's C3 and C5
Delivery Conventions in their folder-moving form, and its `events` reference
leaves `usage` out of the default replay and the `--filter` vocabulary. The
`deliver` reference and user guide name only an archived `_prd.md` where the
queue accepts an Archive Record, and the glossary's Run Event Stream entry
names four filters. This Task makes each statement match the code, as
`_techspec.md` → Current behavior records it, raises the Roundfix skill
version and re-records the Behavior Surface Record (`_prd.md` → Core Features
1-3). It answers the Inbox Entry of 2026-10-20 adopted at
[`references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md`](references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md).

## Requirements

1. MUST, in the canonical `.agents/skills/roundfix/references/review.md` and
   in `docs/user-guide/commands/review.md`, make the list that follows "The
   prompt carries Delivery Conventions version
   `roundfix/delivery-conventions/v3`." quote each convention verbatim from
   `deliveryConventions()` in `internal/cli/review_conventions.go`, as list
   items `- C1. <prompt>` through `- C5. <prompt>`, wrapped at the document's
   usual width but otherwise the exact characters of each prompt (backticks
   included, and no backticks added). Change the lead-in sentence "These
   describe what a delivery writes by design:" to "These describe what a
   delivery writes by design, and the prompt carries them verbatim:". At
   `8fb2da1f` this replaces C3 and C5 in the reference and C3 and C5 in the
   guide; C1, C2 and C4 already match.
2. MUST, in both documents of Requirement 1, replace the C5 eligibility
   sentences (the reference's "C5 applies to files under an archived Spec only
   when its `_prd.md` at the review head records `qa_override: true` and
   non-blank approval and reason. Missing or incomplete records and active
   Specs are ineligible for C5." and the guide's "C5 applies to an Archive
   Record or files under a legacy archived Spec only when the record or
   archived `_prd.md` at the review head records `qa_override: true` and
   non-blank approval and reason. Missing or incomplete records and active
   Specs are ineligible for C5.") with this paragraph:
   "C3 and C5 regions lie only under the archive root; an anchor under the
   active Spec root receives neither. C5 applies to an Archive Record,
   `<archive-root>/<slug>.md`, only when that record at the review head has
   disposition `qa-override` with `qa_override: true` and non-blank
   `qa_override_approval` and `qa_override_reason`. It applies to a file under
   a Legacy Archive Folder, `<archive-root>/<slug>/`, only when that folder's
   `_prd.md` front matter at the review head records `qa_override: true` and
   non-blank `qa_override_approval` and `qa_override_reason`; when
   `<archive-root>/<slug>.md` is an Archive Record, a file under
   `<archive-root>/<slug>/` is judged by that record. Missing or incomplete
   records and active Specs are ineligible for C5."
3. MUST, in the canonical `.agents/skills/roundfix/references/events.md`,
   replace "Default replay emits these public categories in journal cursor
   order: `task-status`, `batch`, `verification`, `outcome`, and
   `agent-selection`." with "Default replay emits these public categories in
   journal cursor order: `task-status`, `batch`, `verification`, `outcome`,
   `agent-selection`, and `usage`.", keep the following `--filter` sentence,
   and add these two rows after the `outcome` row of its stable-fields table:
   - "| `agent-selection` | `schema`, `run_id`, `category`, `time`, `cursor`,
     `scope_kind`, `scope_id`, `scope_identity`, `selection_role`, `status`,
     `runtime`, `model`, `reasoning_effort`, `summary`; optional `batch`,
     `work_item`, `work_category`, `profile_source`, `attempt`,
     `fallback_index`, `next_runtime`, `next_model`, `next_reasoning_effort`,
     `reason_code`, and `reason` |" (on one line);
   - "| `usage` | `schema`, `run_id`, `category`, `time`, `cursor`,
     `scope_kind`, `scope_id`, `runtime`, `model`, `reasoning_effort`,
     `token_basis`, `summary`; optional `batch`, `work_item`, `tokens`,
     `input_tokens`, `output_tokens`, `cached_read_tokens`,
     `cached_write_tokens`, `thought_tokens`, `cost_amount`, and
     `cost_currency` |" (on one line).
4. MUST, in the canonical `.agents/skills/roundfix/references/deliver.md`,
   replace "requires each prerequisite's archived `_prd.md` there." with
   "requires each prerequisite's Archive Record or, for a legacy archive, its
   archived `_prd.md` there.", and in `docs/user-guide/commands/deliver.md`
   replace "A prerequisite is met only when its archived `_prd.md` exists
   there." with "A prerequisite is met only when its Archive Record or, for a
   legacy archive, its archived `_prd.md` exists there.". In both, replace
   "reads the Spec's archived `_prd.md` on the local default branch and proves
   that the delivery commit that added it" with "reads the Spec's Archive
   Record, or a legacy archive's `_prd.md`, on the local default branch and
   proves that the delivery commit that added that file". Locate each sentence
   by its text: Spec 0291, already merged, changed other paragraphs of both
   files, and a later Spec may move them again. No other sentence of either
   file changes.
5. MUST revise **Run Event Stream** in `CONTEXT.md` through
   `domain-modeling`, replacing its sentence "Its stable Supervisor filters
   are task status, Batch boundary, Verification verdict, and terminal
   outcome." with "Its stable Supervisor categories are `task-status`,
   `batch`, `verification`, `outcome`, `agent-selection` and `usage`; replay
   includes every one by default, and a filter selects a subset of them." The
   entry keeps its other sentences and its `_Avoid_` line, and holds no other
   code span.
6. MUST run these commands, because this Task edits the Roundfix skill and
   two guides. Before them, raise both version fields of the Roundfix skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`) by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `review` and `deliver` guides' fingerprints change.
7. MUST search the Roundfix skill (`SKILL.md` and every reference), the
   `review`, `review-report-shape`, `events` and `deliver` user guides and the
   glossary again for a statement that describes the archive as moving the
   Spec folder, names an archived `_prd.md` as the only archive form, or lists
   a category set other than the six, and report each match it leaves with
   its reason.
8. MUST NOT change any Go source or test, any ADR, any text inside a `###
   QA settlement` section, `docs/references/skill-coverage.json`,
   `docs/user-guide/commands/events.md` or any file under `docs/history/`.

## Subtasks

- [ ] Quote C1–C5 verbatim and rewrite the C5 eligibility paragraph in the `review` reference and guide.
- [ ] Add `usage` to the `events` reference's category sentence and the two table rows.
- [ ] Correct the two `deliver` sentences in the reference and the guide, and revise **Run Event Stream**.
- [ ] Raise the skill version, sync the mirror, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The `review` reference, its mirror and the `review` guide carry C3 and C5 exactly as the reviewer's prompt does, and the new eligibility paragraph.
- [ ] The `events` reference and its mirror list six categories and have a stable-fields row for each.
- [ ] The `deliver` reference, its mirror and the `deliver` guide name the Archive Record in both sentences, and the glossary lists the six categories.
- [ ] The Roundfix skill version is raised and recorded, the mirror is in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `internal/cli/review_conventions.go`
- instruction: `internal/runevent/stream.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/events.md`
- interface: `skills/roundfix/references/events.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'C3. The archive commit adds an Archive Record at the archive root and removes the Spec folder, whose bytes stay in Git at source_revision; legacy archives move the folder and stamp its front matter.' || { printf 'missing C3 in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'For an Archive Record, these fields are in the record and the unchanged QA Task and Report stay in Git at source_revision.' || { printf 'missing C5 in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'is judged by that record. Missing or incomplete records and active Specs are ineligible for C5.' || { printf 'missing C5 eligibility in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/events.md skills/roundfix/references/events.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'in journal cursor order: .task-status., .batch., .verification., .outcome., .agent-selection., and .usage.\. ' || { printf 'missing category list in %s\n' "$f" >&2; exit 1; }; grep -q '^| .agent-selection. | ' "$f" && grep -q '^| .usage. | ' "$f" || { printf 'missing table rows in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Archive Record or, for a legacy archive, its archived' || { printf 'missing prerequisite sentence in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'proves that the delivery commit that added that file' || { printf 'missing merged fallback sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Event Stream**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'replay includes every one by default, and a filter selects a subset of them.' || { echo 'missing glossary revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the `review` reference gives the old C3, so the command fails at its first phrase check. Afterwards the `review`, `events` and `deliver` documents and the glossary carry the corrected statements, the mirror is in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map and the glossary contracts among them, pass.

## References

- `_prd.md` → Core Features 1-3; Goals; Success Metric 4; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-3; Build Order 1
