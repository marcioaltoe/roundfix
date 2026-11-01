---
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
status: active
created: 2026-11-01
surfaces: [backend, docs]
---

# The Spec Check leaves the QA gate unprobed and prints its expected Verification

Minimal PRD for a bug fix, minted by `write-techspec`. Two adopters were
blocked by the QA gate's derived Verification on 2026-11-01. Fluxus kept a
failed QA Report, added an honest corrective Task, and
`roundfix spec check --strict --run-verification` still exited 1, because the
probe ran the gate's derived command, which passes on any readable report,
although the Daemon never probes the gate. Quaestor wrote its QA Task from the
`write-tasks` template, got `SC-QA-VERIFICATION-AUTHORED`, and found no way to
learn the text to write short of reading Roundfix's source. After this Spec the
check reports the QA gate's command not probed with reason `qa_gate`, through
the predicate the Daemon uses to route the gate, and the finding prints the
exact `## Verification` block to copy, in text and in JSON, which the
`write-tasks` skill teaches.

## Sources

Two adopter reports, each adopted as a Backlog-style reference by this Spec's
planning change (see `references/_index.md`):

1. [`references/2026-11-01-the-spec-check-probe-judges-the-qa-gates-derived-verification-vacuous.md`](references/2026-11-01-the-spec-check-probe-judges-the-qa-gates-derived-verification-vacuous.md),
   `type: fix`, Fluxus, ai-memory messages `01a1987b` and `01a1987c`.
2. [`references/2026-11-01-no-authoring-path-shows-the-qa-gates-derived-verification.md`](references/2026-11-01-no-authoring-path-shows-the-qa-gates-derived-verification.md),
   `type: fix`, Quaestor, ai-memory message `01a19861`.

Both concern the one command Roundfix derives for every QA gate, so one Spec
owns them. At the authoring base `docs/backlog/` holds no entry and
`docs/findings/` no unresolved Finding on the Spec Check probe or the QA
Verification.

## Prerequisites

None in `requires`. Authored against `340c9004` (v0.127.0), where no other
Spec is active.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the finding code keeps its bytes and the JSON gains one optional
  field inside the existing `roundfix-speccheck/v1` schema. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0148 makes one prober serve the
  Daemon and the authoring check, ADR-0148: "both callers probe the same
  Tasks"; this Spec keeps that rule and amends ADR-0148 to state that
  neither probes the QA gate. ADR-0315 cites ADR-0148 for the shared prober
  and governs pytest selectors at authoring; this Spec does not touch it. No
  new ADR is written. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AT answers of 2026-11-01 quoted in the authorization record, which keep the
  cycle AS rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the two skills' canonical files and
  `SKILL.md` pairs and the Spec Check's coherence detectors. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `internal/speccheck/coherence.go`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- `spec check --run-verification` never runs the QA gate's derived command,
  whatever reports the Spec holds, and reports it not probed with reason
  `qa_gate`, so an honest corrective Task beside a kept failed QA Report
  passes the check.
- The check and the Daemon decide which Task is the QA gate through one
  predicate.
- Every other Task is probed as today: a vacuous, unknown or malformed
  command of a pending, in-progress or failed non-QA Task still exits 1.
- `SC-QA-VERIFICATION-AUTHORED` prints the exact `## Verification` block the
  QA Task must carry, in text and in JSON, and copying it clears the finding.
- The `write-tasks` template and skill, the Roundfix skill's spec reference
  and the spec command guide teach both rules.

## Core Features

1. One QA gate predicate. The Daemon package exports the rule that a Task of
   type `qa` is the QA gate; the Daemon's routing and the Spec Check read it.
2. The probe leaves the gate alone. A pending, in-progress or failed QA gate
   is reported as `- <task>: not probed — "<command>" (QA gate)`, the clean
   verdict line adds `; not probed (QA gate): <n>`, and JSON lists it under
   `verification.notProbed` with reason `qa_gate`. A completed QA gate keeps
   reason `completed`.
3. The expected block. The finding gains an `expected` text: `## Verification`,
   a blank line and one bullet holding the derived command in one code span.
   Text output prints it after the `fix:` line under `  expected block:` as a
   fenced `markdown` block; JSON carries it in the finding's `expected` field,
   omitted on every other finding.
4. The rules in writing. The `write-tasks` template and `SKILL.md`, the
   Roundfix skill's spec reference, the spec command guide, the glossary entry
   for **Vacuous Verification** and an ADR-0148 amendment state them.

## Non-Goals / Out of Scope

- Changing `spec.DerivedQAVerification`, `spec.QAReportEligibility`, the
  audited head or the Daemon's QA step; acceptance of the gate stays where it
  is.
- Letting an author write a stronger QA Verification; a differing command is
  still refused.
- A new command or flag to print the derived block.
- The derived command's fixed `docs/specs/<slug>/qa` path under a non-default
  Spec Root (`_techspec.md` → Risks).
- Any text inside a `### QA settlement` section; the Makefile, the CI
  workflows, `go.mod`, `go.sum`; any Baseline module or guide under
  `docs/agents/`.

## Success Metrics

