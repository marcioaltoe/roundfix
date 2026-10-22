---
task: task_01
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
status: completed
type: docs
complexity: low
---

# Task 01: The guides, the Roundfix reference and the glossary describe the Check Wait

## Overview

The configuration guide has no entry for the Delivery Queue's check wait, and
the `deliver` guide and the Roundfix skill's `deliver` reference say only that
"the existing checks timeout" parks the item as `checks-timeout`. This Task
documents the User Config key `delivery.check_timeout`, its default, bounds and
scope and the new park text in those documents, raises the Roundfix skill's
version, re-records the Behavior Surface Record and adds **Check Wait** to the
glossary (`_prd.md` → Core Feature 5).

## Requirements

1. MUST add to `docs/user-guide/configuration.md`, immediately before the
   `## Linux Lane` heading, a section `## Delivery Queue check wait` with this
   content, the YAML example as a fenced `yaml` block:
   - "`delivery.check_timeout` bounds each Check Wait of the Delivery Queue:
     the wait of an item in `checking` for its current-head checks and a
     mergeable state, and the wait for an existing Pull Request to report the
     candidate head. It is a Go duration, default `60m`, greater than `0` and
     at most `24h`. For example:" followed by the block `delivery:` /
     `  check_timeout: 90m`;
   - "Each wait starts its own deadline, and a re-run of a failed check
     restarts it once. At the deadline the item parks as `checks-timeout:
     current-head checks not settled within delivery.check_timeout
     <duration>`, with the duration in Go's form such as `1h0m0s`. `0`, a
     negative value or a value above `24h` fails configuration validation
     with `delivery.check_timeout must be greater than 0 and at most 24h`."
   - "Only User Config can set `delivery.check_timeout`. A Project Config
     value, in the checkout or on the fetched default branch, is ignored with
     this warning: `config: delivery.check_timeout in Project Config is
     ignored; set delivery.check_timeout in User Config`. A Queue Owner or
     `roundfix deliver retry` reads the value when it starts; an owner already
     running keeps the value it read. An older Roundfix binary refuses a User
     Config that sets `delivery.check_timeout` as an unknown key; the default
     needs no setting."
2. MUST add to the `### Review and Run settings` table of
   `docs/user-guide/configuration.md`, right after the `runs.slot_wait` row,
   the row: `delivery.check_timeout` | `60m` | "Bounds each Delivery Queue
   Check Wait for current-head checks. Go duration greater than 0 and at most
   24h. User Config only; Project Config is ignored with a warning." The Full
   example and the init templates stay as they are.
3. MUST change, identically in `docs/user-guide/commands/deliver.md` and in
   the canonical Roundfix reference `.agents/skills/roundfix/references/deliver.md`:
   - in the paragraph that begins "During `checking`, Roundfix waits", the
     sentence "The existing checks timeout parks the item as
     `checks-timeout`; only a mergeable state with passing checks proceeds to
     merge." to "At the end of the Check Wait the item parks as
     `checks-timeout: current-head checks not settled within
     delivery.check_timeout <duration>`; only a mergeable state with passing
     checks proceeds to merge.";
   - right after that paragraph, a new paragraph: "Each Check Wait lasts
     User Config `delivery.check_timeout`, a Go duration greater than `0` and
     at most `24h`, default `60m`, from its own start; a re-run of a failed
     check restarts it once. Project Config cannot set it. A Queue Owner or
     `roundfix deliver retry` reads it when it starts, so raise it in User
     Config, then run `roundfix deliver retry <slug>` for an item parked
     `checks-timeout`.";
   - the sentence ending "is read again at each check interval up to the
     check timeout." to end "is read again at each check interval up to
     `delivery.check_timeout`.".
   The Park Class table keeps `checks-timeout` in the `environment` row.
