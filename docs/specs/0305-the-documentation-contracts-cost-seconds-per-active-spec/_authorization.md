---
status: approved
granted: 2026-10-22
action: make the removed-history citation pass of spec.HistoryCitationsInSource look each token up in a set with the same History Citations, and make the corpus tests of internal/docscontract share one speccheck.Check per active Spec per test process with every assertion kept, without changing the Makefile or the CI workflows
consuming: 0305-the-documentation-contracts-cost-seconds-per-active-spec
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0305

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of 15:25 BRT were answered by about 18:30 BRT. The answers this
Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is the documentation contract suite half of
  the second item;
- order: "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The operator's briefing for this Spec says to fix inside Go code only, to keep
every contract assertion, and, where the only fix is running the package once
per CI job, to record that as the Spec's follow-up. The follow-up is the open
Backlog Entry
`docs/backlog/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`.

The cycle AO rules that answer carries were approved on 2026-10-22, at about
00:55 BRT ("the same rules as cycle AN"), which carry the cycle AJ grant of
2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and "Merge squash + release + push". The
cycle AO answer "Backlog por enquanto (Recomendado)" kept the adopted entry in
the backlog until this cycle.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0305`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference and the Backlog Entry path it left, and the new
  follow-up Backlog Entry.
- None of those files is governed, so `paths` is empty.

## What is not governed

- `internal/spec/citation_sources.go` and the new
  `internal/spec/history_citation_cost_test.go` (task_01).
- `internal/docscontract/corpus_test.go` (task_02).
- Every file the Tasks read as an instruction.
- Every planning file of this Spec, its adopted reference and both Backlog
  Entry paths.

## Sanctioned regeneration

None. No skill, guide, Baseline module, record or derived file changes, so no
regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any skill, any guide,
  `CONTEXT.md`, the corpus golden, or the `### QA settlement` section of any
  skill.
- No change to any finding, severity, summary or fix line of `roundfix spec
  check`, or to the Contract Relevance of any test.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
