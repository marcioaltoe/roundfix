---
spec: 0264-baseline-group-usage-errors-name-their-own-help
status: active
created: 2026-10-13
surfaces: [cli]
---

# Baseline group usage errors name their own help

This is an engineering-framed minimal PRD. The Spec fixes a defect in the
usage error of the `baseline` dispatcher and enters through the TechSpec as a
bug fix (`docs/agents/spec-routing.md`).

`roundfix baseline capabilities`, `roundfix baseline skills` and
`roundfix baseline assets`, run without one of their commands, exit 2 with a
message that doubles the word `baseline` and a pointer at a help that does not
list the group's commands. The adopted Backlog Entry
([2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md](references/2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md))
records it from the help sweep of Spec 0260.

The defect is a class in one dispatcher. A sweep on 2026-10-13, with the
binary built by `go build -buildvcs=false ./cmd/roundfix` at `3820bdd4`, ran
every group of the root usage with an unknown command word, and the three
baseline groups bare, in a temporary home and an empty working directory.
Sixteen groups answer; every failure is in `baseline`:

- `baseline capabilities`, bare or with an unknown command, exits 2 with
  `roundfix: baseline baseline capabilities failed: baseline capabilities
  requires the check command` and `Run 'roundfix baseline profile --help' for
  usage.`
- `baseline skills` and `baseline assets` fail the same way with their own
  messages.
- `baseline prune`, an unknown `baseline` command, exits 2 with
  `roundfix: baseline baseline failed: unknown baseline command "prune"` and
  the same `baseline profile --help` pointer.

The `baseline profile --help` usage lists only the three `baseline profile`
commands. The other groups, among them `runs`, `spec`, `skills`, `deliver`
and `window`, name the group once and point at a help that lists its
commands, for example `roundfix: unknown runs command "prune"` and
`Run 'roundfix runs --help' for usage.`

## Prerequisites

Spec 0260 (`0260-help-on-every-command-path-exits-zero`) must be merged
first. After it, `roundfix baseline capabilities --help`,
`roundfix baseline skills --help` and `roundfix baseline assets --help` print
the `baseline` usage and exit 0, which is the help this Spec's pointer names.
This Spec's class guard reuses the command-path derivation `helpCommandPaths`
that Spec 0260 adds in `internal/cli/command_help_test.go`, and its skill
sentence follows the sentence Spec 0260 adds to the `roundfix` skill's
baseline reference. The Task Graph names Spec 0260 in `requires`. Every
Verification command still fails on the current tree, because the tests and
the sentence it looks for do not exist yet.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. The change rewrites two lines of
  standard error; nothing reaches a network or a credential. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0256: a Behavior Surface is "the
  help of one command path"; the help text does not change, but the Skill
  Coverage Map names `internal/cli/baseline_profile.go` as a source of the
  `baseline` surfaces, so its authoring rule needs a skills Task, which
  task_01 is. ADR-0184: "Task tests and QA Requirements name the transcript by
  number and reuse its text", which task_02 and task_03 do. ADR-0259: "Every
  top-level test in a Parallel Test Package calls `t.Parallel()` as its first
  statement", which the new tests follow. ADR-0271: "A test bounds by real
  time only what it asserts"; the new tests give no child, goroutine or Git a
  real-time allowance. ADR-0189 and ADR-0233 govern the
  owned skill's version raise, which task_01 records with the repository's
  record command. ADR-0250 does not apply: no Baseline module version or
  coverage record changes. ADR-0244: "A Spec now carries a Glossary
  Declaration"; this Spec declares `None.`. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. Three Governed Paths change, all in the owned
  `roundfix` skill. The maintainer approved cycle AF on 2026-10-12 with the
  same grant as cycle AE for the Governed Paths each Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0264-baseline-group-usage-errors-name-their-own-help/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A usage error of a `baseline` group names the group once and points at the
  help of the command path the user typed, which lists the group's commands.
- An unknown `baseline` command names itself the way the other groups do and
  points at `roundfix baseline --help`.
- A repository test runs every group of the root usage with an unknown
  command, so a usage error that repeats a word or points at a help that does
  not list the group's commands fails `make test`.
