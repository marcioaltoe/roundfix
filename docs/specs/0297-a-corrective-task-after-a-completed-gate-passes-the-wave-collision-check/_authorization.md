---
status: approved
granted: 2026-10-22
action: make every SC-WAVE-COLLISION message of the Spec Consistency Check and the Daemon name the one needs edge that orders the pair, pin the corrective-Task shape across the QA gate precondition, the strict check, the settlement check and the pre-dispatch refusal, and describe the finding in the spec user guide and the Roundfix spec reference
consuming: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
paths:
  - internal/speccheck/coherence.go
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0297

On 2026-10-22, at about 00:55 BRT, the maintainer approved cycle AO through
AskUserQuestion. The answers this Spec rests on were:

- scope, multi-select, all four in this order: "Z: portão e corretiva",
  "Espera de checks", "Retomar item arquivado", "Y: reconcile PID", one Spec
  and one minor each (v0.108–v0.111); this Spec is "Z: portão e corretiva";
- authorization: the same rules as cycle AN, the Governed Paths each Spec
  declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AN rules that answer carries ("the same rules as cycle AM",
2026-10-21, about 01:00 BRT) carry the cycle AM rules ("the same rules as
cycle AL", 2026-10-20, about 18:15–18:35 BRT), the cycle AL rules ("Mesmas
regras do AK (Recomendado)", 2026-10-20, about 13:10 BRT), the cycle AK rules
("Sim, mesmas regras do AJ (Recomendado)", 2026-10-20, about 06:30 BRT), and
those carry the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push". The 2026-09-30 skills grant ("considere autorizado
a ajustar todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, changes
`.roundfixrc.yml`, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0297`, based
on `34b4e246`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted Inbox Entry copy,
  `internal/spec/collision_test.go`, `internal/speccheck/coherence_test.go`,
  `internal/daemon/task_engine_test.go`, `CONTEXT.md`,
  `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Five of those files are governed. Four are changed and listed above.
  `.roundfixrc.yml` is governed, but no Task changes it, and it is not
  bounded.

## Why each governed path is unavoidable

- `internal/speccheck/coherence.go` holds `detectWaveCollisions`, whose
  summary and fix must name the resolving edge (task_01).
- `.agents/skills/roundfix/references/spec.md` covers `roundfix spec check`,
  whose `SC-WAVE-COLLISION` finding it now describes (task_02).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_02).

## What is not governed

- `internal/spec/collision.go`, `internal/daemon/task_engine.go` and the three
  new test files under `internal/spec`, `internal/speccheck` and
  `internal/daemon`.
- The user guide `docs/user-guide/commands/spec.md`.
- The skill reference mirror `skills/roundfix/references/spec.md`,
  `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Inbox Entry copy.

## Sanctioned regeneration

The repository-owned commands resolve the Behavior Surface Record after
task_02 edits the `spec` guide, and the skill mirror and the skill version
after task_02 edits the Roundfix reference. This declaration records the
regeneration that follows the approved edits and adds no source path.

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change.
No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `CONTEXT.md` or the
  `### QA settlement` section of any skill.
- No change to which pairs collide: `Collisions`, `priorRunTaskTouches`,
  `taskDependencyClosure`, `spec.Load`, `speccheck.Check`, the QA gate
  precondition, the settlement checks, the serial-capacity exemption, the
  Precondition Refusal report, the Prior QA Report import and every file
  under `internal/delivery/` stay as they are.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment in this repository. Verification stays
  Daemon-owned, and Task status stays Daemon-written.
