---
spec: 0264-baseline-group-usage-errors-name-their-own-help
prd: _prd.md
created: 2026-10-13
---

# Baseline group usage errors name their own help — Technical Spec

## Executive Summary

`runBaselineCommand` reports a bare or unknown `baseline` group command
through `printBaselineProfileFailure`, a printer written for
`baseline profile`: it prefixes the operation with `baseline` again and always
points at `baseline profile --help`. The dispatcher gets a printer of its own
that writes the message once and points at the help of the command path the
user typed, and a repository test checks the usage error of every group of
the root usage. The trade-off accepted: the pointer for a group names
`baseline <group> --help`, which exits 0 only after Spec 0260, so this Spec
is delivered after it rather than pointing at `baseline --help` and staying
independent.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. Two lines of standard error
  change; nothing reaches a network or credential. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0256: a Behavior Surface is "the
  help of one command path", and its authoring rule needs a skills Task for a
  Task that declares a covered source; `internal/cli/baseline_profile.go` is
  one, and task_01 declares the covering
  `.agents/skills/roundfix/references/baseline.md`. ADR-0184: "Task tests and
  QA Requirements name the transcript by number and reuse its text", which
  task_02 and task_03 do. ADR-0259: "Every top-level test in a Parallel Test
  Package calls `t.Parallel()` as its first statement", which both new tests
  do. ADR-0189 and ADR-0233 govern the owned skill's version raise. ADR-0250
  does not apply: no Baseline module version or coverage record changes.
  ADR-0244 is met by the `None.` declarations. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths, measured with
  `speccheck.GovernedPath` through a `go test -overlay` probe that wrote
  nothing, are granted by the cycle AF approval of 2026-10-12. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0264-baseline-group-usage-errors-name-their-own-help/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`.

## System Architecture

No new package or seam. The change touches two existing places and adds one
test file:

- `internal/cli/baseline_profile.go` — `runBaselineCommand` calls a new
  unexported printer, `printBaselineUsageError`, at its five usage-error
  sites: the `capabilities`, `skills` (twice) and `assets` cases and the
  `default` case. `printBaselineProfileFailure` keeps every other caller, all
  of them `baseline profile` operations, and its text.
- `internal/cli/baseline_group_usage_test.go` — new: the transcript test and
  the class guard over every group of the root usage.
- `.agents/skills/roundfix/references/baseline.md` and its mirror — one
  sentence on the bare-group usage error.

Absence evidence for the new components, at `3820bdd4`:
`grep -rn 'printBaselineUsageError\|TestEveryGroupUsageErrorPointsAtItsCommands\|TestBaselineGroupUsageTranscripts' internal/`
exits 1 with no match, and `ls internal/cli/baseline_group_usage_test.go`
reports no such file.

Relation to Spec 0260, which is being delivered: it edits the same
`capabilities`, `skills` and `assets` cases of `runBaselineCommand` to answer
a help token before the usage error, and it adds a sentence to the same skill
reference. This Spec's task_02 is implemented on the tree after Spec 0260
merges and replaces only the usage-error calls below Spec 0260's help
branches; it does not edit `internal/cli/cli_test.go` or
`internal/cli/command_help_test.go`, which Spec 0260 owns, and reads only the
`helpCommandPaths` helper from the latter.

Current behavior, characterized at `3820bdd4` with a sweep of every group of
the root usage given the unknown word `no-such-command`, plus the three bare
`baseline` groups: 16 groups, 19 calls. Seven calls, all in `baseline`, repeat
`baseline baseline` and point at `baseline profile --help`, whose usage lists
none of the group's commands. The other twelve repeat no word, and every help
they point at exits 0 and lists every command of its group;
`history no-such-command` points at `history sanitize --help`, which lists
the group's only command. `review no-such-command` prints no pointer. No
existing test asserts the old text: a search of `internal/cli/*_test.go` for
`baseline baseline`, `requires the check command`, `requires the sync
command`, `requires the restore or reconcile` and `unknown baseline command`
finds nothing, so no existing test, golden or fixture changes.

## Implementation Design

### Interfaces

```go
// printBaselineUsageError writes a baseline dispatcher usage error to stderr:
// the message once, then the help pointer of the command path the user typed.
func printBaselineUsageError(commandPath, message string, stderr io.Writer)
```

Test-only helper in the new test file:

```go
// usageErrorGroups returns each command path of helpCommandPaths(rootUsage)
// that prefixes another path, mapped to the command paths it prefixes.
func usageErrorGroups(rootUsage string) map[string][]string
```

Invariants:

1. A `baseline` usage error writes nothing to standard output, with or
   without `--format json`.
2. A `baseline` usage error exits 2 (`exitPreflight`), as today.
3. Standard error of a usage error never repeats a word in a row.
4. Every `'roundfix <path> --help'` pointer a usage error prints names a help
   that exits 0 and lists every command of the group the user typed.
5. The `baseline profile` paths, every help text and
   `docs/references/behavior-surfaces.json` keep their bytes.

### Data Models

None. No stored or serialized shape changes.

### API Contracts

1. API Contract 1: `roundfix baseline capabilities`, `baseline skills` and
   `baseline assets`, with no further argument or with a first argument that
   is neither one of the group's commands nor a help token Spec 0260 answers,
   write `roundfix: <message>` and
   `Run 'roundfix baseline <group> --help' for usage.` to standard error,
   nothing to standard output, and exit 2. `<message>` keeps today's text:
   `baseline capabilities requires the check command`,
   `baseline skills requires the restore or reconcile command` and
   `baseline assets requires the sync command`.
2. API Contract 2: `roundfix baseline <word>`, where `<word>` is not a
   `baseline` command and does not start with `-`, writes
   `roundfix: unknown baseline command "<word>"` and
   `Run 'roundfix baseline --help' for usage.` to standard error and exits 2.
3. API Contract 3: `TestEveryGroupUsageErrorPointsAtItsCommands` derives the
   command paths from `roundfix --help` with `helpCommandPaths`, keeps each
   path that prefixes another as a group (at least 12), and runs each group
   with the extra word `no-such-command`. For each it requires a non-zero
   exit, no word of standard error repeated in a row, and, for each
   `'roundfix <path> --help'` pointer in standard error, that
   `<path> --help` exits 0 and its standard output contains
   `roundfix <command>` for every command path the group prefixes. A failure
   names the command line, the repeated word or the pointer and the missing
   command.

### Surface Transcripts

1. Surface Transcript: a bare `baseline capabilities` points at its own help.

   ```transcript
   $ roundfix baseline capabilities
   stdout:
   stderr:
   roundfix: baseline capabilities requires the check command
   Run 'roundfix baseline capabilities --help' for usage.
   exit: 2
   ```

2. Surface Transcript: a bare `baseline skills` points at its own help.

   ```transcript
   $ roundfix baseline skills
   stdout:
   stderr:
   roundfix: baseline skills requires the restore or reconcile command
   Run 'roundfix baseline skills --help' for usage.
   exit: 2
   ```

3. Surface Transcript: an unknown command after `baseline skills` answers the
   same way.

   ```transcript
   $ roundfix baseline skills prune
   stdout:
   stderr:
   roundfix: baseline skills requires the restore or reconcile command
   Run 'roundfix baseline skills --help' for usage.
   exit: 2
   ```

4. Surface Transcript: a bare `baseline assets` points at its own help.

   ```transcript
   $ roundfix baseline assets
   stdout:
   stderr:
   roundfix: baseline assets requires the sync command
   Run 'roundfix baseline assets --help' for usage.
   exit: 2
   ```

5. Surface Transcript: a JSON request still writes nothing to standard
   output.

   ```transcript
   $ roundfix baseline capabilities --format json
   stdout:
   stderr:
   roundfix: baseline capabilities requires the check command
   Run 'roundfix baseline capabilities --help' for usage.
   exit: 2
   ```

6. Surface Transcript: an unknown `baseline` command names itself once.

   ```transcript
   $ roundfix baseline prune
   stdout:
   stderr:
   roundfix: unknown baseline command "prune"
   Run 'roundfix baseline --help' for usage.
   exit: 2
   ```

## Coverage Map

Cross-read from `_prd.md`:

- Goal 1, User Stories 1-2, Core Feature 1 → `printBaselineUsageError` at
  the four group sites of `runBaselineCommand`; Surface Transcripts 1-5; API
  Contract 1; Invariants 1-2.
- Goal 2, User Story 2, Core Feature 2 → `printBaselineUsageError` at the
  `default` site; Surface Transcript 6; API Contract 2.
- Goal 3, User Story 3, Core Feature 3, Success Metrics 1-2 →
  `TestEveryGroupUsageErrorPointsAtItsCommands`; API Contract 3;
  Invariants 3-4.
- Goal 4, Core Feature 4 → the `roundfix` skill's baseline reference.

Every PRD goal, story, feature and metric has a row; no TechSpec component is
unmapped.

## Integration Points

None external. The Behavior Surface Record contract (`make verify-docs`)
fingerprints help, which does not change (Invariant 5).

## Testing Approach

- `TestBaselineGroupUsageTranscripts`, new in
  `internal/cli/baseline_group_usage_test.go`, asserts Surface Transcripts
  1-6 through `runCLI`, one subtest per transcript named by its number: the
  empty standard output, every standard-error line in order and exit 2. For
  transcripts 1-4 and 6 it then runs the pointed command with `--help` and
  requires exit 0, so the transcript proves its pointer resolves.
- `TestEveryGroupUsageErrorPointsAtItsCommands`, new in the same file, is the
  class guard of API Contract 3. On a tree with Spec 0260 and without this
  fix it fails on the seven `baseline` calls the characterization lists; on a
  tree without Spec 0260 the file does not compile, because
  `helpCommandPaths` arrives with that Spec.
- Both tests call `t.Parallel()` first (ADR-0259) and run in the per-test
  environment of `runCLI`; the unknown word reaches only argument parsing.
- Existing `baseline_profile_test.go` cases stay unchanged and keep passing,
  which proves the `baseline profile` half of Invariant 5.

## Build Order

1. The skill sentence and the recorded version raise (task_01). It follows
   the sentence Spec 0260 adds to the same reference, so it requires Spec
   0260 merged.
2. The printer and its five call sites, the transcript test and the class
   guard (task_02) (depends on: 1). It requires Spec 0260 merged; it depends
   on step 1 because the Task that changes the `baseline` command surface
   needs the Task that documents it in its ancestry.
3. Final QA gate (task_03) (depends on: 2).

## Risks & Considerations

- The class guard runs a real command for every group with an unknown word.
  A future group that acts before it parses would run inside `runCLI`'s test
  environment; the failure message names it.
- If Spec 0260 is cancelled, the group pointers name a help that exits 2;
  task_02's transcript test then fails, and this Spec must change its pointer
  to `roundfix baseline --help` before delivery.
- The skill version raise is written only by the record command; task_01
  declares its record file.

## Glossary

None.

## Decisions

- A dedicated printer rather than a flag on `printBaselineProfileFailure`:
  the `baseline profile` callers keep their exact text, and the dispatcher's
  errors take the shape every other group already uses.
- The group pointer names the typed command path. Spec 0260 makes that help
  print the `baseline` usage, which lists the group's commands; pointing at
  `roundfix baseline --help` was the alternative that avoids the dependency,
  rejected because it breaks the convention every other group follows.
- No ADR: no decision here is hard to reverse, and none departs from
  ADR-0256, ADR-0184 or ADR-0259, whose rejected alternatives this design
  does not revive.
- Research: the Secondbrain and Exa consultations are recorded in
  `_prd.md` → Research record; neither changed the design.
