---
type: fix # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-12
spec: 0094-one-history-root-under-docs # Spec slug when status: promoted
reason: null # required when status: declined
---

# The archive root sits beside `docs/` instead of inside it

## Symptom

`_archived/specs/`, `_archived/findings/` and `_archived/adr/` are top-level
directories of the repository, beside `cmd/`, `internal/` and `skills/`. They
hold 156 files and 25 MB of documentation.

Spec 0085 put them there on 2026-08-12, building the single archive root the
maintainer requested on 2026-08-09. The request, recorded verbatim in
`docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`,
named the structure and not its parent:

> mudar a estrutura dos arquiveds de pastas dentro dos folders findings, specs,
> backlog e outras para uma estrutura com `_archived/specs|findings|adr|backlog`
> e assim mantemos os arquivos ativos e o histórico dos inativos e bloqueamos a
> revisão do coderabbit facilmente

The structure is right and the parent is not: an archive of documentation is
documentation, and it belongs under `docs/`.

## Where

`internal/spec/archive.go` — `ArchiveDir` returns the four repository-relative
literals; `ArchiveSpecRoot` joins them from the repository root and needs no
change. `internal/speccheck/citations.go` already reads the resolver rather than
a literal. The six skills that name `_archived/specs` in prose — `archive-spec`,
`write-prd`, `write-tasks`, `write-techspec`, `write-idea`, `roundfix`. The
Baseline modules `spec-workflow.json` and `context-workflow.json`, the
source-baseline corpus and manifest, and the managed regions they render into
`docs/agents/docs-layout.md`. The two archive-layout characterization tests and
the docscontract corpus golden.

## Expected

The archive root resolves under `docs/_archived/`, keeping the four families the
original request named. Both reasons that request gave survive the move:

- Blocking review stays one filter. `.coderabbit.yaml` carries
  `- "!**/_archived/**"`, which matches at any depth, so `docs/_archived/**` is
  excluded with no edit. Only the adjacent comment naming the root goes stale.
- Active stays separate from history. `speccheck` walks `docs/specs` and
  `docs/findings` exactly, not `docs/**`, so the archived trees stay outside the
  active ones by construction, exactly as they are today.

What does not survive untouched is Spec 0085's Goal 2, "retired material leaves
the directories an Agent loads by default". The owning Spec supersedes that goal
explicitly or does not move the bytes; establishing the exception by precedent is
the outcome to avoid.

The fleet needs migrating in the other direction. Roundfix is the only repository
carrying Spec 0085's root `_archived/`, so nothing outside it has to move away
from that. Every other adopted repository still holds the pre-0085 per-tree
layout — `docs/specs/_archived/`, `docs/findings/_archived/`, and the equivalents
inside other `docs/` directories. `roundfix baseline` should detect those and
move their contents into `docs/_archived/{findings,specs,backlog}`, which is what
makes a fix release usable in the rest of the fleet.

## Evidence

`_archived/specs/0085-what-an-agent-reads-before-it-decides/task_04.md`, which
performed the move and lists its bounded scope;
`_archived/specs/0085-what-an-agent-reads-before-it-decides/_prd.md`, Goals 1–3;
`internal/spec/archive_layout_characterization_test.go`, which pins the current
layout.