1. Success Metric: in a temporary repository with a completed Task, a pending
   corrective Task whose command fails, a pending QA gate and a committed
   `qa/qa-report-2026-10-31.md` with `verdict: fail`,
   `spec check --run-verification` exits 0, reports the corrective Task
   `honest` and the gate `not probed` with reason `qa_gate`; at `340c9004`
   it reports the gate `vacuous` and exits 1.
2. Success Metric: the same graph with no QA Report also reports the gate
   `not probed` with reason `qa_gate` and exits 0, and a vacuous pending
   non-QA Task beside a pending QA gate still exits 1.
3. Success Metric: the Daemon's `taskPlanWithoutQAGate` and the Spec Check
   call the same exported predicate, and the Daemon's QA routing tests stay
   green.
4. Success Metric: an authored QA Verification still gives exactly one
   `SC-QA-VERIFICATION-AUTHORED`; its `expected` text equals
   `## Verification`, a blank line and the bullet of the derived command; and
   a QA Task whose section is replaced by that text gives none.
5. Success Metric: the text report prints the block after the `fix:` line
   under `  expected block:` in a fenced `markdown` block, the JSON finding
   carries it in `expected`, and a finding without one has neither the line
   nor the key.
6. Success Metric: the spec command guide and both copies of the Roundfix
   spec reference document the `qa_gate` reason and the `expected` field, and
   both copies of the `write-tasks` template and `SKILL.md` teach copying the
   block.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The two adopter reports quoted in the adopted references, written by the
  Fluxus and Quaestor projects before this Spec existed.
- The rustc book, "JSON Output" (https://doc.rust-lang.org/rustc/json.html,
  read through Exa on 2026-11-01): a diagnostic span carries
  `suggested_replacement`, "An optional string of a suggested replacement for
  this span to solve the issue. Tools may try to replace the contents of the
  span with this text." The `expected` field follows that published shape: an
  optional, exact replacement beside the human message.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-01 against these sources.

- Characterization at `340c9004` with that commit's `bin/roundfix` in a
  scratch clone, recorded in the adopted references: the Fluxus shape exits 1
  with the gate `vacuous`, and an authored QA Verification gives
  `SC-QA-VERIFICATION-AUTHORED` with no expected text.
- Reading at `340c9004`: `Engine.TaskCycle` calls `taskPlanWithoutQAGate` before
  the scheduler, and `verifyTaskPreWork` is called only from
  `executeTaskWorker`, so the Daemon never probes the gate; `runQAGate` calls
  no probe.
- A prototype of both Tasks in a second clone passed the tests
  `_techspec.md` → Testing Approach names, the `./internal/speccheck`,
  `./internal/daemon`, `./skills` and `-tags docscontract
  ./internal/docscontract` suites, and the CLI Spec Check tests; each sabotage
  that section lists failed its named test.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix spec check QA task derived verification vacuous" --all --files --min-score 0.3`.
  It returned only project mirrors (Roundfix's archived Specs 0046 and 0298,
  `docs/user-guide/configuration.md`, and the adopters' `docs/agents/spec-routing.md`);
  none records a decision on probing the QA gate or on printing the derived
  block.
- Exa: the rustc JSON diagnostics page quoted above, which supports carrying
  the exact replacement as its own optional field rather than inside the fix
  prose.

## Glossary

- changes: **Vacuous Verification** — the QA gate is never probed either

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AT brief for Spec 0337 of 2026-11-01 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-01), verbatim:

- Scope: "0337 primeiro, release já" — option text: "Só a 0337, que
  desbloqueia Fluxus e Quaestor, com release logo depois; as outras três num
  ciclo seguinte."
- Derived QA Verification: "spec check mostra (Recommended)" — option text:
  "O achado SC-QA-VERIFICATION-AUTHORED passa a imprimir o bloco exato
  esperado, e o template do write-tasks ensina a copiar dele. Nenhum comando
  novo. Junto: o spec check --run-verification deixa de medir a Verification
  derivada da QA (correção do Fluxus)."
- Authorization: "Mesmas do AS (Recommended)" — option text: "Push e PR do
  plano, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, release v0.128.0 logo
  após o merge e aviso aos adotantes (Fluxus e Quaestor, com resposta às
  mensagens deles). As outras três mensagens viram Backlog Entries com as suas
  decisões, para o próximo ciclo."

Taken from those answers, ADR-0148 and the characterization, with no further
interview:

- The QA gate is never probed, whatever its reports: a pending gate with no
  report used to be probed and judged `honest`, which proved nothing, since
  the command is derived and the Daemon never runs a pre-work probe on it.
  One rule, the Daemon's, replaces the special case.
- The predicate lives in the Daemon package beside `TaskAwaitsWork`, the
  rule ADR-0148 already shares with the check, rather than a second copy of
  the type test in the CLI.
- The expected block is a new optional field on every finding, filled only by
  `SC-QA-VERIFICATION-AUTHORED`, rather than a multi-line fix, so the one-line
  `fix:` and every other finding keep their bytes.
- Two implementation Tasks, one per source, in series because both raise the
  Roundfix skill version and edit the same guide, and a QA gate.

## Open Questions

None.
