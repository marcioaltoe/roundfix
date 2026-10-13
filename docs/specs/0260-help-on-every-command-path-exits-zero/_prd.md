---
spec: 0260-help-on-every-command-path-exits-zero
status: active
created: 2026-10-12
surfaces: [cli]
---

# Help on every command path exits zero

This is an engineering-framed minimal PRD. The Spec fixes a defect in how
four command paths answer a help request and enters through the TechSpec as a
bug fix (`docs/agents/spec-routing.md`).

`roundfix qa-report accept --help` exits 1. The command reads `--help` as the
report path and refuses it: `roundfix: QA Report eligibility refused: report
path "--help" is not in a Spec qa directory`. The adopted Backlog Entry
([2026-10-09-qa-report-accept-help-exits-one.md](references/2026-10-09-qa-report-accept-help-exits-one.md))
records it, and Spec 0251's author met it while fingerprinting command help.
The Behavior Surface Record therefore fingerprints `qa-report accept` from its
parent's help, through a fallback in the contract test that reads the nearest
ancestor whose help succeeds.

The defect is a class, not one command. A sweep on 2026-10-12, with
`bin/roundfix` built by `make build` at `7a8b7f70`, asked `--help` and `-h` of
every command path the root usage names and of every group prefix of one: 70
paths, 140 calls, in a temporary home and an empty working directory. Eight
calls failed, the same way for both flags:

- `qa-report accept --help` and `-h` exit 1 with the refusal above.
- `baseline capabilities --help` and `-h` exit 2 with `roundfix: baseline
  baseline capabilities failed: baseline capabilities requires the check
  command` and `Run 'roundfix baseline profile --help' for usage.`
- `baseline skills --help` and `-h`, and `baseline assets --help` and `-h`,
  exit 2 the same way.

The other 132 calls exit 0 with a usage block on standard output and nothing
on standard error.

## Prerequisites

None. No Verification reads an artifact of another active Spec.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added; help is printed from constants. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0256 defines a Behavior Surface as
  "the help of one command path", so this Spec fingerprints `qa-report accept`
  from its own help and drops the ancestor fallback; its authoring rule needs a
  skills Task for a changed surface's source, which task_01 carries. ADR-0184:
  "A TechSpec now declares numbered Surface Transcripts", which the TechSpec
  does. ADR-0259: "Every top-level test in a Parallel Test Package calls
  `t.Parallel()` as its first statement", which the new tests follow.
  ADR-0189 and ADR-0233 govern the owned skill's version raise, which task_01
  records with the repository's record command. ADR-0250 does not apply: no
  Baseline module version or coverage record changes. ADR-0244: "A Spec now carries
  a Glossary Declaration"; this Spec declares `None.`. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. Four Governed Paths change: the `internal/cli`
  test file that holds the help-exit test and three files of the owned
  `roundfix` skill. The maintainer granted the Governed Paths each cycle AE
  Spec declares on 2026-10-11. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0260-help-on-every-command-path-exits-zero/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- `--help` and `-h` on every command path and every group prefix print a
  usage block on standard output, print nothing on standard error, and exit 0.
- A repository test asks every command path for help, so a new command that
  forgets help fails `make test`.
- The Behavior Surface Record fingerprints each command path from its own
  help, and its contract fails when that help does not exit 0.
- The `roundfix` skill stops telling agents that `accept` has no `--help`.

## User Stories

1. As an agent reading a command's contract, I want
   `roundfix qa-report accept --help` to print its usage and exit 0, so that I
   do not have to know the command misreads the flag.
2. As an operator exploring `baseline`, I want `baseline skills --help` to
   show the commands I can run, so that a help request is not reported as a
   failure.
3. As the maintainer, I want the test suite to ask every command path for
   help, so that this class of defect cannot return unnoticed.

## Core Features

1. `qa-report accept` treats `--help` or `-h`, before or after the report
   path, as a help request: it prints the `qa-report` usage and exits 0
   without reading a report.
