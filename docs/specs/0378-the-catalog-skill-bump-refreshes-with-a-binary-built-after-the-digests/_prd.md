---
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
status: active
created: 2026-11-10
surfaces: [cli]
type: fix
---

# The Catalog Skill Bump refreshes with a binary built after the digests

`make skills-third-party-bump` refreshes the Setup Manifest with a Roundfix binary it built before its own `make baseline-digests` step, so when that step rewrites an embedded Baseline asset the manifest records the catalog digest of the old assets. The bump still exits 0 and prints its final line, because its second refresh runs the same stale binary and converges with itself; only `TestTheSetupManifestIsCurrent` notices, and the author has to run the Managed Refresh by hand, as Spec 0370 had to. After this Spec the bump builds the binary again after `make baseline-digests`, so the Setup Manifest it writes is current with no extra command.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md`](references/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md),
   the Backlog Entry recorded by the operator on 2026-11-10 from Spec
   0370's TechSpec step 5. The operator had not reproduced it outside that
   Spec; this Spec did (see Validation).

At the authoring base `docs/backlog/` holds one other open entry, the
Portuguese not-applicable marker, which Spec 0379 adopts, and
`docs/findings/` none. None of them changes the
bump, its recipe or the Setup Manifest's catalog digest, so none shares this
context.

## Prerequisites

None in `requires`. Authored against `b6cd0f14` (v0.137.9 plus the 0374,
0375, 0376 and 0377 plans); the reproductions ran on `b233761e`, which
differs from it only in those plans' documents. The queue delivers Specs 0373, 0374, 0375, 0377 and 0376 before this
one. This Spec shares two files with them, and the one Task regenerates the
shared one with the repository's command, never by hand:

- the `Makefile`, which Spec 0377 also edits (its `test-budget` target, a new
  `suite-budget-check` target and `SUITE_BUDGET_SECONDS`); this Spec changes
  one line inside the `skills-third-party-bump` recipe only;
- `docs/agents/setup-context.json`, whose verification entries record the
  `Makefile` digest and whose catalog digest moves with every Baseline
  change, which Specs 0373 (Baseline clauses) and 0377 (`Makefile`) also
  regenerate.

It touches no owned skill, no skill version record, `CONTEXT.md` or
`docs/references/behavior-surfaces.json`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed; the Setup Manifest keeps its fields and only receives
  the digest the current catalog computes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the bump's only network steps,
  resolving and fetching the catalog commit, are unchanged, and every test
  uses a local `file://` catalog. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0360 makes the Catalog Skill Bump
  the only writer of the Vendored Skills and fixes its order: asset sync, a
  build, the confirmed restore, `make baseline-digests` and the Managed
  Refresh; this Spec keeps that order and makes the refresh run on a binary
  that embeds what the digests step wrote (ADR-0360: "It then runs the
  existing writers in their order"). ADR-0364 and ADR-0380 are not
  applicable: they cite ADR-0360, and the presence-only check of a required
  external skill and the retirement of a replaced catalog skill are
  unchanged. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answer "Mesmas do BB, teto minor (Recommended)" grants
  "Makefile na 0378 só nos arquivos que elas exigirem"; bounded files:
  `Makefile`, `docs/agents/setup-context.json`. The `Makefile` gains one
  line in the `skills-third-party-bump` recipe, and the Setup Manifest is
  regenerated only by the Managed Refresh command because it records the
  `Makefile` digest. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests/_authorization.md`;
  bounded files: `Makefile`, `docs/agents/setup-context.json`.

## Goals

- After a Catalog Skill Bump the Setup Manifest is current with no extra
  command, also when the change edits a Baseline module whose derived
  embedded assets `make baseline-digests` rewrites inside the bump.
- The refresh the bump runs embeds every asset the bump wrote before it, so
  its own "File changes: 0" check judges the tree the repository holds.

## User Stories

1. As the operator running a Catalog Skill Bump in a change that also edits a
   Baseline module, I want the Setup Manifest the bump writes to be current,
   so that `TestTheSetupManifestIsCurrent` passes without a manual Managed
   Refresh.
2. As a maintainer reading the bump's success line, I want it to mean the
   manifest matches the tree, so that a stale digest cannot hide behind a
   refresh that converged with an old binary.

## Core Features

1. The `skills-third-party-bump` recipe builds `$tmp/roundfix` again after
   its `make baseline-digests` step and before its two Managed Refresh runs.
2. The Catalog Skill Bump Makefile contract proves the order through its
   recording stubs: the pinned call sequence gains that build, and a new test
   proves the refresh runs a binary that embeds what the digests step wrote.

## Non-Goals / Out of Scope

- A Baseline module edited after a bump, or a `make baseline-digests` run on
  its own: every change to an embedded Baseline asset needs the Managed
  Refresh, as the repository rule on derived Baseline files already says;
  the bump is the only command that promises a current manifest.
- Running the Managed Refresh inside `make baseline-digests`, or a currency
  check inside the bump: the root cause is the binary's age, and one build
  removes it.
- Any change to the bump's refusals, its network steps, the restore, the
  asset sync, `make baseline-digests`, the catalog digest, the Setup
  Manifest format, ADR-0360's order or `CONTEXT.md`.
- Running a real bump in Verification or QA: it needs the whole catalog,
  which lives in another repository; the recording stubs prove the recipe,
  and the QA reproduces the stale-binary mechanism with the real binary.

## Success Metrics

1. Success Metric: `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests`
   passes on the audited head and fails, naming the asset the refresh
   embedded, both when the new build is dropped and when it is moved before
   `make baseline-digests`.
2. Success Metric: with the real binary, a Managed Refresh run by a binary
   built before `make baseline-digests` rewrote the formatter golden of an
   edited module leaves `TestTheSetupManifestIsCurrent` failing, and the same
   refresh run by a binary built after it leaves the test passing.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry, recorded by the operator from Spec 0370's
  TechSpec: the bump "refreshes `docs/agents/setup-context.json` with a binary
  it built before `make baseline-digests` rewrote the profile's golden
  digest", so 0370 had to prescribe `go run -buildvcs=false ./cmd/roundfix
  baseline update --repo . --no-skills --no-history --yes --format text` run
  twice by hand.
- The Go `embed` package documentation, read through Exa on 2026-11-10
  (https://pkg.go.dev/embed): a `//go:embed` directive initializes a variable
  "with the contents of files read from the package directory or
  subdirectories at compile time." `internal/baseline` embeds `assets`, and
  `make baseline-digests` rewrites four directories under it, so a binary
  built before that step carries the old bytes.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Reading at `b233761e`: the recipe builds `$tmp/roundfix` before and after
  the asset sync, then runs the restore, `make baseline-digests` and two
  Managed Refresh runs with that binary; `internal/baseline/catalog.go`
  embeds `assets`, and `Catalog.Digest` hashes every embedded asset;
  `DERIVED_DIGEST_PATHS` includes `assets/setups`, `assets/source-baselines`,
  `assets/formatter-fixtures` and `assets/profiles`.
- Reproductions in disposable clones inside the authoring clone, with a
  local catalog repository built from the catalog content the lock pins and
  committed again under a new commit (no network):
  - a bump with no other change, and a bump after a `tui-surface` module
    edit (the `go-cli-tui` profile has no formatter golden): the bump exits
    0 and `TestTheSetupManifestIsCurrent` passes;
  - a bump after a `backend` module clause edit: `make baseline-digests`
    rewrites the `standard-typescript-monorepo` formatter golden and its
    `goldenDigest`, the bump exits 0 with its final line, and
    `TestTheSetupManifestIsCurrent` fails with "Setup Manifest is stale;
    Managed Refresh proposes changes: docs/agents/setup-context.json";
  - the same bump with one `$(GO) build -o "$$tmp/roundfix" ./cmd/roundfix`
    line after `make baseline-digests`: the bump exits 0 and the test passes.
- A prototype of the test change in a `git archive` copy: both named tests
  pass with the change, and each fails when the new build is dropped or
  moved before `make baseline-digests`.
- Published sources (Exa): the Go `embed` documentation quoted above.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "makefile target builds binary before regenerating embedded assets stale digest" --all --files --min-score 0.3`,
  which returned this repository's mirrored decision that a regeneration
  run is not gated on the pins it rewrites, the adopted Backlog Entry, `CONTEXT.md`
  and a note of lessons from this repository's history; none decides when
  the bump builds its binary, and none contradicts this design.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0378 of 2026-11-10, the adopted Backlog Entry, and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~20:55–21:00 BRT), verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among four, "Digest
  do bump de skills" — option text: "O bump de skills deixa um digest velho
  no Setup Manifest e exige um comando manual. Só aparece quando alguém roda
  o bump, o que é raro."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those and the evidence, with no further interview:

- The defect reproduces only when `make baseline-digests` rewrites an
  embedded asset inside the bump, which a module edit with a formatter
  golden does; the fix rebuilds after that step instead of telling the
  author which command to run, because the bump already owns the refresh.
- One build after `make baseline-digests` is enough: every embedded asset
  the bump writes (the snapshots from the asset sync and the derived assets
  from the digests step) exists by then, and the restore writes nothing
  embedded.
- When the author ran `make baseline-digests` before the bump, the bump's
  binary was already current and stays so; with the fix the order of the
  two commands no longer matters.
- No ADR: ADR-0360 already sets the order ending in the Managed Refresh;
  the fix makes that refresh read the tree the order produced.
- One implementation Task and the QA gate. The release is the operator's
  patch right after the merge; no Task releases.

## Open Questions

None.
