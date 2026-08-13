# Final command evidence — Spec 0094 QA gate

Build: `f5c1f5a5a485cdafc05eeb28ef57e14af409cc35`; `rtk make verify`
rebuilt `bin/roundfix` with `-buildvcs=false` from this Run Worktree.

## Preconditions and static gates

- `rtk roundfix spec check 0094-one-history-root-under-docs --strict` exited 0
  with `No findings.`
- The manifest names `task_09` as the terminal QA Task. Every direct dependency
  is `completed`; the PRD has no `## Unreachable Acceptance` declaration.
- `rtk make verify` exited 0. The full Go suite, focused Skill contract,
  Repository Skill Set check, and production build passed.
- `rtk make verify-docs` exited 0. The docs-contract, repository-contract,
  declared regeneration/frozen-boundary checks, corpus budget, and active Spec
  checks passed.

## Tooling authorization, scope, and chronology

Fresh `git diff-tree --no-commit-id --name-only -r <commit>` reads covered the
Task commits that changed protected carriers or their assertions:

- Task 02 `e48d4548b028ce745e266a4c3824cee74e83415a` follows initial authorization
  `385b5f85a8720eb03c92ac1761df7a1cd4017a15`.
- Task 07 `ff4b2e2989fa1388080713f1eae55d2b2ac22eec` follows the same authorization.
  Its extra Source Baseline snapshots, catalog outputs, plan-characterization
  goldens, and seven shipped Skill mirrors are the deterministic outputs of the
  sanctioned regenerators. The fresh `make verify-docs` run proved those
  current byte contracts.
- Task 08 `89864eb7d86d228fde22269c1065462dab2ea2d9` and Task 10
  `d5b3d1b8dbe60abe13d26c9fe2947cb646bbf813` follow the extension authored by
  `4d72840ec51476220d5a0cb8f374da75d18f4868`.
- Task 11 `bda14b926620055a4fd0f0d5fb671fee7228c238` follows the
  `skills/owned_skill_edit_repocontract_test.go` extension at
  `537f4e091985eccbab65278418e9ce1be799a35e`.
- Task 16 `c80e1266658929f68e8046af82f88e13392dc56d` follows authorization and
  Task authoring at `0d7c3ba800185bfaca18040984ce590abee179ca`. Its direct Source
  Baseline paths match the bounded list; its catalog and plan-characterization
  outputs are reproducible sanctioned fallout, freshly proved by
  `make verify-docs`.

One combined `git merge-base --is-ancestor` audit exited 0 for all eight
required edges: the five grant-to-change edges above, Task 07 before Task 10,
and Task 10 before Task 11. Authorization records contain no protected carrier
change. Consequent fixes are separate commits after the change that made them
necessary. No protected change is folded into its authorization commit.

## Outside-evidence fleet migration

Source: `/Users/marcio/dev/vortex` at
`07826b268be7e77fbe2f6231fd5d151b54362cc8`, an adopted repository this Spec
did not build. It held 413 files under `docs/specs/_archived/`. The mutable
target was a fresh `git clone --no-local` at
`/tmp/roundfix-0094-final-vortex.taTdP0/repo`; the source checkout stayed at the
same commit with its legacy directory intact.

- The unconfirmed `baseline update --adopt-suggested --no-skills` exited 3,
  reported `repository bytes are unchanged`, and planned 413 History
  Relocations.
- The approved `--yes` run exited 0 with `Baseline update: verified`, removed
  `docs/specs/_archived/`, and created 413 files under `docs/history/specs/`.
- The sorted per-file SHA-256 value set produced the same aggregate digest before
  and after:

  ```text
  742aab1179690b00c65ff5c5187d0956ac21e47666262516e26c073fe4988e88
  ```

- A fresh unconfirmed rerun exited 0 with `Baseline update: current` and `File
  changes: 0`; it printed no History Relocation block. Warnings for unrelated
  retained orphan reviews remained visible and did not represent a second move.

## Built-in Archive Command flow

The built binary's `archive --help` names
`docs/history/specs/<slug>/` for the built-in Spec Root and reserves
`<spec-root>/_archived/<slug>/` for a configured non-default root.

A fresh disposable clone received completed Spec `9999-qa-history-probe` with a
passing QA Report. `roundfix archive 9999-qa-history-probe` exited 0 and printed:

```text
archived 9999-qa-history-probe -> docs/history/specs/9999-qa-history-probe
```

A fresh path read found the source absent and the destination present. Its
`_prd.md` carried `status: archived`, `archived: "2026-08-13"`, and
`source_slug: 9999-qa-history-probe`. The current repository has no root
`_archived/` and has all five family directories under `docs/history/`.

## Decision and Backlog lifecycle flow

A fresh disposable clone received four documents: a `proposed` ADR, a
`superseded` ADR, an `open` Backlog Entry, and a `declined` Backlog Entry. The
unconfirmed plan named exactly two History Relocations. Approved Baseline update
exited 0 and reported both verified moves:

- `docs/adr/qa-superseded.md` to `docs/history/adr/qa-superseded.md` with
  identity `29e7fc2b…a14`.
- `docs/backlog/2026-08-13-qa-declined.md` to
  `docs/history/backlog/2026-08-13-qa-declined.md` with identity
  `587f49b4…ba5`.

