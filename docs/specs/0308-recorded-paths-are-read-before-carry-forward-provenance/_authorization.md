---
status: approved
granted: 2026-10-22
action: make the reader of a Task file's Daemon-owned Recorded paths section find it when only Daemon-written settlement sections such as Carry-forward provenance follow it, and make the writer replace it in place, so a carried Task's recorded paths count in the corrective retry and in reconcile while the corrective-scope rule stays unchanged
consuming: 0308-recorded-paths-are-read-before-carry-forward-provenance
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0308

On 2026-10-22, at about 19:50 BRT, the maintainer approved this hotfix through
AskUserQuestion, under the cycle AP authorization. The answers this Spec rests
on were:

- "Hotfix antes da 0300 (Recomendado)": this Spec is delivered before Spec
  0300;
- "Hotfix + Oraculum move o teste (Recomendado)": fix only the defect and
  keep the corrective-scope rule unchanged.

The cycle AP authorization (AskUserQuestion of 2026-10-22, questions of 15:25
BRT answered by about 18:30 BRT) is "the same as cycle AO": the Governed Paths
each Spec declares, and still not this repository's Makefile,
`.github/workflows`, `go.mod`, `go.sum`, or lint and formatter configuration.
That answer carries the cycle AO rules ("the same rules as cycle AN",
2026-10-22, about 00:55 BRT), the cycle AN rules ("the same rules as cycle
AM", 2026-10-21, about 01:00 BRT), the cycle AM rules ("the same rules as
cycle AL", 2026-10-20, about 18:15–18:35 BRT), the cycle AL rules ("Mesmas
regras do AK (Recomendado)", 2026-10-20, about 13:10 BRT), the cycle AK rules
("Sim, mesmas regras do AJ (Recomendado)", 2026-10-20, about 06:30 BRT), and
those carry the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, changes
`.roundfixrc.yml`, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0308`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, the consumer
  sources they read, every planning file of this Spec, its adopted Inbox
  Entry copy, `CONTEXT.md` and `.roundfixrc.yml`.
- Only `.roundfixrc.yml` is governed, and no Task changes it, so `paths:` is
  empty.

## What is not governed

- `internal/spec/recorded_paths.go` and the three test files task_01 changes
  in `internal/spec`, `internal/cli` and `internal/worktree`.
- Every planning file of this Spec and its adopted Inbox Entry copy.

## Sanctioned regeneration

None. No guide, skill, Baseline module, record or mirror changes, so no
regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any guide or skill, or
  the `### QA settlement` section of any skill.
- No change to the corrective-scope rule, Task Carry-Forward, the provenance
  section's bytes, or the readers that classify settlement headings by name.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
