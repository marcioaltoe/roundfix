---
status: accepted
created_at: 2026-10-20T21:40:00Z
updated_at: 2026-10-20T21:40:00Z
deprecated_at: null
superseded_by: null
---

# A queue commit names the Spec by number, and a hook refusal parks with the hand commit

On 2026-10-20 an adopter's Delivery Queue parked a Spec with a generic
`delivery-error` at the Pull Request Record Commit. The subject the queue
wrote, `docs: record pull request #280 in the <slug> Archive Record`, was 107
characters long. The adopter's `commit-msg` hook runs commitlint with the
conventional configuration, whose `header-max-length` is 100, so Git refused
the commit and a retry repeated the same subject. The operator committed the
staged change by hand with a short subject and retried, and the queue adopted
that commit, because it proves the record change and never reads the message.
The maintainer chose "Número + falha nomeada (Recomendado)": the subject names
the Spec by its number, and a commit a hook still refuses parks with a reason
that names the hook and the hand commit instead of a generic error.

## Decision

**The Pull Request Record Commit names the Spec by number.** Its subject is
`docs: record pull request #<number> in the <spec number> Archive Record`. The
Spec number is the run of ASCII digits that starts the slug, when one to ten
digits are followed by `-`. A slug without such a prefix gives
`docs: record pull request #<number> in the Archive Record`. Either form is at
most 70 characters for any legal pull request number, whatever the slug. The
queue still adopts an existing commit by proving its change against the
candidate, so a record commit made by hand under any subject, including the
earlier long one, keeps being accepted.

**A queue commit a hook refuses parks with the hook and the hand commit.**
When a `git commit` the Delivery Queue makes fails and Git's own output names a
commit hook, or the repository holds an executable `pre-commit` or `commit-msg`
hook, the failure is a hook refusal. This covers the Pull Request Record Commit,
the archive commit and the derived merge commit. The queue reuses the Daemon's
classification of a Task commit refusal. The park keeps its `delivery-error`
blocker, its `environment` Park Class and the stage it parks at. The reason
after the blocker names the refusing hook, the refused subject, the hook's exit
status and its output, and the recovery: commit the staged change by hand with
a subject the hook accepts, then run `roundfix deliver retry <slug>`. For the
derived merge, whose merge is aborted, the recovery is the hand merge a
conflict park already names. The queue never passes `--no-verify` and never
rewrites a subject on its own.

## Considered Options

- Cutting the slug to fit 100 characters was rejected: a cut slug names no
  Spec reliably, and the number already identifies it within the repository.
- Retrying the commit with `--no-verify` was rejected: it bypasses a check the
  repository chose, which the Baseline forbids.
- A new blocker for hook refusals was rejected: the class and the next action
  of a `delivery-error` already fit, and a new blocker would change every
  reader of the Park Class table for a reason text change.
- Shortening the archive subject and the pull request title in the same change
  was rejected for now: the maintainer's decision covered the record subject
  only. Both still grow with the slug, and an open Backlog Entry tracks them.

## Consequences

- The structural signal the Daemon accepted carries over: a genuine Git
  failure, such as an index lock, in a repository that holds a commit hook is
  reported as a hook refusal. The recovery it names commits by hand, where the
  real failure surfaces again.
- A repository with no executable commit hook keeps the raw Git error.
- The Archive Record, its `pull_request` field and every record already
  written stay as they are; only the message of new record commits changes.
