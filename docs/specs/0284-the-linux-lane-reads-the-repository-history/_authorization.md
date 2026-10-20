---
status: approved
granted: 2026-10-19
action: mount the repository object store read-only in the Linux Lane and register it as an alternate of the container repository, skip a store the lane cannot mount whole with its reason, and describe the change in the deliver user guide, the Roundfix deliver reference and the glossary
consuming: 0284-the-linux-lane-reads-the-repository-history
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0284

On 2026-10-19, at about 16:20 BRT, the maintainer approved cycle AJ through
AskUserQuestion. The answers this Spec rests on were:

- the authorizations, all four selected: "Religar delivery.linux_lane: true",
  "docker pull/build da imagem", "Mexer em Governed Paths declarados" (option
  text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).") and "Merge squash + release + push";
- costs "Regras do ciclo AH/AI (Recomendado)", failures "Regras AE–AH e parar
  na proposta (Recomendado)", and authoring "Opus + strict check, enfileirar
  sem esperar você (Recomendado)".

At about 22:00 BRT the operator added this Spec to the cycle as a patch for a
defect of the cycle's own release: Spec 0282's unproven on-demand lane run
failed on `main` at `dc7e5f72`. The operator placed it within the cycle AJ
authorizations above (Docker, Governed Paths, merge). The 2026-09-30 skills
grant ("considere autorizado a ajustar todas as skills se necessário") also
stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. "Religar
delivery.linux_lane: true" is the operator's own change after this Spec merges
and an on-demand lane run passes; no Task of this Spec changes
`.roundfixrc.yml`. "docker pull/build da imagem" covered the read-only
container probes made while this plan was authored, in the locally present
`golang:1.26` image; no Task, Verification command or QA row pulls or builds
an image.

The governed set was measured on the authoring branch `docs/plan-0284`, based
on `dc7e5f72`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before and
  after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  the adopted Backlog Entry at both paths, the three ADRs the plan writes or
  amends, `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Four of those files are governed. Three are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).
- `.agents/skills/roundfix/references/deliver.md` covers `deliver`, whose
  Linux Lane gains the history mount and three skip reasons (task_01).

## What is not governed

- The production and test sources task_01 changes in `internal/cli`.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guide `docs/user-guide/commands/deliver.md` and the
  ADRs.
- Every planning file of this Spec, including its adopted reference.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill version and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change. No
record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to the lane's image, init, module cache mount, package selection,
  network and pull policy, the variables the container receives, the cap,
  the outcome classification, any blocker or Park Class, or the repository
  Verification's park.
- The object store is mounted read-only and the working tree is never
  mounted.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network or pulls or builds an image;
  - needs Docker, except the one QA row that runs the real lane only when the
    Run sandbox reaches a Docker daemon and is declared unreachable otherwise;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
