---
task: task_02
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
status: pending
type: docs
complexity: medium
---

# Task 02: The Roundfix skill reads its Skill Mode from the request and reports status

## Overview

The Roundfix skill has no modes. This Task gives its entry file a modes
section that names `--follow-up`, `--status` and `--unattended` and the way
Claude Code, Codex and OpenCode pass them: the first word after the skill's
name selects the mode with or without a leading `--` (`/roundfix follow-up`,
`$roundfix status`, `@roundfix --unattended`), which is not the CLI command
`roundfix follow-up`, and any other first word keeps today's behavior. It
exposes the skill as an OpenCode
slash command, and writes the first half of a new topic reference: how a mode
is chosen, how each optional source is detected, and the read-only `--status`
mode. It adds **Skill Mode** to the glossary and a test that holds the
structure (`_prd.md` → Core Features 2, 3 and 6; ADR-0320).

## Requirements

1. MUST answer the Backlog Entry Spec 0306 adopted,
   [`../0306-the-follow-up-command-lists-pending-work-in-priority-order/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md`](../../references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md)
   of 2026-10-22, decisions 1 and 2 and the harness table, and the
   maintainer's decision of 2026-10-24 that `_authorization.md` quotes: the
   first word after the skill's name selects the mode with or without a
   leading `--` in all three harnesses.
2. MUST change the canonical `.agents/skills/roundfix/SKILL.md` exactly as
   `_techspec.md` → API Contract 3 says: the `opencode/slash: "true"` line
   under `metadata`, the description clause, the `## Modes` section with its
   exact first paragraph, one bullet per mode (`--follow-up`: inventory, plan,
   gather decisions, author and deliver; `--status`: a short read-only report;
   `--unattended`: gather every decision in calls of up to three questions)
   and the line linking the reference, and the index row between the
   `implement` and `profiles` rows. Every other byte of the existing sections
   stays, apart from the version fields. The entry file stays within its
   20,000-byte budget and its body never contains a `$` followed by a digit
   or by `ARGUMENTS`, because Claude Code substitutes such text with the
   request's arguments.
3. MUST create `.agents/skills/roundfix/references/modes.md` starting with `#
   Skill Modes` and one introductory paragraph, followed by these sections:
   - `## How a mode is chosen`: Claude Code appends the text typed after the
     skill's name to the skill as an `ARGUMENTS:` line, while Codex and
     OpenCode keep it in the message, so the skill reads the mode from the
     request text in all three; the mode is the first word after the skill's
     name when that word is `follow-up`, `status` or `unattended`, with or
     without a leading `--`, so `/roundfix follow-up`, `/roundfix
     --follow-up`, `$roundfix status` and `@roundfix --unattended` each select
     a mode; this is the Skill Mode, not the CLI command `roundfix follow-up`,
     which the follow-up mode runs as one of its steps; a request whose first
     word is anything else, or that has none, uses the skill as before (for
     example `/roundfix watch 123`), and a mode word later in such a request
     selects nothing; when the words after the mode word name another mode,
     with or without `--`, the skill asks one question that names them,
     recommending `follow-up` when it is named and `status` otherwise; the
     section contains, on whitespace-normalized text, the phrases "first word
     after the skill's name" and "with or without a leading `--`" and the
     code span `` `roundfix follow-up` ``; each mode runs Roundfix's read-only follow-up command, and
     when the installed Roundfix lacks it the mode says so, names `roundfix
     upgrade` and stops.
   - `## Optional sources`: the Secondbrain is a source when the repository
     has `docs/agents/secondbrain.md` and the follow-up command can read the
     Secondbrain its configuration names (Spec 0306 plans the User Config key
     `secondbrain.root` and the inbox namespace key
     `secondbrain.inbox_namespace`); ai-memory is a source when the session
     exposes its MCP tools, and the skill lists this project's messages
     without consuming them, writes them to a JSON file in a temporary
     directory in the message-file schema the command documents (Spec 0306
     plans `roundfix/follow-up-messages/v1`) and hands the file to the
     command, and consumes a message only after it is recorded as a Spec,
     Backlog Entry, Finding or answer; Jev is used by the follow-up command
     when a judge key is configured, and the skill relays the command's own
     report when it says Jev was unavailable; a release-announcement step
     exists only when the user's own instructions define one, and the skill
     never names one; a missing source never stops a mode, and its report
     names each missing source and why.
   - `## --status`: run the follow-up command's status view and report, in a
     few lines, the work in progress, the pending Specs in delivery order, the
     next release and what it brings, the blockers, and what waits on the
     user; it writes nothing, asks nothing, consumes no ai-memory message and
     arms nothing.
