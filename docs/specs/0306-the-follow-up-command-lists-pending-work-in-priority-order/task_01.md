---
task: task_01
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: pending
type: docs
complexity: medium
---

# Task 01: The command guide, the follow-up reference, the configuration guide and the glossary describe the Follow-up Command

## Overview

The Follow-up Command, its JSON, its two User Config keys and its order are
new, so no document describes them. This Task writes the command guide, the
Roundfix skill's new `follow-up` reference with its index row, the
configuration section and key rows, the command index and usage rows, the
Skill Coverage Map entry of the new guide, and the four glossary terms,
before the code that later Tasks add (`_prd.md` → Core Feature 8;
`_techspec.md` → API Contract 9). It raises the Roundfix skill version and
re-records the Behavior Surface Record.

## Requirements

1. MUST create `docs/user-guide/commands/follow-up.md` starting with the
   heading `### follow-up` and a `bash` block holding
   `roundfix follow-up [--status] [--messages <path>] [--format <text|json>]`,
   and describe, in the guide's plain style, everything `_techspec.md` → API
   Contracts 1-8 define: the read-only boundary, with the sentence "It is
   read-only: it never writes a file, commits, pushes or prompts."; the exit
   codes; each source by its JSON name (`secondbrain-inbox`,
   `ai-memory-messages`, `backlog`, `findings`, `repository-inbox`, `specs`,
   `runs`, `delivery-queue`, `pull-requests`, `git`) and its `checked`,
   `missing` and `unavailable` states; the messages file and its schema
   `roundfix/follow-up-messages/v1` with one example; the four questions and
   the classes `unblocking`, `health`, `feature` and `other` with the
   sentence "Without an answer, an item's declared type orders it, and
   nothing is unblocking."; the stop line "Nothing to do: no pending item and
   nothing in progress." and the JSON `state` values `work_found` and
   `nothing_to_do`; the JSON report schema `roundfix/follow-up/v1` with one
   example of each view; the status view's sections, including the line
   "Next release: no unreleased changes"; and that the judge sends nothing
   inside a Run, without a key, at the monthly ceiling or for an item that
   matches a secret pattern.
2. MUST create the canonical Roundfix reference
   `.agents/skills/roundfix/references/follow-up.md` with a heading
   `## Follow-up Command` and the agent-facing contract: when to run
   `roundfix follow-up --format json` and `roundfix follow-up --status
   --format json`, how to hand ai-memory messages in with `--messages <path>`
   in the `roundfix/follow-up-messages/v1` schema, how to read `state`,
   `sources`, `judge` and `items`, the four classes and the type fallback,
   the two keys `secondbrain.root` and `secondbrain.inbox_namespace`, and
   that the command never writes to the Secondbrain, ai-memory or the
   repository. It MUST NOT describe the skill modes, which Spec 0307 adds.
3. MUST add to the reference index of `.agents/skills/roundfix/SKILL.md`,
   between the `events` and `implement` rows, one row whose Reference cell
   is `[follow-up](references/follow-up.md)`, whose Commands covered cell is
   the code span `follow-up`, and whose When to read cell is "Listing pending
   work in priority order, or the status of work in progress.", and change `want 18 command families` and the count `18` it guards in
   `skills/roundfix_layout_test.go` (`TestRoundfixSkillIndexNamesEveryReference`)
   to `19`.
