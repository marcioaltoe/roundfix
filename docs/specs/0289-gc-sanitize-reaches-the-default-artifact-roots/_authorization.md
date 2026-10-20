---
status: approved
granted: 2026-10-20
action: make one resolver serve the storage report, GC sanitation and Run Retention, make roundfix gc sanitize classify the Artifact Root each Run wrote to with its existing proofs and check absence again before each removal it proves, and describe the change in the gc guide, the Roundfix storage reference and the glossary
consuming: 0289-gc-sanitize-reaches-the-default-artifact-roots
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/storage.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0289

On 2026-10-20, at about 13:10 BRT, the maintainer approved cycle AL through
AskUserQuestion. The answers this Spec rests on were:

- sequence "Aprovar como está (Recomendado)" (option text: "Duas Specs
  separadas, um minor cada, nesta ordem; P primeiro porque trava o caminho de
  recuperação documentado.");
- scope of this Spec "Resolver + eu aplico (Recomendado)" (option text: "O
  sanitize resolve a raiz vazia como o Run Retention e usa as provas
  existentes; depois da release eu rodo o dry-run no seu ~/.roundfix e aplico
  só se ele sair limpo.");
- authorization "Mesmas regras do AK (Recomendado)": the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AJ rules that answer carries were approved on 2026-10-19, at about
16:20 BRT, among them "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and "Merge squash + release + push". The
2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix` or runs `roundfix gc
sanitize --apply` outside a temporary home; the maintainer runs it on the real
home himself after the release.

The governed set was measured on the authoring branch `docs/plan-0289`, based
on `7a678ddb`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted Backlog Entry, ADR-0301, the open Backlog Entry this planning
  change writes, `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Four of those files are governed. Three are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_03).
- `.agents/skills/roundfix/references/storage.md` covers `gc` and `storage
  report`, whose sanitation and resolver change (task_03).

## What is not governed

- The production and test sources the Tasks change in `internal/store` and
  `internal/cli`.
- The skill reference mirror `skills/roundfix/references/storage.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/gc.md`.
- Every planning file of this Spec, its adopted reference, ADR-0301 and the
  open Backlog Entry.

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
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change. No
record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/config/` or
  the `### QA settlement` section of any skill.
- No change to the storage report's output, to what the per-repository
  `roundfix gc` sweeps remove, to Run Retention's kept reasons, or to the
  sanitation proofs other than the unresolved proof placed before them and
  the recheck of absence before removal.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
