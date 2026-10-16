---
status: approved
granted: 2026-10-16
action: record how the delivery item binary probe ended in its tests, diagnose the F6 signal through up to ten full suites and the macOS diagnostic report, and make every delivery test replace its item binary by copy and rename instead of writing into it, with a test that rebuilds while a copy of the previous build still runs
consuming: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0275

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion: scope
"A: history só com records, B: ergonomia da fila, C: guardas de autoria, D:
remetente do F6", order "B → A → C → D (Recomendado)", authorizations
"Governed Paths declarados" (the option text: "\"Concedo\" como em AE/AF/AG
para os Governed Paths que cada Spec declarar: módulos da Baseline, registros
de versão, make baseline-digests, skills e testes governados. Makefile,
workflows de CI, go.mod e configuração de lint ficam fora."), "Merge e
releases", "Merge da PR #545", "Avisos por push"; failures "Regras de AE–AG
(Recomendado)". The 2026-09-30 skills grant ("considere autorizado a ajustar
todas as skills se necessário") also stands. For this Spec, item D, the
maintainer chose: "Instrumentar sem sudo (Recomendado) — A Spec registra pid,
pai e sinal do filho do migrate --check pelo próprio teste e roda as suítes
até reproduzir; só se isso não achar o remetente eu te mando por push o
comando eslogger exato e paro a D."

The governed set was measured on the authoring branch `docs/plan-0272` at
`9e723b82` with `speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository (`git status` showed
no change of the probe's own), against every file the Tasks declare and every
planning file. None of them is governed, so `paths` is empty.

## Why each governed path is unavoidable

No governed path changes.

## What is not governed

`internal/cli/deliver_item_binary_test.go`,
`internal/cli/deliver_item_binary_signal_unix_test.go`, the new
`internal/cli/deliver_item_probe_report_test.go` and
`internal/cli/deliver_item_binary_rebuild_test.go`, the diagnosis record
`docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/diagnosis/f6-sender.md`,
the ADR
`docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`,
and this Spec's planning files and adopted Backlog Entry.

## Sanctioned regeneration

None. No skill, Baseline module, version record, coverage record or digest
changes, so no regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any skill, any Baseline module, any
  production Go file or the `### QA settlement` section of any skill.
- No `sudo`, `eslogger`, `dtrace` or `dtruss`; no setting of the system's code
  signing, Gatekeeper or developer mode.
- No test or Verification command reads `~/Library/Logs`, `~/.roundfix` or
  another repository, or reaches the network; only the diagnosis and the QA
  gate read the diagnostic reports, in their Agent turns.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
