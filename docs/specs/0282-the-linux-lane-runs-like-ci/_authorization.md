---
status: approved
granted: 2026-10-19
action: run the Linux Lane under an init on the image delivery.linux_lane_image names; make a test that needs an External Test Tool skip naming it and refuse a new literal lookup of one; tell a Lane Container Failure from a Lane Change Failure in the lane log, the blocker and the Park Class; add roundfix deliver lane; and describe each change in the user guides, the Roundfix deliver reference, the repository rules and the glossary
consuming: 0282-the-linux-lane-runs-like-ci
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - docs/agents/specific-repository.md
  - internal/cli/baseline_release_gate_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0282

On 2026-10-19, at about 16:20 BRT, the maintainer approved cycle AJ through
AskUserQuestion. The answers were:

- decomposition "Três Specs (Recomendado)": J (0281), then K (0282, this
  Spec), then L (0283), one minor release each;
- this Spec's design "--init + imagem configurável + skip com motivo
  (Recomendado)" (option text: "O lane roda com --init; a imagem passa a ser
  configurável em .roundfixrc; testes que precisam de ferramenta externa
  pulam com motivo quando ela falta; falha de container é distinguida de
  falha da mudança. Sem build de imagem própria.");
- the authorizations, all four selected: "Religar delivery.linux_lane: true",
  "docker pull/build da imagem", "Mexer em Governed Paths declarados" (option
  text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).") and "Merge squash + release + push";
- costs "Regras do ciclo AH/AI (Recomendado)", failures "Regras AE–AH e parar
  na proposta (Recomendado)", and authoring "Opus + strict check, enfileirar
  sem esperar você (Recomendado)".

The 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. "Religar
delivery.linux_lane: true" is the operator's own change after this Spec
merges and an on-demand lane run passes; no Task of this Spec changes
`.roundfixrc.yml`. "docker pull/build da imagem" covered the single `docker
pull golang:1.26` and the read-only probes made while this plan was authored;
no Task, Verification command or QA row pulls or builds an image.

The governed set was measured on the authoring branch `docs/plan-0282`, based
on `f382480b`, with `speccheck.GovernedPath`.
- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status` was identical before and after, and
  the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  the adopted Backlog Entry at both paths, the two ADRs the plan writes and
  `.roundfixrc.yml`.
- Six of those files are governed. Five are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_02, task_03,
  task_04).
- `.agents/skills/roundfix/references/deliver.md` covers `deliver`, whose
  Linux Lane gains `--init` and the configured image (task_02), the two
  failure outcomes and the new blocker (task_03), and the `deliver lane`
  subcommand (task_04). It also covers the `config keys` surface that the new
  key moves.
- `docs/agents/specific-repository.md` holds the repository rules; it gains
  the External Test Tool rule outside any setup marker (task_01).
- `internal/cli/baseline_release_gate_test.go` holds the rust-cli journey,
  which resolves `rustfmt` through `testfixture.RequireTool` (task_01).

## What is not governed

- The production and test sources the Tasks change or create in
  `internal/testfixture`, `internal/config`, `internal/cli` (except the file
  above), `internal/delivery` and `internal/agent`, including every new file.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json`,
  `docs/references/behavior-surfaces.json` and
  `docs/references/skill-coverage.json`.
- `CONTEXT.md`, the user guides under `docs/user-guide/` and the ADRs.
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
platform set changes, so the coverage record is not re-recorded. No record or
mirror is hand-edited; the new `command deliver lane` entry of
`docs/references/skill-coverage.json` is authored by task_04, because the map
is a reviewed source, not a generated record.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to the lane's skips, its package selection, its network and pull
  policy, the variables the container receives, the module cache mount, the
  cap's length, any other blocker or Park Class, or the repository
  Verification's park.
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
