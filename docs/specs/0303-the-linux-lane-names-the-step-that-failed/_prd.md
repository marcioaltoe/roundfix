---
spec: 0303-the-linux-lane-names-the-step-that-failed
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# The Linux Lane names the step that failed

Minimal PRD for a bug fix, minted by `write-techspec`. The Linux Lane selects
the candidate's packages with `go list` before it runs `go test`, but a failed
`go list` ends the lane script with its own exit code, so the Delivery Queue
gate and `roundfix deliver lane` report `go test exited <n>` for a test run
that never started, and a `go list` ended by a signal is reported as the
container's signal. After this Spec the lane script marks a failed `go list`
with its exit code, the outcome names `go list` and that code, and its Park
Class follows the cause, while every `go test` outcome keeps its bytes. The
TechSpec holds the design.

## Sources

One Backlog Entry, written and adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-22-the-linux-lane-reports-a-failed-go-list-as-go-test.md`](references/2026-10-22-the-linux-lane-reports-a-failed-go-list-as-go-test.md),
   type `fix`, from Quaestor's ai-memory report of 2026-10-21 at 19:02 BRT: a
   Codex review on Quaestor pull request #7 of the vendored Roundfix `deliver`
   reference found that a `go list` failure "never reaches `go test` and is
   misclassified". The operator verified it against `linuxLaneScript` and
   `classifyLinuxLaneOutcome` on 2026-10-22.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` held only Backlog Entries about the
documentation contract suite and the QA refusal report, which cycle AP gives
to Specs 0305 and 0302, and `docs/findings/` held no entry.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0), with
Specs 0299 and 0300 active, and delivers after them and after Specs 0301 and
0302 of cycle AP. None of them changes `internal/cli/deliver_linux_lane.go`,
its tests or the Linux Lane sections of the `deliver` documents; they raise
the Roundfix skill's version, so this Spec's documentation Task raises it over
whatever the tree holds when it starts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the lane script gains one output line and one exit code, and the
  outcome text gains one cause. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the lane keeps `--network none`, and every
  test uses temporary repositories and compiled Docker and Go stand-ins.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0316: "The Linux Lane names a
  failed go list apart from go test", which amends ADR-0294: "The Linux Lane
  runs like CI", the record that tells a Lane Container Failure from a
  Lane Change Failure; ADR-0294 amends ADR-0288: "The Delivery Queue gate runs
  the changed Go packages in a Linux Lane" and is amended by ADR-0297: "The
  Linux Lane reads the repository history read-only", whose script lines this
  Spec keeps; ADR-0303: "A Linux Lane skipped for its environment is shown
  where the delivery is read" keeps every skip unchanged. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0303-the-linux-lane-names-the-step-that-failed/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config change.

## Goals

