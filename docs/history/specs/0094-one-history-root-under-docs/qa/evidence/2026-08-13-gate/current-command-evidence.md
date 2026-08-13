# Current command evidence — Spec 0094 QA gate

Build: `9e85a41e99dfb1f0ec7d1e3128ed69269e4f55fc`, rebuilt as
`bin/roundfix` with `go build -buildvcs=false`.

## Preconditions and static gates

- `rtk ./bin/roundfix spec check 0094-one-history-root-under-docs --strict`
  exited 0 with `No findings.`
- Every dependency of terminal QA Task `task_09` was `completed`.
- `rtk make verify` exited 2. `TestCatalogCompatibility` and four
  `TestBaselinePlanCharacterization` shapes found committed catalog digest
  `sha256:507ceef…` where the embedded catalog produces `sha256:9773febe…`.
- `rtk make verify-docs` exited 2. The docs-contract package passed; the
  repository-contract regeneration checks failed on the same stale catalog and
  plan-characterization artifacts.
- The focused reproduction failed at current HEAD and passed in a Git archive
  of Task 16's parent. Task 16 changed the embedded Source Baseline, while
  `internal/baseline/testdata/catalog.digest` and four plan-characterization
  goldens retained the parent digest.

## Outside fleet migration

Source: `/Users/marcio/dev/vortex` at
`07826b268be7e77fbe2f6231fd5d151b54362cc8`, an independently adopted
repository this Spec did not build. The source held 413 files and 13,600 KiB
under `docs/specs/_archived/`; it remained untouched.

The only target was `/tmp/roundfix-0094-qa-vortex`, created with
`git clone --no-local`. An unconfirmed `baseline update --adopt-suggested
--no-skills` exited 3, reported `repository bytes are unchanged`, and planned
413 History Relocations. Approved `--yes` exited 0 with `Baseline update:
verified`, removed the legacy tree, and placed 413 files under
`docs/history/specs/`.

The sorted content-digest set before and after was identical:

```text
742aab1179690b00c65ff5c5187d0956ac21e47666262516e26c073fe4988e88
```

A fresh rerun exited 0 with `Baseline update: current`, `File changes: 0`, and
no History Relocation block.

## Lifecycle and Review Artifact flow

In `/tmp/roundfix-0094-qa-flows`, approved Baseline update moved a
`superseded` ADR, a `declined` Backlog Entry, and the `round.md` of a finished
orphan Review Artifact. A `proposed` ADR and `open` Backlog Entry stayed active
with their original SHA-256 identities. A reachable non-ancestor review stayed
live with a `baseline.history.review.live` reason; a no-head review stayed with
the distinct `baseline.history.review.undecidable` reason. A Review Artifact
inside the active Spec stayed at `<spec>/reviews/`.

The finished review's file moved with identity
`7144940d2294a0ebb277103e1a838ad6d5a22ae698ce0c0601244559ecaf66d9`,
but its two empty source directories remained. A fresh Baseline update then
reported `docs/specs/reviews/qa-finished` as undecidable because the moved
`round.md` could no longer be read. Source inspection confirms
`relocateHistoryMove` renames files and syncs parents but never removes emptied
source directories.

## Collision and output flow

A disposable dual-layout tree held one occupied destination and one free
sibling. Approved `baseline update --yes` exited 3, named both colliding paths,
`already exists`, and `not every history relocation was performed`. The sibling
moved with its original identity; both colliding files retained their distinct
bytes. A fresh text and JSON rerun stayed `plan_ready`, carried the one
outstanding `historyMoves` entry, and never reported current. The JSON entry had
only `ordinal`, `from`, `to`, and `contentIdentity`; no file content appeared.

## Archive, review configuration, and vocabulary

- `roundfix archive --help` names `docs/history/specs/<slug>/` for the built-in
  Spec Root and `<spec-root>/_archived/<slug>/` only for a non-default root.
- A disposable completed `9999-qa-history-probe` archived successfully to
  `docs/history/specs/9999-qa-history-probe`; a fresh read found the source gone
  and archive metadata stamped.
- Focused `TestReviewHistoryConfiguration*` checks all passed, including
  negative controls that remove `!docs/history/**` or
  `!docs/specs/**/reviews/**` and inject rule sources reaching either tree.
- `CONTEXT.md` defines `History Root` and `History Relocation` and documents
  `HistoryRelocation`, `HistoryMove`, and `historyMoves`; no glossary edit is
  needed.
- Focused lifecycle, liveness, inventory-budget, collision, reporting, resolver,
  and Review Artifact root tests all passed.

## Tooling chronology and scope

Fresh `git diff-tree --no-commit-id --name-only -r <commit>` reads confirmed the
initial and earlier extended grants precede Tasks 02, 07, 08, 10, and 11, and
that Tasks 10 and 11 follow the carrier change that caused them.

Task 16 is different. At parent `65c51eb^`, neither `task_16.md` nor an
authorization for the routing carrier, `baseline.json`, or `index.json`
existed. Commit `65c51eb` created Task 16 and changed those protected paths plus
the Source Baseline manifest. Commit `e300c73` later added the authorization and
Task Graph node. The ancestry commands proved the order:

```text
git merge-base --is-ancestor e300c73 65c51eb  # exit 1
git merge-base --is-ancestor 65c51eb e300c73  # exit 0
```

The current Task 16 bounded list also omits the manifest changed by `65c51eb`.
The same commit did not regenerate all artifacts claimed by
`make baseline-digests`, which is the static-gate failure above.

## Evidence hygiene

All mutable probes stayed under `/tmp`; the external source checkout was
read-only. The built binary remains under ignored `bin/`. No commit, push, Pull
Request, or external mutation occurred.
