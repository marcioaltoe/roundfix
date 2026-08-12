---
granted: 2026-08-12
action: Relocate every retired documentation family under one `docs/history/` root including Review Artifacts, migrate the pre-0085 per-tree layout when the Baseline runs, and exclude the history tree from review completely.
consuming: 0094-one-history-root-under-docs
paths:
  - internal/spec/archive.go
  - internal/spec/archive_layout_characterization_test.go
  - internal/spec/archive_test.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/docscontract/doc.go
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - docs/agents/docs-layout.md
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/brainstorming/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .coderabbit.yaml
  - .roundfixrc.yml
  - skills/baseline_skill_contract_test.go
---

# Tooling authorization — the archive root under `docs/` (2026-08-12)

On 2026-08-12 the maintainer read the finding that `_archived/` at the repository
root is Spec 0085's deliberate choice, and directed the correction anyway, with
its scope and its authorization in one message:

> 1. O unico repositório que está usando essa implementação da spec 0085 é o
> roundfix, podemos ajustar isso localmente e depois commitar. É por isso que
> estou dando prioridade e pedindo o release de fix, assim posso usar o roundfix
> nos demais repositórios. O que eu quero em paralelo é que o roundfix detecte a
> existencia do _archived dentro de findings/specs/backlog e outras pastas de
> docs/ e faça, automaticamente, a movimentação desses arquivos para a pasta
> docs/history/{findings|specs|backlog} quando for executado o baseline. Isso
> deve entrar nessa spec de fix.

> 2. Está autorizado a alteração de tooling para que tenhamos o objetivo
> alcançado, inclusive atualizar o `.roundfixrc.yml` e `.coderabbit.yaml`
> (coderabbit deve ignorar completamente as pastas _archived no review).

## What this covers

**Retired documentation moves under one `docs/history/` root.** `spec.ArchiveDir`
returns the repository-relative literals and is the single resolver Spec 0085
Task 03 built for exactly this kind of change; `ArchiveSpecRoot` joins it from the
repository root and needs no edit. Every consumer already reads the resolver
rather than a literal. The bytes move with `git mv`, and the skills, Baseline
modules, and setup-owned guide that name the old path in prose follow it.

**Review Artifacts become a retired family.** A review whose owning Spec is known
already writes inside that Spec at `<spec>/reviews/` and travels with it; that
path is unchanged. The orphan case — a Pull Request with no owning Spec, fifty
folders in this repository on 2026-08-12 — gains `docs/history/reviews/` as its
terminal home and loses its underscore in the live root, so both read the same.
Which orphan is finished is decided by local Git reachability under ADR-0123,
never by the hosting provider, because the migration must run offline.

**The Baseline migrates the old per-tree layout.** Repositories adopted before
Spec 0085 still hold `docs/specs/_archived/`, `docs/findings/_archived/` and the
equivalents inside other `docs/` directories. `roundfix baseline` detects those
and moves their contents into `docs/history/{specs,findings,backlog,adr,reviews}`.
This is new product behaviour rather than a tooling mutation, and it is recorded
here because the maintainer attached it to this same Spec.

**CodeRabbit ignores the history tree completely**, which the maintainer restated
on 2026-08-12:

> Lembre-se de adicionar docs/history/ no .coderabbit.yaml para não ser revisado
> nunca.

Because the root carries no underscore, the existing prefix-matched exclusion no
longer reaches it, so `path_filters` gains a path-anchored `docs/history/**`
exclusion. The same anchoring closes a hole the prefix match never covered: a
review directory inside an active Spec is not matched by `!**/_reviews/**` and is
therefore reviewed today. This grant also authorizes confirming that no
`knowledge_base.filePatterns` entry reaches the history tree — that file's own
comment records that negation is undocumented there, so non-reachability is the
only protection — and correcting the stale comment naming the old root.

## Authorized paths

- `internal/spec/archive.go`, limited to the `ArchiveDir` literals and the
  doc comments that describe the layout.
