---
status: approved
granted: 2026-10-19
action: read the fetched default branch for Git Locator reachability, the QA Delivery Base, the Archive Command's source_revision warning and the pre-PR review's default base; measure the storage report in one Storage Snapshot that groups pointer-map and lock-byte pages; and describe each change in the user guides, the Roundfix skill references and the glossary
consuming: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/storage.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0283

On 2026-10-19, at about 16:20 BRT, the maintainer approved cycle AJ through
AskUserQuestion. The answers were:

- decomposition "Três Specs (Recomendado)": J (0281) → K (0282) → L (0283), one
  minor each;
- authorizations, all four selected: "Religar delivery.linux_lane: true",
  "docker pull/build da imagem", "Mexer em Governed Paths declarados" (option
  text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).") and "Merge squash + release + push";
- costs "Regras do ciclo AH/AI (Recomendado)"; failures "Regras AE–AH e parar
  na proposta (Recomendado)"; authoring "Opus + strict check, enfileirar sem
  esperar você (Recomendado)".

The 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
se necessário") also stands. This record grants `implement`, `commit`, `push`,
`pull_request` and `merge`; the operator releases, so `release` is not
granted here. This Spec uses neither the Linux Lane switch nor Docker.

The governed set was measured on the authoring branch `docs/plan-0283`, based
on `f382480b`, with `speccheck.GovernedPath`.
- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status` was identical before and after, and
  the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and the ADRs the plan writes.
- Five of those files are governed, all changed and listed above.
  `.github/workflows/ci-verify.yml` is governed too; the TechSpec cites it as
  read-only evidence of the CI checkout, and no Task declares or changes it.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_02,
  task_03).
- Three Roundfix references are governed:
  - `.agents/skills/roundfix/references/archive.md` covers `archive`, whose
    Git Locator reads the Default Branch Head (task_01) and whose
    `source_revision` warning checks against it (task_02);
  - `.agents/skills/roundfix/references/review.md` covers `review`, whose
    default base becomes the Default Branch Head (task_02);
  - `.agents/skills/roundfix/references/storage.md` covers `storage report`
    and `gc`, whose report measures one Storage Snapshot (task_03).

## What is not governed

- The production and test sources the Tasks change or create in
  `internal/preflight`, `internal/spec`, `internal/speccheck`,
  `internal/daemon`, `internal/cli` and `internal/store`, including the new
  `internal/preflight/default_branch_head.go`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides under `docs/user-guide/commands/` and the
  ADRs.
- Every planning file of this Spec, including its adopted references.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill version
and the Behavior Surface Record. This declaration records the regeneration
that follows the approved edits and adds no source path.

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
baseline-digests` has nothing to rewrite. No test is removed or renamed and
no platform set changes, so the coverage record is not re-recorded. No record
or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, or the `### QA
  settlement` section of any skill.
- No change to `internal/cli/implement.go`, `internal/worktree/worktree.go`,
  `internal/cli/deliver.go`, `internal/cli/deliver_workflow.go` or
  `internal/cli/deliver_linux_lane.go`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network or runs Docker;
  - fetches from or pushes to any remote other than a temporary local
    repository or this checkout's own objects;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
