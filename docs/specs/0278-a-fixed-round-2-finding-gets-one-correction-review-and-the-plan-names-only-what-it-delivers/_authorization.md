---
status: approved
granted: 2026-10-16
action: give a fixed round-2 finding one Correction Review of its delta and let the Delivery Retry of such a park return the item to review whatever paths the fix touched; add the Spec Consistency Check code SC-VERIFY-PATH-UNDECLARED with its corpus golden and pin; add a documentation contract test that fails while the Setup Manifest differs from a fresh Managed Refresh; and describe each change in the user guides, the Roundfix skill references, the write-tasks skill and the glossary
consuming: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - internal/docscontract/publicdocs_test.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0278

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion. The
answers were the authorization "Governed Paths declarados" (option text:
"\"Concedo\" como em AE/AF/AG para os Governed Paths que cada Spec declarar:
módulos da Baseline, registros de versão, make baseline-digests, skills e
testes governados. Makefile, workflows de CI, go.mod e configuração de lint
ficam fora.") and the failure rules "Regras de AE–AG (Recomendado)".

On 2026-10-18, at about 09:20 BRT, after v0.88.0, the maintainer chose "Triagem
e Spec G (Recomendado)" (option text: "Faço a triagem e, se houver itens,
escrevo, entrego e lanço uma Spec G (v0.89.0) com até 4 tasks e as mesmas
regras. Candidatos: um round de correção quando a revisão do round 2 acha algo
real e a correção é estrutural, e um detector de critério da QA sem task que o
entregue. Paro na proposta seguinte."). For a real round-2 finding, the choice
was "Corrigir e mergear (Recomendado)". The 2026-09-30 skills grant
("considere autorizado a ajustar todas as skills se necessário") also stands.

The governed set was measured on the authoring branch `docs/plan-0278`, based
on `9471f698`, with `speccheck.GovernedPath`.
- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status` was identical before and after, and
  the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and the ADRs the plan writes.
- Twelve of those files are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_02).
- Three Roundfix references are governed:
  - `.agents/skills/roundfix/references/review.md` covers `review`, whose
    lineage gains the Correction Review (task_01).
  - `.agents/skills/roundfix/references/deliver.md` covers `deliver retry`,
    whose correction proof gains the round-2 exemption (task_01).
  - `.agents/skills/roundfix/references/spec.md` covers `spec check`, which
    gains the new code (task_02).
- `.agents/skills/write-tasks/SKILL.md` and `skills/write-tasks/SKILL.md` are
  the `write-tasks` skill and its mirror. They gain the `Verification paths`
  rule (task_02), with the version fields they carry.
- `internal/speccheck/coherence.go` and `internal/speccheck/constraints.go`
  stage and call the new detector (task_02).
- `internal/docscontract/testdata/corpus-golden.json` and its pin
  `internal/spec/archive_layout_characterization_test.go` record the new
  code's count of 0 (task_02).
- `internal/docscontract/publicdocs_test.go` is the always-relevant
  documentation contract file that receives the Setup Manifest test (task_04).
  A new file would break the nine-file pin, and a Repository Contract Test
  would need a Makefile entry.

## What is not governed

- The production and test sources the Tasks change in `internal/cli`,
  `internal/speccheck` and `internal/docscontract`, other
  than the files above, including every new file.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, the user guides under `docs/user-guide/commands/` and the
  ADRs.
- Every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
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
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, or the
  `### QA settlement` section of any skill.
- No change to round 1, round 2, Review Segments, the review bound, the
  corrective-Spec park, `internal/judge`, any judge question or threshold, the Jev
  Ceiling, the release plan or the Managed Refresh.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