- A non-zero exit of the lane script's `go list` is reported as `go list
  exited <n>`, never as `go test exited <n>`, by the Delivery Queue gate, its
  parked blocker, the lane log and `roundfix deliver lane`.
- A `go list` exit below `128` parks in the `finding` Park Class and makes
  `roundfix deliver lane` exit `1`; one of `128` or more, a signal, parks in
  `environment` and exits `3`.
- Every other lane outcome keeps its bytes, `go test exited <n>` included.
- The `deliver` user guide, both copies of the Roundfix `deliver` reference
  and the glossary describe the `go list` outcome.

## Core Features

1. The script marker. A failed package-selection `go list` prints
   `roundfix-linux-lane: go list: exited <n>` and the lane script exits `4`;
   `go test` still runs only after a successful `go list`.
2. The classification. Exit `4` reads the last well-formed marker line: a code
   below `128` is a Lane Change Failure with cause `go list exited <n>`; `128`
   or more is a Lane Container Failure with cause `go list exited <n>: the
   process ended by signal <n-128>`; no well-formed line is a Lane Change
   Failure with cause `the lane script exited 4 without reporting the go list
   exit code`.
3. Documentation. The `deliver` guide, both reference copies and the
   **Lane Change Failure** glossary entry name the `go list` outcome.

## Non-Goals / Out of Scope

- Changing which packages the lane selects, the container's arguments, mounts
  or environment, or any skip.
- Reclassifying a `go test` failure, including one caused by a module the
  read-only cache lacks, which `go list -e` does not fail for.
- Parsing `go list`'s error text, or running `go list` outside the container.
- Running a real container in any Verification or QA row.

## Success Metrics

1. Success Metric: run under `sh -ceu` with a stand-in `go`, the lane script's
   package-selection tail exits `4` and prints `roundfix-linux-lane: go list:
   exited 1` when `go list` exits `1`, and `exited 137` when it exits `137`,
   without starting `go test`; when `go list` succeeds it prints the package
   count and ends with `go test`'s exit code, `1` for a failing stand-in, and
   prints no `go list` line.
2. Success Metric: with the stand-in Docker exiting `4` after the marker for
   code `1`, the gate's reason is `Linux Lane change failure: go list exited
   1; see <lane log>`, its blocker is `gate-failed: ` followed by that reason,
   the lane log ends with `Linux Lane change failure: go list exited 1`, and
   the parked item is classified `finding`.
3. Success Metric: with the marker for code `137`, the gate's reason is
   `Linux Lane container failure: go list exited 137: the process ended by
   signal 9; see <lane log>`, its blocker starts `linux-lane-container: `, and
   the parked item is classified `environment`.
4. Success Metric: `roundfix deliver lane --all` prints exactly
   `roundfix: Linux Lane change failure: go list exited 1; see <lane log>`
   with empty standard error and exits `1`, and for code `137` prints the
   container line and exits `3`.
5. Success Metric: the classification table proves the boundaries (`127` a
   change failure, `128` a signal), the last marker winning, a marker embedded
   in a line or with a non-decimal code ignored, and every exit other than `4`
   keeping its existing summary, `go test exited 5` included.
6. Success Metric: the guide, both reference copies and the glossary name the
   `go list` cause, the skill version is recorded, and the Behavior Surface
   Record is current.
7. Success Metric: after release, `roundfix deliver lane --all` in a scratch
   clone whose candidate commit holds a `go.mod` that does not parse prints
   `roundfix: Linux Lane change failure: go list exited 1; see <lane log>` and
   exits `1` in the real `golang:1.26` container.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- Quaestor's report of 2026-10-21, recorded in the adopted Backlog Entry
  before this Spec existed: "every started nonzero exit is reported as `go
  test exited <n>`, but package selection runs `go list` first; a `go list`
  failure never reaches `go test` and is misclassified".
- The go command documentation, "List packages or modules"
  (https://pkg.go.dev/cmd/go, also `go help list`): "The -e flag changes the
  handling of erroneous packages, those that cannot be found or are
  malformed. [...] With the -e flag, the list command never prints errors to
  standard error and instead processes the erroneous packages with the usual
  printing."

## Unreachable Acceptance

- criterion: Success Metric 7, after release, `roundfix deliver lane --all`
  in a scratch clone whose candidate commit holds a `go.mod` that does not
  parse prints `roundfix: Linux Lane change failure: go list exited 1; see
  <lane log>` and exits `1` in the real `golang:1.26` container.
  reason: no Verification or QA row may start a container; the lane needs
  Docker and the local `golang:1.26` image, which the Run sandbox does not
  reach.
  satisfied-by: the operator's `roundfix deliver lane --all` in a scratch
  clone of this repository, on a throwaway commit that breaks `go.mod`, after
  this Spec's release.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2`. The current script tail, run under `/bin/sh`
  and `/bin/dash` with a stand-in `go`, exited `1`, `2` and `137` when `go
  list` exited `1`, `2` and `137`, with no line naming `go list`;
  `classifyLinuxLaneOutcome` reads the first two as `go test exited <n>` and
  the third as `docker run exited 137: the container process ended by signal
  9`. The proposed tail exited `4` with the marker in each case, and kept `go
  test`'s exit when `go list` succeeded. `go list -e` with the lane's
  environment exited `1` for a `go.mod` that does not parse, a `go` directive
  above the toolchain and a `go.work` naming a missing module, and `0` for an
  absent required module, an unknown import and a Go syntax error. The
  TechSpec records the commands.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "Linux Lane go
  list failure misclassified as go test" --all --files --min-score 0.3`, which
  returned the Roundfix mirror's Linux Lane ADRs (ADR-0288, ADR-0294,
  ADR-0297), two unrelated ADRs on module versions and the selective gate, the
  Archive Record of Spec 0282 and a Go lint skill; none records
  a decision on the lane's `go list` step, and ADR-0294 is the record this
  Spec amends. No Inbox Entry was
  written, because the briefing forbids writing there.
- Exa: the go command documentation above, fetched from pkg.go.dev, whose
  `-e` passage matches `go help list` of go 1.26.5. Challenging: a module
  absent from the read-only cache could have made `go list` an environment
  failure, but `-e` turns it into an erroneous package and exit `0`.

## Glossary

- changes: **Lane Change Failure**

## Decisions

- 2026-10-22 (about 15:25–18:30 BRT): the cycle AP scope "Correções de
  adotantes (Recomendado)", the order "Correções, follow-up, Python
  (Recomendado)", and the authorization "same as cycle AO"; this Spec is the
  operator's bullet "0303 the Linux Lane names the step that failed".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AP briefing for Spec 0303 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- A failed `go list` has its own reserved exit, `4`, and a marker line with
  its code; the exit decides the step, the line carries the code.
- A `go list` exit below `128` is a Lane Change Failure (`finding`), because
  every measured cause is the candidate's module files; `128` or more is a
  Lane Container Failure (`environment`), as for any process ended by a
  signal.
- One new ADR, ADR-0316, records the rule and amends ADR-0294; the glossary
  revises **Lane Change Failure**.

## Open Questions

None.
