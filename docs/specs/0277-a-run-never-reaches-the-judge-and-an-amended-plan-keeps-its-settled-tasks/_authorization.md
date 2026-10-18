---
status: approved
granted: 2026-10-16
action: give every Agent Session and Verification command of a Run the Run Marker and make the judge send nothing while it is set, and withhold TYPESAFE_API_KEY from Runs; remove the race of the bootstrap serialization test; move every Unix-only test into unix-constrained files, re-record the coverage record and add a Repository Contract Test that compiles every test for Windows; let Task Carry-Forward carry a settled Task over a planning amendment that leaves its Binding Text unchanged; make the history cleanup read every form of ADR reference adopters use; and describe each change in the user guides, the Roundfix skill references and the glossary
consuming: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/reconcile.md
  - .agents/skills/roundfix/references/spec.md
  - Makefile
  - docs/references/coverage-record.json
  - internal/baseline/repository_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0277

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion: the
authorization "Governed Paths declarados" (the option text: "\"Concedo\" como
em AE/AF/AG para os Governed Paths que cada Spec declarar: módulos da
Baseline, registros de versão, make baseline-digests, skills e testes
governados. Makefile, workflows de CI, go.mod e configuração de lint ficam
fora.") and the failure rules "Regras de AE–AG (Recomendado)". On 2026-10-17,
at about 15:15 BRT, the maintainer chose the items of this Spec, "Chave do juiz
entra no Run, Teste instável no worktree, Testes no Windows, Carry-forward e
emendas", the next step "Escrever, entregar e lançar (Recomendado)" and, for
the adopters, "Só nota no inbox (Recomendado)". The 2026-09-30 skills grant
("considere autorizado a ajustar todas as skills se necessário") also stands.

The governed set was measured on the authoring branch `docs/plan-0277`, based
on `23ab0904`, with `speccheck.GovernedPath`, through a `go test -overlay`
probe in `internal/speccheck` that wrote nothing to the repository (`git
status` was identical before and after, and the probe file never existed in
the tree), against every file the Tasks declare, every planning file of this
Spec and the ADRs it writes. Eight of them are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_03,
  task_04).
- `.agents/skills/roundfix/references/spec.md` covers `spec judge` and
  `.agents/skills/roundfix/references/archive.md` covers `archive` and
  `history sanitize`, whose judge output changes (task_01);
  `.agents/skills/roundfix/references/reconcile.md` covers carry-forward
  (task_03); `.agents/skills/roundfix/references/archive.md` and
  `.agents/skills/roundfix/references/baseline.md` describe the ADR reference
  forms of the history cleanup (task_04).
- `internal/baseline/repository_test.go` holds a test that calls
  `syscall.Mkfifo` without a build constraint; that test moves to a
  `unix`-constrained file (task_02).
- `docs/references/coverage-record.json` lists the platforms of every test and
  is re-recorded by its sanctioned command because the moved tests leave the
  Windows list (task_02).

## What is not governed

The production and test sources the Tasks change in `internal/judge`,
`internal/agent`, `internal/daemon`, `internal/cli`, `internal/worktree`,
`internal/spec` and `internal/verifyselect`, including the new
`unix`-constrained test files and the new Repository Contract Test; the skill
reference mirrors under `skills/roundfix/references/`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; `CONTEXT.md`; the user guides under
`docs/user-guide/commands/`; the ADRs; and every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions,
the Behavior Surface Record and the coverage record. This declaration records
the regeneration that follows the approved edits and adds no source path.

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

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No Baseline module changes, so no module version is recorded; an owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No record or mirror is hand-edited.

## Limits

- No change to the Makefile other than the one `REPO_CONTRACT_TESTS` entry
  the amendment of 2026-10-18 authorizes, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, or the `### QA settlement` section of any skill.
- No change to the bootstrap permit or any production code of the Worktree
  Bootstrap, the offline archive plan, the judge's questions, the Jev Ceiling,
  or the environment of Git hooks, the QA format command and the Worktree
  Bootstrap.
- No edit to any adopter repository; no test, Verification command or QA row
  reads or writes the real `~/.roundfix`, reads a shell startup file or a
  runtime's state directory, reaches the network, sends a request to a
  provider or the judge, or reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

## Amendment of 2026-10-18

The delivery's pre-PR review of 2026-10-17 found that task_02 had loosened
`internal/suiteguardcontract/repository_gate_test.go` so that a textual match
on the Makefile's Full Contract Run excused a repocontract test missing from
`REPO_CONTRACT_TESTS`, because the Makefile was outside this grant. On
2026-10-18 the maintainer answered through AskUserQuestion "Autorizar 1 linha
no Makefile (Recomendado)" (option text: "Acrescento só
TestEveryTestCompilesForWindows ao REPO_CONTRACT_TESTS, com emenda datada no
_authorization.md. O teste roda no gate de contratos (main e release), sem
custo nas tasks."). This amendment adds `Makefile` to `paths` for exactly that
change: `TestEveryTestCompilesForWindows` joins the `REPO_CONTRACT_TESTS`
list, and no other Makefile byte changes. The loosened audit is restored to
its delivery-base bytes in the same fix. It lands on the default branch
before the delivery consumes it.
