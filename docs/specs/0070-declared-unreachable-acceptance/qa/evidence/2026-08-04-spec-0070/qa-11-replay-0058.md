# QA-11 — Spec 0058 replay

Status: pass

A disposable Git repository copied the real archived
`0058-npm-trusted-publishing-and-release-preflight` artifact tree, removed the
historical override/archive stamps, added the Spec-0070 declaration overlay,
and added the provenance-bearing accepted replay report.

`roundfix archive 0058-npm-trusted-publishing-and-release-preflight` exited 0.
A fresh read of the disposable archive found:

```yaml
unproven:
    - a maintainer publishes a tagged release and records the run
```

A second invocation exited 2 because the active source no longer existed.
`rtk git -c core.fsmonitor=false status --short --
docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight`
was empty in the real worktree, proving the source archived corpus remained
byte-untouched by the replay.
