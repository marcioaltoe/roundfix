---
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# The review drops a Result Status Note, and a stale fingerprint names its command

Minimal PRD for two bug fixes, minted by `write-techspec`. The pre-PR review
parks a delivery on a finding that a Task's `## Result` reports the status
`in_progress` while the Daemon settled the Task `completed`, although
Delivery Convention C2 describes that pair as designed, because the sealed
validator decides it and twice read C2 as a requirement the Result failed.
After this Spec C2 says who writes the Result and when, and Roundfix dismisses
that finding by validation itself, with a recorded reason, while every real
contradiction still reaches the validator. Separately, a stale Behavior
Surface fingerprint is reported as `<id>: stale fingerprint` alone; after this
Spec the problem names the command that re-records it. The TechSpec holds the
design.

## Sources

Two Backlog Entries written into `references/` at birth and adopted by this
Spec (see `references/_index.md`):

1. [`references/2026-10-22-the-pre-pr-review-parks-on-a-result-status-note-that-convention-c2-allows.md`](references/2026-10-22-the-pre-pr-review-parks-on-a-result-status-note-that-convention-c2-allows.md),
   type `fix`: the reviewer raised, and the validator kept as `stands`, the
   Result-versus-front-matter status pair on Spec 0290 (review head
   `61dec34e`) and Spec 0297 (review head `1aa6ae3a`, 2026-10-22); the
   operator dismissed both with evidence citing C2.
2. [`references/2026-10-22-a-stale-behavior-surface-fingerprint-does-not-name-its-re-record-command.md`](references/2026-10-22-a-stale-behavior-surface-fingerprint-does-not-name-its-re-record-command.md),
   type `fix`: `TestTheSkillCoverageMapIsCurrent` reports only `<id>: stale
   fingerprint`; on 2026-10-22 a Spec 0298 item parked `gate-failed` and a Spec
   0299 Run ended Unresolved for a stale fingerprint after a guide edit.

At the authoring base `docs/backlog/` held one open entry, the QA refusal
report's audited commit, which Spec 0302 of this cycle adopts; `docs/findings/`
held none. Neither shares this context.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0), with
Specs 0299 and 0300 active, and delivers after them and after the cycle AP
Specs 0301 to 0303. None of them changes the review conventions, the convention
validator, the review guide or `skill_coverage_test.go`; several raise the
Roundfix skill's version, so this Spec's first Task raises it over whatever the
tree holds when it starts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the review record keeps its fields and the existing rule value
  `convention:C2`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; every test uses temporary repositories and fake runners,
  and no provider is called. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0317: "A Result Status Note of a
  completed Task is dismissed without the validator", which refines the
  validation of pre-PR review findings and rests on the Daemon's ownership of
  Task status; both earlier decisions stay unchanged and are cited in
  ADR-0317's own text. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config change.

## Goals

- A pre-PR review finding that names a Result Status Note against a
  `completed` non-QA Task is dismissed by validation with rule
  `convention:C2` and a fixed reason, without the sealed validator, and a
  review whose findings are all dismissed reads `findings-dismissed`.
- A Task whose status is not `completed`, a QA Task, a Result without a
  Result Status Note, an anchor outside the C2 region and a finding that names
  another failure keep reaching the validator or standing as today.
- C2 states that the Agent writes the Result while the Daemon holds the Task
  `in_progress`, and the review guide and the Roundfix review reference carry
  the new text and the dismissal.
- The stale fingerprint problem names the command `.roundfixrc.yml` declares
  for the Behavior Surface Record; every other problem keeps its bytes.

## Core Features

1. C2 version 4. Delivery Conventions become
   `roundfix/delivery-conventions/v4`; C2 no longer says the Daemon writes the
   Result and says that a Result reporting `in_progress`, unchanged or
   Daemon-owned does not contradict `completed`.
2. The Result Status Note dismissal. Before asking the validator, Roundfix
   dismisses a finding whose anchor lies in a Task file's C2 region, whose
   Task is `completed` and not QA, whose Result holds a Result Status Note and
   whose text names the pair. With nothing left to ask, the review settles on
   the findings it has.
3. Documentation. The review guide and both copies of the Roundfix review
   reference describe C2 v4 and the dismissal; the glossary adds
   **Result Status Note** and revises **Delivery Convention**.
4. The named re-record. `TestTheSkillCoverageMapIsCurrent` reports a stale
   fingerprint with the re-record command, and a contract test holds that
   command equal to the `.roundfixrc.yml` declaration.

## Non-Goals / Out of Scope

- Changing the Agent prompt or the `implement-task` skill so Results stop
  reporting the status: Results already written keep the sentence, and the
  dismissal covers both.
- Changing C1, C3, C4 or C5, the anchor validation, the validator's prompt or
  grammar, or how the operator disposes a finding.
