---
spec: 0034-release-plan
prd: _prd.md
created: 2026-07-16
---

# Release Plan — Technical Spec

## Executive Summary

The Release Plan Command adds a read-only `roundfix release plan` surface backed by a deterministic classifier over local Git history. The command uses Conventional Commit signals as minimum impact, recognizes a narrow maintenance-only path boundary for no-release changes, and stops for manual classification whenever neither source can prove the impact. Its primary trade-off is conservative friction: some legitimate internal changes require an explicit classification, but Roundfix never understates a release by inventing semantic meaning from filenames or source text. ADR-0048 keeps planning separate from every release mutation.

## System Architecture

The design extends existing seams without changing the tag-triggered release workflow:

- `internal/cli` owns `release plan` parsing, help, stdout/stderr separation, exit-code mapping, and JSON/text selection through the existing `Run(args, stdout, stderr) int` contract.
- A cohesive `internal/releaseplan` package owns version parsing, Git-derived input normalization, Conventional Commit classification, manual-impact validation, version-zero mapping, and the final Release Plan model.
- The existing Git process boundary supplies committed tags, revisions, commits, and changed paths. The package receives normalized records through a small interface and never shells out itself.
- The root help, release runbook, root Agent pointer, and canonical/embedded Roundfix skill describe the same command and approval contract. `make skills-sync` preserves the owned skill copies.
- `.github/workflows/release.yml`, npm package staging, GitHub Release creation, and the Upgrade Command remain unchanged.

The command creates no Run, opens no Run Database, reads no Roundfix configuration, and contacts no external service.

## Implementation Design

### Interfaces

```go
type Request struct {
    From, To string
    ManualImpact Impact
    ManualReason string
}

type GitSource interface {
    ResolveRange(context.Context, string, string) (Range, error)
    Commits(context.Context, Range) ([]Commit, error)
}

func Build(Request, GitSource) (Plan, error)
```

`internal/cli` supplies an adapter over the repository's existing Git runner. Tests can provide in-memory commit records to the classifier and real temporary repositories to the CLI boundary.

### Data Models

```go
type Plan struct {
    SchemaVersion string
    State State
    Base VersionRef
    Target RevisionRef
    Classification Classification
    ProposedVersion string
    Approval Approval
    Changes []ChangeEvidence
}
```

Stable values:

- `State`: `ready`, `approval_required`, `manual_classification_required`, `no_release`.
- `Impact`: `none`, `patch`, `minor`, `major`; `breaking` is a separate boolean because a breaking version-zero plan maps to a minor increment.
- `Classification.Source`: `conventional_commit`, `maintenance_only`, `manual`, or `mixed`.
- JSON schema identifier: `roundfix.release-plan/v1`.

Every `ChangeEvidence` carries the commit SHA, subject, parsed Conventional Commit type, detected breaking marker, automatic impact, and whether changed paths cross the maintenance-only boundary. The JSON result includes evidence for every commit in the range; the text result prints the decision first and then only the commits that determine or block it.

### API Contracts

```text
roundfix release plan [--from <tag>] [--to <revision>]
                      [--impact <none|patch|minor|major> --reason <text>]
                      [--format <text|json>]
```

- `--from` defaults to the latest stable semantic-version tag reachable from `--to`.
- `--to` defaults to `HEAD` and must resolve to a commit.
- `--format` defaults to `text`; `json` emits exactly one JSON object on stdout.
- `--impact` supplies the human classification for otherwise ambiguous changes. It is valid only with a non-empty `--reason`, cannot be lower than the automatic minimum, and does not count as approval for the resulting version.
- `--reason` without `--impact`, unsupported impacts or formats, non-semantic base tags, missing or reversed ranges, and dirty repository state are invalid input.
- The first release accepts stable `vMAJOR.MINOR.PATCH` bases only. Pre-release bases fail with a next action instead of being normalized silently.

The command inspects a clean working tree. Any tracked or untracked change blocks the plan because the uncommitted content is absent from the classified release range; stderr names the paths and tells the caller to commit, stash, or remove them.

Automatic classification uses the entire commit subject and body:

- a `!` breaking marker or `BREAKING CHANGE:` footer sets `breaking=true` and minimum impact `major`;
- `feat` sets minimum impact `minor`;
- `fix` and `perf` set minimum impact `patch`;
- all other or malformed commit types have no automatic increment.

For commits without automatic impact, a maintenance-only changed-path set can classify them `none`. The maintenance-only set is limited to repository planning evidence and verification-only artifacts: Spec/ADR/finding/handoff documents, test files and fixtures, and CI definitions that do not change release staging. Documentation outside planning remains no-release when it changes no shipped skill, package, configuration, CLI contract, or release asset. Any commit touching command/runtime code, configuration examples or schemas, shipped skills, distribution packages, release staging, or public compatibility documentation requires manual classification unless its Conventional Commit signal already establishes a minimum.

The final impact is the maximum automatic and manual impact. Version calculation follows:

- `none` → no proposed version and state `no_release`;
- `patch` → increment patch, state `ready`, approval not required;
- `minor` → increment minor and reset patch, state `approval_required`;
- breaking at major zero → increment minor and reset patch, state `approval_required`, `breaking=true`;
- breaking at major one or later → increment major and reset minor/patch, state `approval_required`, `breaking=true`.

An unresolved ambiguous commit produces no proposed version and state `manual_classification_required`. The text next action names the explicit rerun shape. An approval-required plan prints the exact question `Approve the <major|minor> increment to <version>?`; a version-zero breaking plan uses `minor` in the question and separately labels the compatibility break.

Exit codes use the existing Roundfix contract:

- `0`: `ready` or `no_release` plan produced;
- `2`: invalid arguments, Git/repository failure, dirty tree, invalid manual classification, or malformed version/range;
- `3`: `approval_required` or `manual_classification_required` plan produced successfully but a human decision remains.

Requested output goes to stdout. Diagnostics and the next corrective action go to stderr. JSON decision states remain on stdout even when the command exits `3`; invalid-input exits emit no partial JSON.

## Coverage Map

- Goal 1 and Story 1 → GitSource adapter, commit classifier, version calculator.
- Goal 2 and Story 3 → maintenance-only boundary, manual-classification state, validated manual impact.
- Goal 3 and Story 2 → approval model, version-zero mapping, ADR-0048 boundary.
- Goal 4 and Story 4 → CLI parser, versioned JSON schema, deterministic text renderer, exit codes.
- Goal 5 and Story 5 → documentation/Agent pointer updates and maintenance-only classification.
- Core Feature 7 → read-only integration tests that compare worktree, refs, and remote configuration before and after every state.

## Integration Points

- Local Git executable through the existing context-aware Git runner. No GitHub or npm access.
- `cog.toml` supplies the repository's Conventional Commit convention as documentation, not as a runtime dependency; the parser implements the required subset in stdlib Go.
- The release runbook and tag-triggered workflow remain the downstream human-operated release path.
- The canonical Roundfix skill and generated embedded copy receive the new command recipe under the existing skill-sync contract.

## Cross-Spec dependencies

This Spec has no functional dependency on Specs 0032 or 0033 and can be implemented independently after the deterministic verification fix lands. It must merge before Spec 0035 to keep public CLI dispatch, help text, documentation, and Roundfix Skill synchronization changes on a stable base.

## Testing Approach

Classifier unit tests use table-driven stdlib tests for stable version parsing, Conventional Commit subjects and bodies, impact ordering, manual-impact validation, version-zero behavior, and ambiguous/no-release boundaries. Every positive classification has a negative or ambiguous companion.

CLI integration tests create real temporary Git repositories with commits and tags. They exercise the exact acceptance ranges from the PRD, mixed commit order, explicit `--from`/`--to`, dirty state, missing tags, malformed commits, JSON schema fields, stdout/stderr isolation, and exit codes. Each scenario snapshots files, refs, tags, and remotes before invocation and asserts byte-for-byte state equality afterward.

Documentation contract tests pin the root help, command help, release runbook pointer, and canonical/embedded Roundfix skill synchronization. The release workflow compatibility test remains unchanged and proves the Release Plan did not alter artifact naming or tag-driven publication.

## Build Order

1. Release Plan domain model, semantic-version calculator, and table tests.
2. Commit and changed-path classification with manual-impact validation (depends on: 1).
3. Local GitSource adapter and temporary-repository integration fixtures (depends on: 1, 2).
4. `release plan` CLI parsing, text/JSON rendering, help, and exit contracts (depends on: 1, 2, 3).
5. Release runbook, root Agent pointer, and canonical/embedded Roundfix skill updates (depends on: 4).
6. Full acceptance matrix and read-only mutation audit (depends on: 3, 4, 5).

## Risks & Considerations

- Conservative ambiguity creates extra maintainer work. The mitigation is an explicit, non-mutating `--impact` plus reason rather than a heuristic that can under-version a breaking change.
- Conventional Commit history can be inaccurate. The changed-path boundary prevents unknown or nonstandard commits from silently becoming no-release when they touch shipped surfaces, but it cannot understand source semantics; manual classification remains part of the contract.
- A global manual impact can hide which ambiguous commit required it. The plan retains every ambiguous commit as evidence and records the supplied reason; per-commit overrides are deferred until evidence shows they are needed.
- Requiring a clean tree excludes tentative release analysis over local edits. That is intentional: only committed content can be tagged, reproduced, and published.
- The new command is public API. Command names, flags, JSON fields, states, and exit codes require compatibility treatment after release.

## Decisions

- Release planning is read-only and confirmation-gated. See ADR-0048.
- The public surface is `roundfix release plan`, not a command named `release` that appears to publish.
- Classification is deterministic and conservative; ambiguous shipped-surface changes require a supplied impact and reason.
- The first version supports stable release bases only and performs no network access.
- The implementation uses stdlib `flag.FlagSet` dispatch and existing Git seams; no CLI or semantic-version dependency is added.