- The `roundfix` skill tells agents what a bare `baseline` group answers.

## User Stories

1. As an operator who typed `roundfix baseline skills`, I want the error to
   point at a help that lists `restore` and `reconcile`, so that I find the
   command I meant without guessing.
2. As an agent reading standard error, I want the message to name the group
   once, so that I do not mistake `baseline baseline` for a command path.
3. As the maintainer, I want the test suite to check every group's usage
   error, so that this class of defect cannot return unnoticed.

## Core Features

1. `baseline capabilities`, `baseline skills` and `baseline assets`, bare or
   followed by a word that is not one of their commands, write
   `roundfix: <message>` and `Run 'roundfix baseline <group> --help' for
   usage.` to standard error, nothing to standard output, and exit 2. The
   message keeps today's wording.
2. An unknown `baseline` command writes `roundfix: unknown baseline command
   "<word>"` and `Run 'roundfix baseline --help' for usage.` and exits 2.
3. A repository test runs every group of the root usage with an unknown
   command word and requires that standard error repeats no word in a row
   and that every help it points at exits 0 and lists every command of the
   group.
4. The `roundfix` skill's baseline reference describes the bare-group usage
   error and its pointer.

## Non-Goals / Out of Scope

- The `baseline profile` paths: their errors already point at
  `baseline profile --help`, which lists their commands, and keep their text.
- JSON output for a group usage error. Today a group usage error writes
  nothing to standard output even with `--format json`, and it keeps doing so.
- Exit codes, help text, the user guide and the Behavior Surface Record's
  bytes.
- Usage text of its own for a `baseline` group, a `help <command>`
  subcommand, or a suggestion of the closest command.

## Success Metrics

1. Success Metric: the group sweep of the built binary, every group of the
   root usage with an unknown command word plus the three bare `baseline`
   groups, records 0 usage errors that repeat a word in a row or point at a
   help that does not list the group's commands, against 7 at `3820bdd4`.
2. Success Metric: the repository class-guard test fails, naming the group,
   on a tree where the `baseline` dispatcher again points at
   `baseline profile --help`, and passes at the audited head.

## Acceptance evidence

The outside-evidence row rests on published conventions this Spec did not
write, both found and read through Exa during authoring:

- Command Line Interface Guidelines (https://clig.dev/), section Errors:
  "Catch errors and rewrite them for humans. If you're expecting an error to
  happen, catch it and rewrite the error message to be useful. Think of it
  like a conversation, where the user has done something wrong and the
  program is guiding them in the right direction." A pointer at a help that
  does not list the command the user needs does not guide them, which is why
  the pointer names the typed group's help.
- Kieran Potts, "TS-16: Command Line Interfaces (CLIs)"
  (https://kieranpotts.com/standards/016.html): "Error messages MUST be sent
  to stderr whenever a command returns a non-zero exit code. The error message
  MUST provide sufficient information to help the user fix the problem
  identified by the exit code." This is why the error stays on standard error
  with exit 2 and why its pointer must lead to the group's commands.

## Research record

- Secondbrain: read `wiki/index.md`, then ran
  `qmd query "CLI usage error message names wrong help command" --all --files --min-score 0.3`.
  It returned
  `projects/skills/mirror/skills/03-engineering-design/agentic-cli-design/references/anti-patterns.md`
  (score 0.75), whose error-handling anti-patterns concern machine-readable
  failures, and the Roundfix mirror of the `doctor` command guide (0.38). No
  prior decision covers usage-error pointers. The results did not change the
  design.
- Exa: one search found the Command Line Interface Guidelines and the TS-16
  standard quoted under Acceptance evidence. Both also recommend suggesting
  the closest command; that is left out of scope, because the defect is the
  wrong pointer, not a missing suggestion.

## Glossary

None.

## Decisions

- The pointer names the command path the user typed, as every other group
  does, and relies on Spec 0260 for that help to exit 0. No ADR: the change
  is a message fix that follows the repository's existing pointer convention.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the fix, the transcripts, the class
guard and the build order.
