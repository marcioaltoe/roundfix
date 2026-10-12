---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An Archive Record names its pull request, and readers fetch the pull request head

`roundfix archive` runs on the delivery branch, before the pull request exists,
because the Baseline delivery order archives the candidate before it opens the
pull request. The record's `source_revision` is that branch's commit. A squash
merge leaves the commit out of the default branch, so after the branch is
deleted it survives only under the forge's `refs/pull/<N>/head`. On
2026-10-12 none of this repository's records for Specs 0252, 0253 and 0254
had a `source_revision` in the ancestry of `main`, and a fresh fetch of
`main` could not read their folders
(`docs/references/2026-10-12-squash-merged-archive-source-revisions.md`). An
adopter met the same on 2026-10-09 and switched its Spec pull requests to
merge commits.

On 2026-10-11 the maintainer decided "PR + refs/pull, com aviso":

- The record stores its delivery pull request as `pull_request`. Because the
  number exists only after the pull request opens, the Delivery Queue writes it
  in a **Pull Request Record Commit**: after it opens the pull request, it
  sets the field in the record that the branch's archive commit added, commits
  that one-field change and pushes it before checks are judged. An operator
  records the number with `roundfix archive <slug> --pull-request <n>`. The
  Delivery Queue proves the commit changes nothing but the record's
  `pull_request` and, when the record outgrows its target size, a shortened
  outcome paragraph. That proof stands in for a second repository gate.
- A reader that opens the folder at `source_revision` and finds the commit
  missing performs **Source Revision Recovery**. It takes the pull request
  number from the record, or, for a record without one, from the trailing
  `(#N)` of the subject of the commit that added the record. ADR-0247 already
  names that commit as the delivery evidence. The reader then runs
  `git fetch --no-tags --no-write-fetch-head <remote> refs/pull/<N>/head`
  non-interactively and within a time bound, and reads the folder again. The
  fetch writes no ref. When it fails, the reader keeps its existing
  unavailable behavior.
- `roundfix archive` warns on standard error when `source_revision` is not an
  ancestor of the default branch.
- Nothing pushes a new ref, and records already archived stay byte-identical.
  No backfill rewrites them; recovery reaches them through the subject of
  their delivery commit.

This extends ADR-0247, ADR-0247: "Readers find them as the default-branch
commit that added the record". It does not change the Archive Record schema
string: `pull_request` was already an optional field that the History Sanitize
Command writes.

## Considered Options

- Pushing a dedicated archive ref for each Spec was rejected by the
  maintainer ("No new refs pushed"). GitHub's `refs/pull/` namespace is
  read-only to pushes.
- Merging Spec pull requests with merge commits changes each repository's
  merge policy. It also fails as soon as one pull request is squashed.
- Archiving from the default branch after the merge, or opening the pull
  request before the archive commit, reverses the Baseline delivery order.
- Backfilling `pull_request` into existing records would rewrite archived
  records. The subject fallback reaches them without rewriting them.

## Consequences

A Spec's delivery pull request carries one more commit than before. Required
checks run again on it, and its review evidence covers the candidate before the
archive. Both the archive commit and the Pull Request Record Commit are
accepted by proof, not by review. A reader can now reach the network, so the
fetch is non-interactive, bounded and attempted only when the commit is
missing and a pull request number is known. Tests use a local bare repository
as the forge.