4. MUST name the follow-up command, its status view, its message-file
   option, the message-file schema and the configuration keys exactly as `go run -buildvcs=false ./cmd/roundfix --help`, the
   command's own `--help` and its Roundfix reference state in the tree this
   Task starts on (Spec 0306 delivers them before this Spec and plans them
   as `roundfix follow-up [--status] [--messages <path>] [--format
   <text|json>]`), link the command's Roundfix reference from the modes
   reference for its flags, and MUST NOT name a command or flag the tree
   lacks. When the tree has no follow-up
   command, stop and report instead of describing one.
5. MUST raise `TestRoundfixSkillIndexNamesEveryReference` in
   `skills/roundfix_layout_test.go` to expect one more reference file than
   the value the tree holds when this Task starts, in both the count and its
   message, and change nothing else in that file.
6. MUST create `skills/roundfix_modes_test.go` with `roundfixSection` and
   `TestTheRoundfixSkillReadsItsModeFromTheRequest` as `_techspec.md` →
   Interfaces and Testing Approach describe them, over
   `roundfixLayoutSources()`, reading the frontmatter `metadata` with
   `gopkg.in/yaml.v3`. The top-level test and each subtest call `t.Parallel()`
   first.
7. MUST add **Skill Mode** to `CONTEXT.md` through `domain-modeling`, right
   after the **Supervisor** entry, as "**Skill Mode**:" followed on the next
   line by "A way of running the Roundfix skill that the user's request
   selects with the first word after the skill's name, `follow-up`,
   `status` or `unattended`, with or without a leading `--`; the skill reads
   it from the request text in every harness, it is not the CLI command
   `roundfix follow-up`, and a request whose first word names no mode uses
   the skill as before." with a link to
   ADR-0320 in the glossary's link form, and the line "_Avoid_: Command flag,
   CLI mode, Follow Mode".
8. MUST run these commands, because this Task edits an owned skill. Before
   them, raise both version fields (`metadata.version` and the top-level
   `version`) of the canonical `SKILL.md` by one patch over the value the tree
   holds when this Task starts, because the record command refuses changed
   content under an unchanged version. The Agent runs them itself before
   handing back, in this order, and again after any later edit of the skill;
   no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
9. MUST NOT add the `## --follow-up` or `## --unattended` sections (task_03
   writes them), change any other owned skill, `skills/roundfix/agents/` or
   `.agents/skills/roundfix/agents/`, any user guide page, the follow-up
   command's reference, or any text inside a `### QA settlement` section.

## Subtasks

- [ ] Add the metadata, description, modes section and index row to the entry file.
- [ ] Write the first half of the modes reference with the tree's follow-up command names.
- [ ] Raise the layout count, add the mode test and the glossary term.
- [ ] Raise the skill version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The canonical and embedded skill name the three modes and invocation forms, each mode word bare and with `--`, say the mode is not the CLI command `roundfix follow-up` and that any other first word keeps today's behavior, link the modes reference, set the OpenCode slash metadata and carry no argument placeholder.
- [ ] The modes reference explains how a mode is chosen, how each optional source is detected and what `--status` reports, with the tree's command names.
- [ ] `CONTEXT.md` defines **Skill Mode**; the skill version is recorded and the mirrors are in sync.

## Consumers

- the Roundfix skill's reference index — consumers: `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0320-a-question-call-carries-up-to-three-independent-questions.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/modes.md`
- creates: `skills/roundfix/references/modes.md`
- interface: `skills/roundfix_layout_test.go`
- creates: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Verification

- `go test -list '^TestTheRoundfixSkillReadsItsModeFromTheRequest$' ./skills | grep -qx TestTheRoundfixSkillReadsItsModeFromTheRequest || { echo 'missing mode test' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Skill Mode**:' || { echo 'missing glossary term' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestEveryCommandIsNamedInTheRoundfixSkill)$' ./internal/docscontract` — expected: exit 0. Before this Task the mode test does not exist, so the command fails at its first check. Afterwards the test exists, the glossary defines **Skill Mode**, the mirrors are in sync, the skills package passes with the raised reference count, the mode structure and the recorded version, and the documentation contracts that parse the reference index pass.

## References

- `_prd.md` → Goals; User Stories 1, 3 and 4; Core Features 2, 3 and 6; Success Metric 3; Glossary
- `_techspec.md` → Current behavior 5-6; API Contracts 3, 4 and 7; Invariant 4; Testing Approach; Build Order 2
- ADR-0320
