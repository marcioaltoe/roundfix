---
type: fix
status: open
created: 2026-10-12
spec: null
---

# An Archive Record points at a branch commit that a squash merge drops

## Symptom

`roundfix archive` runs on the delivery branch, as the Baseline order
requires, and records that branch's commit as `source_revision`. A squash
merge then leaves that commit out of the default branch's ancestry. The record
says the folder is "kept in Git at `<source_revision>`". After the branch is
deleted, the commit survives only in the forge's `refs/pull/N/head`, which a
plain clone does not fetch.

- Pantheon, 2026-10-09: Spec 0043's record pointed at `b5390aab`, a commit
  that only the branch had. The maintainer switched Spec PRs to merge commits
  as a workaround.
- Roundfix itself, checked on 2026-10-12: the `source_revision` of the 0252,
  0253 and 0254 records is not an ancestor of `origin/main`.
  `c74b91f9` (0254) is reachable from `refs/pull/487/head`. The 0254 record
  carries no `pull_request` field.

## Where

`roundfix archive`, the Archive Record (`internal/spec/archive_record.go`), its readers (`history`, `review`), and the Delivery Queue's archive step.

## Expected

Maintainer decision, 2026-10-11 ("PR + refs/pull, com aviso"):

- The record always stores `pull_request`.
- Readers (`history`, `review`, anything that opens the folder at
  `source_revision`) fetch `refs/pull/N/head` when the commit is missing
  locally.
- `roundfix archive` warns when `source_revision` is not an ancestor of the
  default branch.
- No new refs are pushed. Archived records stay byte-identical.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-09-archive-record-aponta-para-commit-de-branch-sob-squash.md`
  (pantheon, 2026-10-09).
