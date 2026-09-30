---
type: fix
status: promoted
created: 2026-09-30
spec: 0186-citation-scans-that-read-only-what-they-mean-to
reason: null
---

# The citation scan caches parent-directory checks, leaving a symlink-swap race

## Symptom

Spec 0184's Relocation Citation scan checks each tracked path's parent directories with `os.Lstat` and caches the directories it has seen. If a concurrent process replaces an already-checked tracked directory with a symbolic link, later `Lstat` and open calls under it resolve through the link. The scan could then read a file outside the repository, while the final component's `O_NOFOLLOW` and `os.SameFile` checks still pass.

## Where

`internal/baseline/history_citations.go`, the cached parent check at about line 290.

## Expected

The scan opens each path relative to a directory handle it opened itself (`openat`-style traversal with no-follow at every component), or it re-verifies the parent chain without a cache. A directory swapped mid-scan can then never redirect a read.

## Evidence

The second pre-PR review of Spec 0184's candidate `2423efe2` on 2026-09-30. It was dismissed at the review ceiling with this entry as the follow-up.