4. MUST add to `docs/user-guide/configuration.md` a section `## Secondbrain
   inbox` after `## Jev monthly ceiling`, with a YAML example of both keys,
   the Project Config warning `config: secondbrain.root in Project Config is
   ignored; set secondbrain.root in User Config`, both validation messages of
   `_techspec.md` → API Contract 2, the default namespace rule, and the
   sentence "An older Roundfix binary refuses a User Config containing these
   keys as unknown keys."; and add the rows `secondbrain.root` (default
   unset) and `secondbrain.inbox_namespace` (default the repository's name)
   to the `## Key reference` table.
5. MUST add to the command index of `docs/user-guide/commands.md`, between
   the `fetch` and `gc` rows, one row whose Command cell is the code span
   `follow-up` and whose Guide cell is `[command guide](commands/follow-up.md)`,
   and to the `## Command reference` table of `docs/user-guide/usage.md`,
   right after the `release plan` row, one row whose first cell is the code span `follow-up` and
   whose second cell is "List pending work in priority order, or show status
   (read-only)".
6. MUST add to `docs/references/skill-coverage.json` the entry
   `guide docs/user-guide/commands/follow-up.md` whose `skills` is
   `.agents/skills/roundfix/references/follow-up.md` and whose `sources` is
   the guide itself; add `.agents/skills/roundfix/references/follow-up.md` to
   the `skills` of the `config keys`, `guide docs/user-guide/configuration.md`,
   `guide docs/user-guide/commands.md` and `guide docs/user-guide/usage.md`
   entries; and record the Coverage Review of `_prd.md` → Skills by setting
   each of those four entries' `review` to `2026-10-22 — ` followed by that
   surface's reason in `_prd.md` → Skills. It MUST NOT add the
   `command follow-up` entry, which task_03 adds with the command.
7. MUST add to `CONTEXT.md` through `domain-modeling`, right after the
   **Triage** entry, four entries, each as "**Term**:" with its definition on
   the next line and an `_Avoid_` line:
   - "**Follow-up Command**:" — "The read-only `roundfix follow-up` command
     that lists a repository's Follow-up Items in Follow-up Priority order
     beside its work in progress, names each source it checked or missed,
     reports `nothing_to_do` when it finds neither, and with `--status` shows
     the work in progress, the pending Specs in delivery order, the next
     release and what waits on the maintainer. It never writes, and the judge
     never changes its exit code (ADR-0319)." — `_Avoid_: Triage command,
     dashboard, inbox sync`.
   - "**Follow-up Item**:" — "One pending unit the Follow-up Command lists: a
     pending Inbox Entry of the repository's inbox namespace, a message the
     caller hands in, an open Backlog Entry, an unresolved Finding, a note in
     `docs/_inbox/` or an active Spec, with its kind, reference, type, date,
     title and the project that sent it." — `_Avoid_: Work Item, ticket,
     to-do`.
   - "**Follow-up Priority**:" — "The order of Follow-up Items in four
     classes, `unblocking`, `health`, `feature` and `other`, computed in code
     from four Jev yes/no answers per item, with small items and then the
     oldest first inside a class; without an answer an item's declared type
     places it, and never in `unblocking` (ADR-0319)." — `_Avoid_: Jev
     ranking, severity, score`.
   - "**Secondbrain Root**:" — "The local Secondbrain directory that User
     Config `secondbrain.root` names; the Follow-up Command reads the
     pending Inbox Entries under its `inbox/<namespace>/` and never writes,
     fetches or pulls it." — `_Avoid_: Brain path, knowledge workspace`.
   Each ADR reference takes the glossary's link form, to
   `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`.
8. MUST run these commands, because this Task edits an owned skill and
   guides. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `.agents/skills/roundfix/SKILL.md`
   by one patch over the value the tree holds when this Task starts, because
   the record command refuses changed content under an unchanged version.
   The Agent runs them itself before handing back, in this order, and re-runs
   the last one after the last guide edit; no mirror or record is
   hand-edited:
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     new guide and the edited `commands.md`, `usage.md` and
     `configuration.md` change their fingerprints.
9. MUST NOT change any Go source other than the count in
   `skills/roundfix_layout_test.go`, any ADR, any other guide page or
   reference, or any text inside a `### QA settlement` section.

## Subtasks

- [ ] Write the command guide and the configuration section and key rows.
- [ ] Write the `follow-up` reference, its index row and the layout count.
- [ ] Add the command index and usage rows and the Skill Coverage Map entries.
- [ ] Add the four glossary terms.
- [ ] Raise the skill version, sync the mirrors, record the version and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The guide and both copies of the reference describe the sources, the order, the stop rule, both views, the JSON and the keys.
- [ ] The reference index names the `follow-up` family once, and the layout test counts 19 references.
- [ ] `CONTEXT.md` defines **Follow-up Command**, **Follow-up Item**, **Follow-up Priority** and **Secondbrain Root**, linking ADR-0319.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record and the Skill Coverage Map are current.

## Consumers

None. This Task changes documentation, the skill mirror, the layout test's
count, the skill version record, the Skill Coverage Map and the Behavior
Surface Record; the suites that read them run in its own Verification.

## Context

- instruction: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/user-guide/commands/cancel.md`
- creates: `docs/user-guide/commands/follow-up.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- creates: `.agents/skills/roundfix/references/follow-up.md`
- creates: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix_layout_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `test -f docs/user-guide/commands/follow-up.md || { echo 'missing follow-up guide' >&2; exit 1; }; g="$(tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md)"; for p in 'roundfix follow-up [--status] [--messages <path>] [--format <text|json>]' 'It is read-only: it never writes a file, commits, pushes or prompts.' 'roundfix/follow-up-messages/v1' 'roundfix/follow-up/v1' 'Without an answer, an item' 'Nothing to do: no pending item and nothing in progress.' 'Next release: no unreleased changes' 'nothing_to_do' 'unblocking' 'secondbrain-inbox'; do printf '%s\n' "$g" | grep -qF -- "$p" || { printf 'missing in guide: %s\n' "$p" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do grep -q '^## Follow-up Command$' "$f" || { printf 'missing heading in %s\n' "$f" >&2; exit 1; }; r="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'roundfix follow-up --format json' 'roundfix follow-up --status --format json' 'roundfix/follow-up-messages/v1' 'secondbrain.root' 'secondbrain.inbox_namespace' 'nothing_to_do'; do printf '%s\n' "$r" | grep -qF -- "$p" || { printf 'missing in %s: %s\n' "$f" "$p" >&2; exit 1; }; done; done; for f in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do grep -qE -- '^[|] \[follow-up\]\(references/follow-up.md\) [|] .follow-up. [|]' "$f" || { printf 'missing index row in %s\n' "$f" >&2; exit 1; }; done; grep -qE -- '^[|] .follow-up. [|] \[command guide\]\(commands/follow-up.md\) [|]' docs/user-guide/commands.md || { echo 'missing command index row' >&2; exit 1; }; grep -qE -- '^[|] .follow-up. [|]' docs/user-guide/usage.md || { echo 'missing usage row' >&2; exit 1; }; grep -q '^## Secondbrain inbox$' docs/user-guide/configuration.md || { echo 'missing configuration section' >&2; exit 1; }; c="$(tr -s '[:space:]' ' ' < docs/user-guide/configuration.md)"; for p in 'config: secondbrain.root in Project Config is ignored; set secondbrain.root in User Config' 'secondbrain.root must be an absolute path or start with ~/' 'An older Roundfix binary refuses a User Config containing these keys as unknown keys.'; do printf '%s\n' "$c" | grep -qF -- "$p" || { printf 'missing in configuration guide: %s\n' "$p" >&2; exit 1; }; done; grep -qF '"guide docs/user-guide/commands/follow-up.md"' docs/references/skill-coverage.json || { echo 'missing guide map entry' >&2; exit 1; }; test "$(grep -c 'references/follow-up.md' docs/references/skill-coverage.json)" -ge 5 || { echo 'follow-up reference not added to the config entries' >&2; exit 1; }; grep -q 'want 19 command families' skills/roundfix_layout_test.go || { echo 'layout count not raised' >&2; exit 1; }` — expected: exit 0. Before this Task no follow-up guide exists, so the command fails at its first check. Afterwards the guide, both reference copies, the index rows, the configuration section, the map entries and the layout count are present.
- `x="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; for t in '**Follow-up Command**:' '**Follow-up Item**:' '**Follow-up Priority**:' '**Secondbrain Root**:'; do printf '%s\n' "$x" | grep -qF -- "$t" || { printf 'missing glossary term %s\n' "$t" >&2; exit 1; }; done; printf '%s\n' "$x" | grep -qF -- 'docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md' || { echo 'glossary does not link ADR-0319' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary lacks **Follow-up Command**, so the command fails at its first check. Afterwards the four terms link ADR-0319, the mirrors are in sync, the skill tests (layout count and recorded version among them) pass, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 8; Goals; User Story 3; Glossary; Skills
- `_techspec.md` → API Contracts 1-9; Build Order 1
- ADR-0256
- ADR-0319
