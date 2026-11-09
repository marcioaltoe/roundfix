---
status: approved
granted: 2026-11-09
action: make doctorDiskSpaceAncestor walk past a file that stands where a directory is expected (ENOTDIR) as it walks past a missing path, and make every disk-full test path that reaches the filesystem lie under the test's own temporary directory
consuming: 0368-the-disk-full-tests-read-no-real-host-path
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

# Approved authority for Spec 0368

On 2026-11-09 the operator recorded the Backlog Entry "The disk-full tests
read a real host path"
(`references/2026-11-09-the-disk-full-tests-read-a-real-host-path.md`,
recorded by the operator), after an empty `/tmp/x` on the maintainer's
MacBook parked Spec 0365's first Run. The maintainer's answer this Spec
rests on (AskUserQuestion via `/roundfix --follow-up`, 2026-11-09 ~11:20
BRT) is, verbatim:

- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto
  fica para amanhã."

"Mesmas do BA" keeps the cycle BA rules of 2026-11-09 ("Mesmas do AZ
(Recommended)" — "Release patch assim que cada Spec ficar pronta, teto de 2
corretivas por Spec, merge pela fila, avisos ao Jarvis, sem notas de release
às caixas."), and through them the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados". This Spec is the item "D testes de disco cheio"
of that question. The design needs no Governed Path and no Makefile,
workflow, Go module, lint or formatter change, although the answer would
allow them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0368`,
based on `cfc08683`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file task_01 declares, `internal/cli/deliver_disk_full.go`,
  every planning file of this Spec, its adopted reference and the removed
  Backlog Entry path. None is governed, so `paths` is empty.

## Sanctioned regeneration

None. No skill, guide, record or generated file changes, and no top-level
test is removed or renamed, so no record is re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide, skill, ADR,
  `CONTEXT.md`, any record, `internal/cli/deliver_disk_full.go` or
  `internal/cli/deliver_free_space.go`.
- No test, Verification command or QA row:
  - creates, removes or changes a file in `/tmp`, `/private/tmp` or `$HOME`
    outside a temporary directory it made;
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
