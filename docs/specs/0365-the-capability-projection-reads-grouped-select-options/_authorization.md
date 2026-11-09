---
status: approved
granted: 2026-11-09
action: make the capability projection read a grouped ACP select option as its groups' values under the bounds of a flat select, pin the two codex-acp payloads recorded by the operator as test fixtures, and prove the Exact Agent Selection Proof over both, as ADR-0375 decides
consuming: 0365-the-capability-projection-reads-grouped-select-options
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0365

On 2026-11-09 the operator recorded that codex-acp 2.6.0 makes every Codex
profile fail its preflight with `capability_evidence_invalid`
(`references/2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md`).
The maintainer's answers this Spec rests on (AskUserQuestion, cycle BA,
2026-11-09 ~08:35–08:45 BRT) are, verbatim:

- Question "Defeito novo e bloqueante: o codex-acp 2.6.0 (08/11, instalado
  hoje pelo dev_upgrade) manda a opção `mode` em grupos ACP; o Roundfix só lê
  opções planas e o preflight recusa TODO perfil codex (doctor: profiles
  failed, capability_evidence_invalid). Afeta todos os adotantes que
  atualizaram. Como trato?" → "Spec C hotfix primeiro (Recommended)" —
  option text: "Spec de correção: a projeção de capacidades lê select options
  agrupadas (achata os grupos), com teste do payload real do 2.6.0. Entra na
  frente da B e da A e sai como patch v0.137.2 logo após o merge."
- Rules: question "Que autorizações você concede de saída para eu levar o
  ciclo BA até o fim sem parar (implementações, merges e releases inclusos
  como no AZ)?" → "Mesmas do AZ (Recommended)" — "Release patch assim que
  cada Spec ficar pronta, teto de 2 corretivas por Spec, merge pela fila,
  avisos ao Jarvis, sem notas de release às caixas."; "Scripts/Makefile/CI" —
  "Mexer em ~/.roundfix-operator, Makefile e .github/workflows se a migração
  ou a Spec B precisarem."; "go.mod/lint" — "Mexer em go.mod/go.sum e na
  configuração de lint/formatter se uma Task precisar."; "Avisar
  Quaestor/Fiscus" — "Mensagem ai-memory curta a cada um dizendo qual Spec
  adotou o finding (não é nota de release)."

"Mesmas do AZ" keeps the cycle AZ rules, and through them the cycle AY, AX,
AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas
as skills se necessário"). The design needs no Makefile, workflow, Go
module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts v0.137.2 right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0365`,
based on `2e59aebe`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted references
  and ADR-0375.
- No path is governed, so `paths` is empty.

## What is not governed

- `internal/agent/selection_capabilities.go`, which task_01 changes.
- The fixtures `internal/agent/testdata/codex-acp-2.6.0-config-options.json`
  and `internal/agent/testdata/codex-acp-2.5.0-config-options.json`, and the
  test files `internal/agent/selection_capabilities_grouped_test.go` and
  `internal/agent/selection_proof_grouped_mode_test.go`, which tasks 01 and
  02 create.
- ADR-0375, every planning file of this Spec and its adopted references.
