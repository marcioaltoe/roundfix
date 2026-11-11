---
status: approved
granted: 2026-11-10
action: make the suite-time budget of a push to main one measured number, 420 s in the Makefile, have make test-budget record its verdict in the file SUITE_BUDGET_VERDICT names and a new make suite-budget-check judge it, wire the CI Verify workflow so the push Verify step records the verdict and a last step judges it after Verify docs and Verify every contract, add the test that reads both files, write the Amended by ADR-0387 notice into ADR-0259, and regenerate the setup manifest
consuming: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
paths:
  - .github/workflows/ci-verify.yml
  - Makefile
  - docs/agents/setup-context.json
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

# Approved authority for Spec 0377

On 2026-11-10 the operator recorded that the `CI | Verify` gate of a push to
`main` fails its 240 s suite-time budget on most runs while every test
passes, and skips `Verify docs` and the Full Contract Run
(`references/2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md`,
recorded by the operator). The maintainer's answers this Spec rests on
(AskUserQuestion via `/roundfix --follow-up`, 2026-11-10 ~20:55–21:00 BRT)
are, verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among others,
  "Portão do CI (Recommended)" — option text: "O push no main falhou 9 de 12
  vezes só pelo orçamento de 240 s da suíte (os testes passam) e pula o Full
  Contract Run. Spec pequena no workflow de CI."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

"CI na 0377 … só nos arquivos que elas exigirem" grants
`.github/workflows/ci-verify.yml`, the file that sets the push budget and
orders the steps. That answer granted the `Makefile` to 0378 only, so the
operator asked the maintainer again (AskUserQuestion, 2026-11-10 ~22:35 BRT):
question "O plano da Spec 0377 (portão do CI) precisa mexer no `Makefile`,
além do workflow de CI. Só o alvo `test-budget` sabe se a suíte falhou ou só
demorou: ele passa a gravar o veredito do orçamento num arquivo, o orçamento
vira 420 s, e um alvo novo julga esse veredito no fim do CI. Sua autorização
das 21:00 liberou \"CI na 0377 e Makefile na 0378\". Libero o `Makefile` (e o
`docs/agents/setup-context.json` regenerado) também na 0377?" → "Liberar na
0377 (Recommended)" — option text: "A 0377 muda o `Makefile` e o workflow,
como o plano prevê. Um push com testes verdes passa a rodar sempre o Verify
docs e o Full Contract Run, e um estouro real de tempo continua reprovando."
This answer grants the `Makefile` and the regenerated
`docs/agents/setup-context.json` to 0377. The design needs them: only the
`test-budget` target that runs the suite knows whether the suite failed or
only took too long, so the deferred verdict and the one budget number live
there. `docs/agents/setup-context.json` records the `Makefile` digest and
changes only through the regeneration below (lessons 42 and 55 of the
operator's brief). "Mesmas do BB" keeps the cycle BB rules ("Mesmas do BA,
teto minor (Recommended)") and, through them, the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados".

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0377`,
based on `60ac8142`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its references and the
  adopted entry's old path.
- The three paths above are governed: the workflow, the `Makefile` and the
  setup manifest task_01 regenerates. The new test file, the ADR files and
  every planning file are not governed; the `docs/agents/` guides task_01
  reads (`agent-instructions.md` and `spec-routing.md` are governed) are read
  only.

## Sanctioned regeneration

The repository-owned command writes the setup manifest after the approved
`Makefile` edit. This declaration records that regeneration and adds no
source path.

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

The digest is never hand-edited. No top-level test is removed or renamed, so
the coverage record is not re-recorded.

## Limits

- No change to `.github/workflows/release.yml` or any other workflow, any
  `Makefile` target other than `test-budget`, the new `suite-budget-check` and
  the `.PHONY` line, `go.mod`, `go.sum`, the lint, formatter or test-runner
  configuration, `.roundfixrc.yml`, any Baseline module, profile, template or
  derived Baseline file other than the setup manifest, any guide under
  `docs/agents/`, `CONTEXT.md`, any skill, any Go production code, any ADR
  other than the notice in ADR-0259, any Archive Record, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - starts, dispatches, reruns or cancels a GitHub workflow, or pushes;
  - reads or writes the real `~/.roundfix`;
  - creates files in `/tmp`, `/private/tmp` or `$HOME` outside its own
    temporary directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network or runs `gh`;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
