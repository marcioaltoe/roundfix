# QA-11 — Spec 0058 replay

Status: pass

A disposable Git repository copied the complete real archived
`0058-npm-trusted-publishing-and-release-preflight` artifact tree, removed the
historical archive and override stamps, added Spec 0070's declaration overlay,
and added the provenance-bearing accepted replay report.

The built Archive Command exited 0. A fresh read found no `qa_override` and
found:

```yaml
unproven:
    - a maintainer publishes a tagged release and records the run
```

A second invocation exited 2 because the active source no longer existed.
`rtk git -c core.fsmonitor=false status --short --
docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight` was
empty in the real worktree, proving the source corpus remained untouched.