4. MUST add **Check Wait** to `CONTEXT.md` through `domain-modeling`, right
   after the **Run Slot Wait** entry, as "**Check Wait**:" followed on the
   next line by "The bounded wait of a Delivery Queue item in `checking` for
   its current-head checks and a mergeable state, and of its publication for
   an existing pull request to report the candidate head. Each wait lasts
   `delivery.check_timeout` (User Config, 60 minutes by default) from its own
   start, restarts once after a check re-run, and ends in a `checks-timeout`
   park that names the key and the duration (ADR-0311)." and the line
   "_Avoid_: CI timeout, check deadline". The ADR reference takes the
   glossary's link form, to
   `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`.
5. MUST run these commands, because this Task edits an owned skill and two
   guides. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `.agents/skills/roundfix/SKILL.md`
   by one patch over the value the tree holds when this Task starts, because
   the record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     fingerprints of the `deliver` and configuration guides change.
6. MUST search the guides, the agent guides and the skill copies for wording
   that states the check wait's value or calls it "the existing checks
   timeout", and report what it found; an older decision record that quotes a
   past `checks-timeout` park stays as written.
7. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page, and MUST NOT change text inside any `### QA
   settlement` section of any skill.

## Subtasks

- [ ] Add the check wait section and the key-reference row to the configuration guide.
- [ ] Change the park sentence and add the Check Wait paragraph in the `deliver` guide and the Roundfix `deliver` reference.
- [ ] Add **Check Wait** to the glossary.
- [ ] Raise the skill version, sync the mirrors, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The configuration guide gives the key, `60m`, the bounds, the refusal message, the User Config scope with its warning and the park text.
- [ ] The guide and both copies of the reference give the park text and the key with its default and bounds.
- [ ] `CONTEXT.md` defines **Check Wait** and links ADR-0311.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
- instruction: `docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md`
- instruction: `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/configuration.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'current-head checks not settled within delivery.check_timeout <duration>' || { printf 'missing park text in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'greater than' || { printf 'missing bound in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'config: delivery.check_timeout in Project Config is ignored; set delivery.check_timeout in User Config' || { echo 'missing warning' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'delivery.check_timeout must be greater than 0 and at most 24h' || { echo 'missing refusal' >&2; exit 1; }; grep -qE '^[|] .delivery[.]check_timeout. [|] .60m. [|]' docs/user-guide/configuration.md || { echo 'missing key row' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Each Check Wait lasts User Config' || { printf 'missing check wait paragraph in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Check Wait**:' || { echo 'missing glossary term' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document gives the park text, the key's bounds, its warning or its row, and the glossary lacks **Check Wait**, so the command fails at its first phrase check. Afterwards the configuration guide, the `deliver` guide and both reference copies give the key and the park text, the glossary defines **Check Wait**, the mirrors are in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 5; Goals; Success Metric 6; Glossary
- `_techspec.md` → API Contract 6; Build Order 1
- ADR-0236
- ADR-0304
- ADR-0311

## Result

Implemented the documentation slice for the Delivery Queue Check Wait:

- Added the `delivery.check_timeout` section and `60m` reference row to the
  configuration guide, including bounds, User Config scope, Project Config
  warning, validation refusal, and park text.
- Updated the `deliver` guide and both Roundfix `deliver` reference copies with
  the Check Wait park text, duration contract, retry guidance, and configured
  key in the existing Pull Request polling sentence.
- Added **Check Wait** to `CONTEXT.md` immediately after **Run Slot Wait**, with
  the ADR-0311 glossary link and avoided terms.
- Raised the canonical Roundfix skill from `0.1.126` to `0.1.127`, ran
  `make skills-sync`, and confirmed the canonical and embedded skill copies
  are identical.
- Re-recorded `skills/testdata/owned-skill-versions.json` with
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` (`1 passed`).
- Re-recorded `docs/references/behavior-surfaces.json` with
  `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  (`1 passed`).

Focused evidence:

- `git diff --check` passed.
- The canonical and embedded Roundfix skill files compare identical with
  `cmp`.
- A repository search found no remaining `existing checks timeout` or
  `checks timeout parks` wording in the guides, agent guides, skill copies, or
  glossary.

The Daemon-owned `## Verification` commands were not run in this handoff.