- Changing `.roundfixrc.yml`, the Behavior Surface Record's format or the
  other problems `TestTheSkillCoverageMapIsCurrent` reports.

## Success Metrics

1. Success Metric: in a temporary repository, a review whose only finding has
   the shape of the Spec 0297 finding, the Spec 0290 finding or a front-matter
   anchor naming the same pair exits 0 with outcome `findings-dismissed`,
   `validator` `not-needed`, no sealed request, rule `convention:C2`, the
   fixed reason, and one stderr line naming the finding, rule and reason.
2. Success Metric: with the same Result, a `failed` Task, a QA Task, a Result
   whose note was removed, a note only outside the Result and a finding that
   names another failure each send exactly one question with
   `convention:C2` to the fake validator, and the finding stands with its
   verdict; an anchor outside the C2 region stands without a question.
3. Success Metric: each removed guard of the dismissal (completed status, QA
   type, Result section, Result Status Note, finding text, C2 region, settling
   without a question) makes a named test case fail.
4. Success Metric: the review prompt and the validator input carry
   `roundfix/delivery-conventions/v4` and the new C2 text, and the guide and
   both reference copies carry the same text.
5. Success Metric: a stale fingerprint problem reads `<id>: stale
   fingerprint; re-record docs/references/behavior-surfaces.json after the
   last edit of its sources with: <command>`, where `<command>` equals the
   `regenerate` `.roundfixrc.yml` declares for that record; the missing and
   extra problems keep their bytes.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- The Task files the real deliveries wrote: Spec 0297's `task_02.md` and Spec
  0290's `task_01.md` at the `source_revision` their Archive Records name
  (`docs/history/specs/0297-*.md`, `docs/history/specs/0290-*.md`). Their
  Results hold "Task status remains `in_progress`" (0297, lines 127-129) and
  "leaving the Daemon-owned status `in_progress` unchanged" (0290, lines
  116-117) under a `status: completed` front matter, the pair the review kept.
- Anthropic's cookbook "Building Evals"
  (https://github.com/anthropics/anthropic-cookbook/blob/main/misc/building_evals.ipynb),
  on code-based grading: "This is by far the best grading method if you can
  design an eval that allows for it, as it is super fast and highly
  reliable."

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2` with `go test -overlay` probes that wrote
  nothing to the repository. Through `roundfix review` with a fake reviewer
  and a fake sealed validator, a finding of the 0297 shape was sent to the
  validator with the eligible rules `[convention:C2]`, the validator input
  carried C2's v3 sentence "A Result that calls status Daemon-owned agrees
  with a `completed` status.", and the `stands` verdict of the real run left
  the review at exit 1, outcome `findings`. A prototype of this design, in the
  same overlay form, passed the new tests and failed a named case for each of
  seven sabotages of the dismissal and three of the stale fingerprint
  problem; the TechSpec records them.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "pre-PR review
  validator dismisses Task Result status Daemon-owned in_progress completed
  convention" --all --files --min-score 0.3`, which returned the Roundfix
  mirror's decisions that the Daemon sets `in_progress` and writes the
  terminal status, and that the review validator fails closed, and `qmd query "LLM
  judge false positive deterministic rule before model validator fail
  closed" --all --files --min-score 0.3`, which returned the concept
  `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` and eval
  sources. No prior decision on this pair exists. No Inbox Entry was written,
  because the briefing forbids writing there.
- Exa: the Anthropic cookbook passage above, which favors a code check over a
  model grader wherever the case allows one. Challenging: the same source
  says many evaluations do not allow code grading, which is why every finding
  outside the narrow pair still goes to the validator.

## Glossary

- adds: **Result Status Note**
- changes: **Delivery Convention**

## Decisions

- 2026-10-22 (questions of about 15:25 BRT, answered by about 18:30 BRT): the
  cycle AP scope "Correções de adotantes (Recomendado)", "Ruído C2 e suite
  docscontract", "Modos follow-up A e B", "Preset de baseline Python"; the
  order "Correções, follow-up, Python (Recomendado)"; and the authorization
  same as cycle AO.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados".

Autonomous route: no person present; decisions taken from the operator's cycle AP briefing for Spec 0304 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The fix is both a C2 rewording and a mechanical dismissal; the rewording
  alone leaves the verdict to the model.
- The dismissal reads the finding's text as well as the head tree, so a
  finding anchored on the Result that names another failure still reaches the
  validator.
- The dismissed finding keeps the existing rule `convention:C2`, so every
  reader of the rule, `review dispose` among them, keeps its behavior.
- The re-record command is held equal to the `.roundfixrc.yml` declaration
  the Delivery Queue's derived merge runs, so the operator, the Agent and the
  queue use one command.
- One new ADR, ADR-0317, records the dismissal; the glossary adds **Result
  Status Note** and revises **Delivery Convention**.

## Open Questions

None.
