# R09 — Pull Request environment boundary

Authoritative prompt fact: `Pull Request: none open; Pull Request journeys are
environment-blocked.` The Run Worktree branch is never pushed and cannot be
resolved as a Pull Request of its own.

Equivalent fresh evidence:

- Six focused prompt/Daemon cases passed across `internal/agent` and
  `internal/daemon`.
- The cases distinguish proven absence from unresolved lookup, resolve against
  the Spec target branch rather than the Run branch, preserve checkout facts,
  and keep the QA prompt usable when lookup fails.
- The repository loop contract orders authored QA before Pull Request creation;
  later review/approval/check/Merge-Ready behavior is unchanged by this Spec.

Unblocking action: open the Spec target branch Pull Request only after a
passing gate, then exercise approval, checks/status, unresolved threads,
Merge-Ready evidence, and review-artifact ancestry read-only.
