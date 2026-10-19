---
status: accepted
created_at: 2026-10-19T20:00:00Z
updated_at: 2026-10-19T20:00:00Z
deprecated_at: null
superseded_by: null
---

# Reachability reads the fetched default branch, and origin names it when origin/HEAD is unset

On 2026-10-18 the Verification gate of pull request #580 failed in GitHub
Actions because `SC-HISTORY-RETAINED` called a Git Locator in Spec 0279's
Archive Record invalid, although its commit was on `main`. The CI checkout
fetches every branch into `refs/remotes/origin/*` and every tag, checks out the
pull request head detached, and leaves `refs/remotes/origin/HEAD` unset; the
actions/checkout issue #910 ("`origin/HEAD` Not Present", open since 2022)
records the same. Roundfix's default-branch detection then answered
"undetermined", so only tags proved reachability. The same day a scratch clone
with a stale local `main` called a valid locator invalid, because the Git
Locator read `refs/heads/<default>` before `refs/remotes/origin/<default>`.
On 2026-10-19 the maintainer asked for one Spec that proves reachability the
same way locally and in CI and never trusts a stale local branch over the
fetched remote.

## Decision

**origin names the default when origin/HEAD is unset.** Default-branch
detection keeps its order: `refs/remotes/origin/HEAD` when set, then a current
branch named `main` or `master`. When both fail, exactly one fetched
`refs/remotes/origin/main` or `refs/remotes/origin/master` names the default,
with the detection source `origin-branch`. Both present, or neither, leaves the
default undetermined; Roundfix never guesses between them.

**One Default Branch Head serves every reachability and merge-base read.** Its
commit is the fetched `refs/remotes/origin/<default>` when that ref exists, and
the local `refs/heads/<default>` only when no remote-tracking ref exists. The
Git Locator's reachability and path history, the QA Delivery Base, the Archive
Command's `source_revision` warning and its record-only pull request form, and
the pre-PR review's default base all read it. A local default branch that
holds unpushed commits no longer proves reachability, because CI would not
see them either.

## Considered Options

- Running `git remote set-head origin --auto` was rejected: it contacts the
  remote and writes a ref, and a check must stay read-only and offline.
- Treating any remote-tracking branch that contains the commit as proof was
  rejected: a live work branch would keep a commit alive that the default
  branch and the tags do not hold, which ADR-0278 excludes.
- Adding the remote-tracking ref only as a second candidate after the local
  branch was rejected: a stale local branch still answers first, which is the
  defect.
- Preferring `main` when both `origin/main` and `origin/master` exist was
  rejected: a repository that migrated keeps both, and a wrong default turns
  a missing proof into a wrong one.
- Reading a CI environment variable such as `GITHUB_BASE_REF` was rejected: the
  check must give the same answer on every machine from the same refs.

## Consequences

- A record archived between two releases passes the Spec Consistency Check in
  a CI checkout, and a stale local `main` no longer hides a fetched commit.
- Commands that already read the fetched remote after a fetch, such as the
  Delivery Queue and the Linux Lane, keep their reads; they now also find a
  default when origin/HEAD is unset, instead of refusing.
- Preflight Validation's question "is this branch the default?" is unchanged
  for a named branch; a detached HEAD still matches no default.
- This refines ADR-0278, whose Git Locator names a commit "the default branch
  or a tag reaches", by defining that branch's head, and ADR-0276, whose
  Delivery Base is the merge base with the default branch.
