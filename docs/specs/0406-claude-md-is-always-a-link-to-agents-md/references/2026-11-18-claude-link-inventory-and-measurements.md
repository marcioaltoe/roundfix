# The root `CLAUDE.md` of the adopters, and the measurements of the design (2026-11-18)

The author of Spec 0406 recorded these facts on 2026-11-18, about 06:40 BRT, on the authoring host. Every read of another checkout was read only: nothing was committed there, no message was sent to another project, and each checkout's `git status --porcelain` was the same before and after.

## Inventory of the adopted checkouts

Method: for each checkout under `~/dev` and `~/srv` that holds `docs/agents/setup-context.json`, `ls -l CLAUDE.md AGENTS.md` and `git ls-files -s -- CLAUDE.md AGENTS.md`, which prints the mode Git records. Mode `120000` is a symbolic link; `100644` is a regular file.

| Checkout | Commit | `CLAUDE.md` on disk | `CLAUDE.md` in Git | `AGENTS.md` |
| --- | --- | --- | --- | --- |
| `~/dev/argus` | `fc53eb7c` | link to `AGENTS.md` | `120000` | regular file |
| `~/dev/quaestor` | `4cb991ec` | link to `AGENTS.md` | `120000` | regular file |
| `~/srv/quaestor` | `658272b0` | link to `AGENTS.md` | `120000` | regular file |
| `~/srv/fiscus` | `55bac1b1` | link to `AGENTS.md` | `120000` | regular file |
| `~/srv/fiscus-mcp` | `55bac1b1` | link to `AGENTS.md` | `120000` | regular file (second checkout of the same repository and commit) |
| `~/srv/indexus` | `3e48f0d0` | link to `AGENTS.md` | `120000` | regular file |
| `~/dev/roundfix` | `eb888135` | link to `AGENTS.md` | `120000` | regular file (this repository) |

Every adopted checkout on the host already keeps the standard: `CLAUDE.md` is a relative symbolic link whose target is exactly `AGENTS.md`, and Git records it as a link.

Two checkouts on the host hold no Setup Manifest and are not adopters. They show the two shapes the standard meets outside the fleet:

- `~/dev/secondbrain` keeps `CLAUDE.md` as a regular 11-byte file, `@AGENTS.md` and a newline, the import line Claude Code reads. Git records it as `100644`.
- `~/dev/skills` keeps it as a link to `AGENTS.md`.

Vortex, Fluxus, Oraculum, Conexus and Pantheon have no checkout on this host, so they are not in this inventory.

## What the Baseline does with `CLAUDE.md` today

Read from `internal/baseline` at `eb888135`:

- The inventory (`repository.go`) captures the root `AGENTS.md` and `CLAUDE.md` as mutable paths. A regular carrier is read; a symbolic link is an alias whose in-repository target is read. An alias whose target is absolute, escapes the repository, is missing or forms a cycle blocks planning with `baseline.inventory.unsafe-alias`.
- Adoption (`greenfield` and `preservation`) backs up every safe root carrier it would rewrite as `AGENTS.<digest>.md` or `CLAUDE.<digest>.md` (ADR-0070). The plan writes its root blocks only into `AGENTS.md`; no artifact path is `CLAUDE.md`, so adoption leaves a regular `CLAUDE.md` as it was and only writes its backup. Spec 0405's measurement saw the same.
- A managed refresh writes no root backup and verifies that every byte outside the Managed Regions of each changed carrier is unchanged (ADR-0100, `validateManagedRefreshPreservation`).
- The plan's postimages are regular files or deletions. The file transaction can already write a symbolic link (`createDestinationTemporary` uses `os.Root.Symlink` and renames the temporary over the destination), but only where the preimage is missing: `captureTransactionEntries` and the journal validation refuse a link postimage over an existing path, and `validateMutationDestination` refuses a destination that is a link. `baseline skills restore` is the only producer of link postimages today, for the Claude Skill Links (ADR-0341).
- The Misplaced Repository Rule detector (ADR-0411) reads a regular root `CLAUDE.md` and reports every passage outside its Managed Regions, so a `CLAUDE.md` holding only `@AGENTS.md` is reported as one passage to move. A link is skipped.

## Symbolic links and Windows

- Roundfix ships Windows binaries through its npm platform packages, and every test file compiles for Windows (ADR-0285); no test runs on Windows.
- The Baseline creates links in two places: the Claude Skill Links of a restore, only where `.claude/skills` already holds such a link (ADR-0341, which rejected linking a directory of copies "kept where symbolic links do not work"), and `roundfix skills install`, which creates `.claude/skills/roundfix` after a prompt.
- Git on a host that does not create symbolic links (`core.symlinks=false`, the default of Git for Windows without Developer Mode) checks a link out as a regular file whose bytes are the link's target, here `AGENTS.md` with no newline, and keeps mode `120000` in the index.