Fresh reads found the proposed ADR and open Backlog Entry in their active
directories with their original identities (`7f8dfc52…6b15` and
`094916fa…adf8`). Both moved documents retained their original identities. A
fresh JSON rerun exited 0 with state `current`, an empty `fileChanges` array,
and no `historyMoves` field.

## Orphan Review Artifact lifecycle flow

A disposable local-Git repository contained four cases: a finished orphan whose
recorded head was an ancestor of the default branch, a reachable non-ancestor
head, a newest Round with no recorded head, and a Spec-owned Review Artifact.
The unconfirmed plan was read-only and named only the finished orphan as a
History Relocation. It emitted distinct `baseline.history.review.live` and
`baseline.history.review.undecidable` warnings for the retained cases.

The approved update exited 0 and verified this move:

```text
docs/specs/reviews/qa-finished/round-001/round.md
-> docs/history/reviews/qa-finished/round-001/round.md
sha256:7144940d2294a0ebb277103e1a838ad6d5a22ae698ce0c0601244559ecaf66d9
```

Fresh path reads found the finished source Review directory absent, proving the
empty Round and Review shells were pruned. The destination retained the exact
SHA-256 identity. The live, undecidable, and Spec-owned Review Artifacts all
remained at their original paths. A fresh rerun reported `Baseline update:
current` and `File changes: 0`; it emitted the live and undecidable warnings but
no output for `qa-finished`, so the retired review was neither relocated nor
reported twice.

## Collision and partial-progress flow

A fresh disposable clone received a legacy `qa-collision` Spec with two files.
Its intended history destination already contained a different `_prd.md`, while
the sibling `task_01.md` destination was free. Their initial identities were:

- source `_prd.md`: `97493ae6…a5a7`
- occupied destination `_prd.md`: `458f1886…e4b2`
- source `task_01.md`: `bb02ca02…10a4d`

The approved update exited 3 with `Baseline update: action required`. Its error
named the exact source and destination, explained that the destination already
exists, and stated that the source was not moved. The free sibling nevertheless
moved to `docs/history/specs/qa-collision/task_01.md` with identity
`bb02ca02…10a4d`. Fresh reads proved both colliding `_prd.md` files remained at
their original paths with their distinct identities.

Fresh text and JSON plans both remained non-current and contained exactly the
one unresolved relocation. The JSON `historyMoves` row contained only
`ordinal`, `from`, `to`, and `contentIdentity`; it carried no document content.
Its state was `plan_ready`, its exit status was 3, and the text plan again
reported that repository bytes were unchanged.

## Review-source exclusion contract

Status: `blocked (environment: no open Pull Request)`. The supplied QA
environment states that this Spec has no open Pull Request, so no live CodeRabbit
review can be requested without creating external state forbidden by this gate.

Equivalent observed evidence passed. A fresh focused command,
`go test -count=1 -tags docscontract ./internal/docscontract -run
'^TestReviewHistoryConfiguration' -v`, exited 0. It exercised the positive
configuration and both negative-control suites:

- removing either `!docs/history/**` or `!docs/specs/**/reviews/**` is rejected;
- adding a rule-source glob that reaches either protected tree is rejected.

A fresh configuration read found both exclusions at `.coderabbit.yaml:215-216`.
The rule-source patterns reach neither protected tree. The broader fresh
`make verify-docs` static gate independently passed the same repository-owned
contract. Unblock action: open a Pull Request on the Spec target branch and
observe the CodeRabbit reviewed-file set; no product or QA code change is
required.

## Glossary, help, and documentation carriers

Fresh reads found `History Root` and `History Relocation` in `CONTEXT.md`; the
latter owns the emitted `HistoryRelocation`, `HistoryMove`, and `historyMoves`
spellings and defines identity-only ledger rows plus collision behavior. Positive
search found `docs/history/specs/` across the repository guides, shipped Skills,
managed module guidance, and Source Baseline corpus. A negative search for the
stale built-in spelling `_archived/specs` across those carriers returned no
matches.

Fresh built help states that `roundfix archive` sends the built-in `docs/specs`
root to `docs/history/specs/<slug>/` and retains `<spec-root>/_archived/<slug>/`
only for a configured non-default root. It also states that Archive creates no
Run and never pushes. The live archive flow and outside-evidence migration
matched that published behavior.

## Non-Goal sweep

Fresh flow and source evidence preserved all eight boundaries:

- Archive still required a completed Spec and passing QA Report; only its
  resolved destination changed.
- Every fleet, lifecycle, review, and collision identity check proved bytes move
  unread. Decision status vocabulary and recorded values were unchanged.
- The Spec-owned Review Artifact remained under the active Spec's `reviews/`
  directory. Only orphan resolution moved to the live orphan root.
- `internal/spec/review_liveness.go` imports and executes local Git plus local
  filesystem/YAML operations; it contains no hosting-provider, network, or
  credential path.
- A Spec-range path audit found no changed internal consistency-check or detector
  implementation.
- The outside repository was mutated only through a disposable clone; its source
  checkout retained its commit and legacy tree.
- The only repository configuration change was `.coderabbit.yaml`, whose changes
  are exactly the history-root and Spec-owned-review exclusion contract.

The fresh Archive, lifecycle, review, and full verification flows independently
confirmed these scope statements instead of crediting them from source review
alone.
