---
status: done
created_at: 2026-07-27
updated_at: 2026-07-28
---

# Verification — sandboxed Agents cannot reach the default Go build cache, and nothing tells them so (2026-07-27)

`make verify` is the authoritative gate for this repository, but an Agent
running inside the ACP sandbox cannot always reach the default `GOCACHE`
(`~/Library/Caches/go-build` on macOS). When the sandbox denies it, the gate
fails **before compilation**, producing a failure that looks like a broken build
rather than a denied path.

Nothing in `docs/agents/go.md`, `AGENTS.md`, or the Roundfix Skill mentions
`GOCACHE`, so each Agent rediscovers the problem independently — or does not,
and reports a false gate failure.

## Evidence

- Spec 0037 QA gate, 2026-07-27: the Agent ran
  `rtk env GOCACHE=/private/tmp/roundfix-qa-0037-gocache make verify` and the
  gate passed. It had worked the problem out on its own.
- Spec 0038 QA gate rerun, 2026-07-27: a different Agent ran `make verify`
  directly and reported it "blocked before compilation by sandbox denial of the
  host Go cache." Flow QA rows were skipped as a consequence, so the verdict
  rested on an unexecuted gate.
- Every supervisor-delegated repair in this session had to pass a portable
  `GOCACHE` explicitly to get a clean run.
- The repository already carries scar tissue from the same class of problem:
  commit `f98a12f`, "fix: use portable Go cache in CLI test."

## Why this matters more than a flaky command

The gate is what decides whether a Task settles `completed` and whether QA
passes. A gate that fails for environmental reasons produces two bad outcomes
that are hard to tell apart from real ones: a Task settled `failed` whose work
was correct, and a QA verdict derived from checks that never ran. Both cost a
full recovery cycle, and both look exactly like product defects in the report.

## Suggested resolution

1. Make the cache location deterministic for Agents rather than discovered:
   have the Makefile default `GOCACHE` to a repository-local, gitignored path
   when it is unset, so `make verify` behaves identically inside and outside the
   sandbox. This needs maintainer authorization, since the Makefile and ignore
   files are protected tooling.
2. Failing that, document the requirement once in `docs/agents/go.md` and in the
   Roundfix Skill's verification guidance, so every Task and QA Agent sets it
   the same way.
3. Teach the Daemon to distinguish an environment denial from a genuine
   verification failure. A gate that never compiled should not settle a Task
   `failed` with a reason that implies the code is broken; it should surface the
   denied path and the remediation.

## Suggested acceptance checks

- `make verify` succeeds from inside a Task Worktree under the ACP sandbox with
  no environment variables set by the Agent.
- A denied cache path produces a distinct, actionable diagnostic rather than a
  compilation failure.
- Two Agents running the gate on the same commit reach the same verdict.

## What worked — keep

- Passing an explicit portable `GOCACHE` is a reliable workaround today and
  should stay documented even after the default is fixed, for anyone running the
  gate in an unusual sandbox.

## Addendum — 2026-07-28 — Routed to Spec 0054

The repository-local `GOCACHE` default in the Makefile (maintainer-authorized)
is owned by
[Spec 0054 — Tooling task and verification hygiene](../specs/0054-tooling-task-and-verification-hygiene/_prd.md).