2. `baseline capabilities`, `baseline skills` and `baseline assets` followed
   by `--help` or `-h` print the `baseline` usage, which lists their
   commands, and exit 0.
3. A repository test derives every command path and group prefix from the
   root usage and requires both help flags to exit 0 with a usage block and
   an empty standard error; the help-exit test names the four fixed paths.
4. The Behavior Surface Record's contract reads only the command path's own
   help and fails when it exits non-zero.
5. The `roundfix` skill's settle and baseline references describe the new
   help behavior.

## Non-Goals / Out of Scope

- The usage-error text of a bare `baseline capabilities`, `baseline skills`
  or `baseline assets`, which doubles the word `baseline` and points at
  `baseline profile --help`. That error path keeps its text and exit 2; a
  Backlog Entry records it.
- New usage text for any command, or a `help <command>` subcommand.
- Any change to the user guide, the exit-code table or the Behavior Surface
  Record's bytes: the `qa-report accept` usage is the `qa-report` usage, so
  its fingerprint does not change.
- Fingerprinting group prefixes in the Behavior Surface Record.

## Success Metrics

1. Success Metric: the help sweep of the built binary, `--help` and `-h` over
   every command path and group prefix of the root usage, records 0 failing
   calls at the audited head, against 8 at `7a8b7f70`.
2. Success Metric: the repository help sweep test fails, naming the command,
   on a tree where `qa-report accept` again reads `--help` as a path, and
   passes at the audited head.
3. Success Metric: the Behavior Surface Record contract passes at the audited
   head with the record unchanged, and fails, naming `qa-report accept`, on a
   tree where that command's help exits non-zero.

## Acceptance evidence

The outside-evidence row rests on published conventions this Spec did not
write, both found and read through Exa during authoring:

- GNU Coding Standards, "Standards for Command Line Interfaces"
  (https://web.mit.edu/gnu/doc/html/standards_18.html): programs should
  support a help option "which prints option usage information on standard
  output and exits successfully", and "These options should inhibit the normal
  function of the command; they should do nothing except print the requested
  information." This is why help exits 0, prints to standard output and reads
  no report.
- Command Line Interface Guidelines (https://clig.dev/): "Show full help when
  `-h` and `--help` are passed", and "Ignore any other flags and arguments that
  are passed—you should be able to add `-h` to the end of anything and it
  should show help." This is why `-h` is swept with `--help`, why a group
  prefix answers help, and why a help flag after the report path still prints
  help.

## Research record

- Secondbrain: read `wiki/index.md`, then ran
  `qmd query "CLI subcommand --help exit code zero usage convention" --all --files --min-score 0.3`.
  It returned the Roundfix mirror of the adopted entry
  (`projects/roundfix/mirror/docs/backlog/2026-10-09-qa-report-accept-help-exits-one.md`),
  the mirror of the `spec` command guide, and two upstream skill mirrors:
  `projects/skills/mirror/skills/05-implementation-loop/golang-cli/SKILL.md`,
  whose common-mistakes table says to "save full usage for `--help`", and
  `projects/skills/mirror/skills/03-engineering-design/agentic-cli-design/SKILL.md`.
  No prior decision covers help exit codes. The results confirmed the defect
  and supported printing usage only when asked; they changed nothing in the
  design.
- Exa: two searches and fetches found the GNU Coding Standards and the
  Command Line Interface Guidelines quoted under Acceptance evidence. They
  set the exit status, the stream, and the rule that a help flag anywhere
  wins, which is why `-h` and a trailing flag are covered too. The official
  `gnu.org` page could not be fetched; the MIT mirror carries the same
  passage.

## Glossary

None.

## Decisions

- A help flag on a leaf command prints that command's usage; on a group
  prefix with no usage of its own it prints the nearest usage that lists the
  group's commands. No ADR: this applies the conventions above and ADR-0256's
  definition of a Behavior Surface.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the fixes, the transcripts, the
sweep test and the build order.