- `internal/spec/archive_layout_characterization_test.go` and
  `internal/spec/archive_test.go`, limited to the relocated layout.
- `internal/docscontract/testdata/corpus-golden.json`, re-recorded with its
  reason stated, never regenerated silently.
- `internal/docscontract/doc.go`, limited to the stale package comment naming
  `docs/**/_archived`.
- `internal/baseline/assets/modules/spec-workflow.json` and
  `internal/baseline/assets/modules/context-workflow.json`, limited to the
  archive-location clauses.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
  and `.../manifest.json`, limited to the same subject.
- `docs/agents/docs-layout.md`, which carries the setup-owned managed regions
  rendered from the catalog above.
- `.agents/skills/archive-spec/SKILL.md`, `.../write-prd/SKILL.md`,
  `.../write-tasks/SKILL.md`, `.../write-techspec/SKILL.md`,
  `.../write-idea/SKILL.md`, `.../brainstorming/SKILL.md`, and
  `.../roundfix/SKILL.md`, limited to the archive path they name. All seven are
  Roundfix-owned per `skills/skills.go`.
- `.coderabbit.yaml`, limited to the archive comment, the path-anchored
  `docs/history/**` exclusion, and any `path_filters` or
  `knowledge_base.filePatterns` entry that would otherwise reach the history
  tree or a Spec-owned review directory.
- `.roundfixrc.yml`. Authorized by the maintainer's message; no required edit is
  identified today, and this grant does not oblige one. If the Spec finds none,
  the file stays untouched rather than edited to honour the grant.
- `_archived/**` moving to `docs/history/**`, preserving history.
- `docs/specs/_reviews/**` moving to `docs/history/reviews/**` for orphan reviews
  that local reachability proves finished, preserving history.
- `internal/config/config.go`, limited to the Review Artifact root resolution.
  This is production Go rather than protected tooling and needs no grant; it is
  listed so the Spec's bounded-scope audit has one complete set to read.
- `skills/baseline_skill_contract_test.go`, limited to the assertion that pins
  the archive path literal. Added 2026-08-12 after the authorized skill edits
  landed and left it stale, on the maintainer's decision:

  > Estender o grant em um arquivo

  The test asserts that `write-prd`'s skill contains
  ``Exclude `_archived/specs/` from automatic link rewrites``. The authorized
  edit correctly moved that path, so the assertion now names a location the
  contract no longer has. This is a consequent fix of an authorized change and
  lands in its own commit after it, never folded into it. It is not a licence to
  change what the contract requires — only where it says the archive lives.

Derived Baseline pins and skill mirrors rewritten by `make baseline-digests` and
`make skills-sync` are sanctioned fallout under ADR-0081, not separate targets.

## Bounded by purpose

This grant covers the archive root's location, the Baseline migration of the
pre-0085 per-tree layout, and the review exclusion of archived trees. It does not
authorize other Baseline clauses, other modules, other guides, other skills, any
`.coderabbit.yaml` key outside `path_filters` and
`knowledge_base.filePatterns`, or any `.roundfixrc.yml` key unrelated to the
archive.

## Chronology

This record lands as its own commit, before the commit that consumes it. A
prerequisite fix repairing something already red may land before either. A
consequent fix made necessary by the authorized change lands after it, never
folded into it.

## Consuming Spec

`0094-one-history-root-under-docs`, the priority correction the
maintainer directed to ship before every other Spec in
`docs/design/2026-08-12-the-spec-set-the-evidence-asks-for.md`, with a fix
release after it.

## What it supersedes

Spec 0085's Goal 2, "retired material leaves the directories an Agent loads by
default", for the archive root specifically. The consuming Spec states that
supersession in its PRD; leaving it unstated would establish the exception by
precedent. Goals 1 and 3 — one archive root, one path filter — are preserved
unchanged, which is why this is a correction to where 0085 anchored the root
rather than a reversal of what it built.
