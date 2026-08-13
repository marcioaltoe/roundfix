# Command evidence — Spec 0094 QA gate

Build: `30be7bb727f9997372b06be86d3ba30021199390` (`bin/roundfix`, rebuilt in
the Run Worktree).

## Static gates

- `rtk ./bin/roundfix spec check 0094-one-history-root-under-docs --strict`
  exited 0 and printed `No findings.`
- `rtk make verify` exited 0. All Go packages passed; the focused Skill
  contract, Repository Skill Set check, and production build passed.
- `rtk make verify-docs` exited 0. The docs-contract, repository-contract,
  corpus-budget, and Spec consistency checks passed.

## Archive Command

A disposable clone received a completed `qa-history-probe` Spec and passing QA
Report. `bin/roundfix archive qa-history-probe` exited 0 and printed:

```text
archived qa-history-probe -> docs/history/specs/qa-history-probe
```

The source disappeared, the destination existed, and `_prd.md` carried
`status: archived`, `archived: "2026-08-13"`, and
`source_slug: qa-history-probe` on a fresh read.

The same binary's `archive --help` contradicted that result:

```text
repository's default _archived/specs/<slug>/ when the Spec Root is the built-in
docs/specs
```

`internal/cli/archive.go:22` owns the stale text.

## Outside fleet migration

Source: `/Users/marcio/dev/vortex` at
`07826b268be7e77fbe2f6231fd5d151b54362cc8`, an adopted repository this Spec
did not build. The source held 413 files and 13,600 KiB under
`docs/specs/_archived/`. The source checkout remained untouched.

The only migration target was a `git clone --no-local` under `/tmp`. The first
unconfirmed `baseline update --adopt-suggested --no-skills` exited 3, reported
`repository bytes are unchanged`, and planned 413 History Relocations. The
approved `--yes` run exited 0 with `Baseline update: verified`, reported all 413
moves, removed `docs/specs/_archived/`, and created 413 files under
`docs/history/specs/`.

The sorted per-file SHA-256 value sets before and after both hashed to:

```text
742aab1179690b00c65ff5c5187d0956ac21e47666262516e26c073fe4988e88
```

A fresh unconfirmed rerun exited 0 with `Baseline update: current`,
`File changes: 0`, and no History Relocation block.

## Decision and intent lifecycle

In a disposable clone, an added `superseded` ADR moved to
`docs/history/adr/`, while an added `proposed` ADR remained under `docs/adr/`.
An added `declined` Backlog Entry moved to `docs/history/backlog/`, while an
added `open` entry remained under `docs/backlog/`. Apply output named each move
with its SHA-256 identity, and fresh path and hash reads confirmed the result.

The current Roundfix checkout is not itself current. An unconfirmed
`bin/roundfix baseline update --repo . --no-skills` exited 3 and planned this
outstanding move:

```text
docs/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md
-> docs/history/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md
```

The source still declares `status: declined` in the active Backlog directory.

## Review Artifact liveness

A disposable clone used local Git only. A Review Artifact whose recorded head
was the local `main` commit moved from `docs/specs/reviews/` to
`docs/history/reviews/`. A reachable non-ancestor head remained live and was
reported as `baseline.history.review.live` with its reachability reason. A
missing head remained live and was reported distinctly as
`baseline.history.review.undecidable` with the failed local `git rev-parse`
reason. Fresh reads confirmed only the finished review's files moved.

No provider request, credential, or network lookup was used.

## Collision and reporting

A disposable clone held two legacy sources. One destination was occupied and
one was free. Approved `baseline update` exited 3, named the colliding source,
destination, and `already exists`, and stated that not every History Relocation
was performed. The sibling reached its destination with the recorded identity;
both colliding files retained their distinct bytes. A fresh unconfirmed rerun
again exited 3 with `History moves: 1` and did not print `current`.

## Review configuration

No Pull Request is open, so the live Review Source journey is environment
blocked. Equivalent repository evidence passed through `make verify-docs`:
the docs-contract suite loads `.coderabbit.yaml`, proves
`!docs/history/**` and `!docs/specs/**/reviews/**`, and rejects rule-source
patterns that reach either protected tree.

## Tooling chronology and scope

`git diff-tree --no-commit-id --name-only -r <commit>` was read for the tooling
commits. The initial authorization commit
`385b5f85a8720eb03c92ac1761df7a1cd4017a15` precedes the layout and carrier
commits. The extended authorization carried into the PRD/TechSpec at
`53762d7e99f673c02cfd42cd4b3640dcbba4285e` precedes the pinned assertion and
review-configuration commits. The 2026-08-13 extension at
`537f4e091985eccbab65278418e9ce1be799a35e` precedes the final derived-artifact
test fix. `git merge-base --is-ancestor` exited 0 for every checked edge, and
the consequent fixes landed after the carrier change rather than inside it.

Task 07's extra Baseline snapshots and Skill mirrors are deterministic output
of the sanctioned regenerators. `make verify-docs` freshly proved their current
byte contracts. No current QA delta exists outside this report and its evidence.

## Documentation and vocabulary sweep

`CONTEXT.md` defines both `History Root` and `History Relocation`, including the
`HistoryRelocation`, `HistoryMove`, and `historyMoves` emitted tokens. Active
guides, modules, and the seven canonical Roundfix-owned Skills name
`docs/history/specs/`.

Two shipped carriers remain stale:

- `internal/cli/archive.go:22` says the built-in root is
  `_archived/specs/<slug>/`.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md:22`
  says completed Specs archive under `_archived/specs/`.

## Non-Goal sweep

The outside-repository value-set hashes and the lifecycle probe's per-file
hashes confirmed byte-preserving moves. The implementation diff did not change
the Spec-owned `docs/specs/<slug>/reviews/` location, Archive Command
eligibility, or the detector population beyond the authorized history
classifications. Review liveness used local Git only; no provider request,
credential, or network lookup was introduced or exercised.

## Evidence hygiene

All scratch repositories and built test artifacts stayed under `/tmp` or the
ignored `bin/` build output. `git ls-files -s` found no `160000` gitlink under
the Spec's QA directory. The QA-authored repository delta is limited to this
evidence file and the seeded report.
