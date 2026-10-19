---
type: fix
status: promoted
created: 2026-10-18
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
reason: null
---

# A Git Locator in a freshly archived record fails the check in CI

## Symptom

On 2026-10-18 the Verification gate of PR #580 failed
`TestCheckActiveCorpusHasNoErrors`: `SC-HISTORY-RETAINED` called the Git
Locator `docs/specs/0279-…@b29c5430` in Spec 0279's Archive Record invalid.
That commit is on `main`. The GitHub Actions checkout has no
`refs/remotes/origin/HEAD` and no local `main`, so the default branch is
undetermined, and only tags prove reachability. `b29c5430` was newer than
the last tag (v0.89.0), and the run passed once the v0.90.0 tag existed.

A related case showed up in a scratch clone: a stale local `main` made a
valid locator look invalid, because the default head is read from
`refs/heads/<name>` before `refs/remotes/origin/<name>`.

## Where

`internal/speccheck/history_retained.go` (`invalidHistoryLocators`) and
`internal/spec` `ResolveGitLocator`, through `preflight.DetectDefaultBranch`.

## Expected

The check proves reachability the same way locally and in CI, for example by
also trying `refs/remotes/origin/main` or the remote's HEAD when the name is
undetermined. It does not trust a stale local branch over the fetched
remote. A record archived between two releases then passes the check.