## The prototype

The author built a prototype on a scratch worktree of `eb888135`, removed afterwards.

### The guide sentence

The intro of `internal/baseline/assets/templates/guides/agent-instructions.md` gained the sentence of API Contract 4. The template's version (4 to 5) and the guide's version in `core.json` (17 to 18) were raised, then these ran in order:

1. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which raised `core` from 31 to 32;
2. `make baseline-digests`;
3. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, whose second run reported `current`.

Files they changed besides the three edited ones:

- `internal/baseline/module-versions.json`;
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and the formatter golden `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`;
- `internal/baseline/testdata/catalog.diagnostics.golden.json`, `catalog.digest` and `catalog.normalized.json`;
- the four plan goldens under `internal/baseline/testdata/plan-characterization/`: `advisory-only-divergences`, `clean-adoption`, `idempotent-replan-after-verified-apply` and `same-baseline-changed-profile-and-catalog-digests`;
- `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`.

No Source Baseline, Setup Snapshot, other profile, formatter golden `AGENTS.md`, `AGENTS.md` or `docs/agents/specific-repository.md` changed. With only this edit, `go test ./internal/baseline -count=1` passed (1,820 tests).

### The link

The prototype then planned `CLAUDE.md` as a link postimage with target `AGENTS.md` in `assemblePostimages`, for adoption always and for a managed refresh when the file was absent or held no passage of its own, and let the transaction replace a regular file or a link at that one path and restore it on rollback. After `make baseline-digests` regenerated the four plan goldens:

- `go test ./internal/baseline -count=1`: one existing test failed, `TestApplyExactDigest`, because its helper `assertAppliedPostimage` in `internal/baseline/apply_test.go` reads every postimage as a regular file; the plan goldens it would also have moved are the four the regeneration rewrites.
- `go test ./internal/cli -count=1`: passed (223 s). The prototype added no Doctor line; the line this Spec adds moves the Doctor tests that count lines, named in the TechSpec.

Runs of the prototype binary:

- A temporary repository with a regular `CLAUDE.md` of three lines, adopted with `roundfix baseline plan --profile go-cli-tui --decision preservation.mode=greenfield` and the other decisions of this repository's Setup Manifest that the profile selects, after copying this repository's `.agents/skills` and `skills-lock.json` and adding a `Makefile` with `verify` and `verify-changed` targets, then `roundfix baseline apply`: the plan listed `update CLAUDE.md` and the backup `CLAUDE.<digest>.md`; the apply verified 13 postimages; `CLAUDE.md` became a link to `AGENTS.md`; the backup held the three lines; a later `roundfix baseline update` reported `current`.
- On this repository with `CLAUDE.md` replaced by a regular file holding `@AGENTS.md`, `roundfix baseline update --yes` replaced it with the link and a second run reported `current`. Before the prototype exempted that conversion, the same run failed with `validate managed-refresh preservation: non-managed region digest mismatch at "CLAUDE.md"`.
- With `CLAUDE.md` removed, the update created the link.
- With `CLAUDE.md` a regular file holding a heading and a bullet, the plan left it unchanged.
- Read only, without `--yes`, in each adopted checkout of the inventory: no plan line named `CLAUDE.md`, and every `git status --porcelain` was unchanged.

## Governed Paths

Measured with `speccheck.GovernedPath` through a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository, over 58 paths.

- **Governed and changed by the design:** `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/index.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/setup-context-driven/SKILL.md`, `skills/roundfix/SKILL.md` and `skills/setup-context-driven/SKILL.md`.
- **Governed and not changed:** the formatter golden `AGENTS.md`, `internal/baseline/assets/profiles/go-cli-tui.json`, `internal/baseline/assets/setups/go.json`, `internal/baseline/assets/source-baselines/index.json`, `docs/agents/specific-repository.md`, `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, and the controls `Makefile`, `.github/workflows/ci.yml` and `go.mod`.
- **Not governed:** `AGENTS.md`, `CLAUDE.md`, `CONTEXT.md`, every Go source and test file of `internal/baseline` and `internal/cli` the design changes or adds except the two governed tests above, `internal/baseline/module-versions.json`, the catalog records and plan goldens under `internal/baseline/testdata/`, `docs/user-guide/commands/baseline.md` and `doctor.md`, `docs/references/behavior-surfaces.json`, the mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json`, the ADRs and this Spec's files.
