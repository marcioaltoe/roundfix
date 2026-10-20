---
status: approved
granted: 2026-10-20
action: make the Archive Command and the Cancel Command compute every repository-relative path from one Repository Root whatever spelling of the checkout they run from, so the archive rewrites relative citations, relinks Live Documents and names paths in its refusals from a linked worktree, a subdirectory or a path that crosses a symbolic link, and say so in the archive and cancel user guides with the Coverage Review of their surfaces
consuming: 0288-the-archive-rewrites-relative-citations-from-any-worktree
paths:
  - internal/spec/archive.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0288

On 2026-10-20, at about 13:10 BRT, the maintainer approved cycle AL through
AskUserQuestion. The answers this Spec rests on were:

- sequence "Aprovar como está (Recomendado)" (option text: "Duas Specs
  separadas, um minor cada, nesta ordem; P primeiro porque trava o caminho de
  recuperação documentado."): P (0288) then Q (0289), one minor each;
- authorization "Mesmas regras do AK (Recomendado)": the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AK rules that answer carries were approved on 2026-10-20, at about
06:30 BRT ("Sim, mesmas regras do AJ (Recomendado)"), and they carry the cycle
AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and "Merge squash +
release + push".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, changes
`.roundfixrc.yml`, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0288`, based
on `7a678ddb`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference and `CONTEXT.md`.
- Three of those files are governed. One is changed and listed above.
  `internal/spec/archive_test.go` and
  `.agents/skills/roundfix/references/archive.md` are governed too, but the
  Tasks only read them (task_01 runs tests from the first, task_02 reads the
  second for the Coverage Review), so they are not bounded.

## Why the governed path is unavoidable

- `internal/spec/archive.go` holds `Archive`, which computes the Spec folder,
  the Spec Root, the Archive Root and the Archive Record as repository-relative
  paths and hands them to the Archive Relink and the Legacy Link Rebase. It is
  where the defect lives (task_01).

## What is not governed

- The other production and test sources the Tasks change in `internal/spec`
  and `internal/cli`.
- The user guides `docs/user-guide/commands/archive.md` and
  `docs/user-guide/commands/cancel.md`.
- `docs/references/skill-coverage.json` and
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned command re-records the Behavior Surface Record after
task_02 adds a sentence to the `archive` and `cancel` user guides. This
declaration records the regeneration that follows the approved edits and adds
no source path.

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No owned skill or Baseline module changes, so no skill or module version is
recorded, `make skills-sync` and `make baseline-digests` have nothing to
rewrite, and the coverage record is not re-recorded because no test is
removed or renamed and no platform set changes. No record or mirror is
hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module, any skill,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to what the Archive Command or the Cancel Command decide: the QA
  eligibility, the QA Archive Override, the refusals and their order, the
  Archive Record schema and fields, the Archive Relink's candidates, the
  Legacy Link Rebase's rules and every Archive Record already written stay as
  they are. Only the spelling of the paths they compare changes.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
