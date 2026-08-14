---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md
line: 111
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiAoa,comment:PRRC_kwDOS0qyts7fC8Qq
review_hash: adfa9e18d2ee6c906e495637739718aa4aaffc40bbc9fb51718392b7b668ffd0
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:48Z"
---

# Issue 030: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Use the cache that `make verify` actually controls.**

This handoff says bare `go clean -testcache` made the gate trustworthy. The accompanying backlog states that `make verify` exports the repository-local `.gocache`, so bare cleanup leaves the gate's cache intact. A reader can repeat the stale-cache false green. Replace the command with `GOCACHE="$PWD/.gocache" go clean -testcache` or an authorized target, and update both statements.

Based on the supplied cache-analysis backlog, `make verify` uses a repository-local cache that bare cleanup does not invalidate.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md` around
lines 99 - 111, Update the handoff’s cache-cleaning guidance to invalidate the
repository-local cache controlled by make verify, using GOCACHE="$PWD/.gocache"
go clean -testcache or an authorized equivalent. Replace both references to bare
go clean -testcache while preserving the instruction to rerun the specific test
uncached when appropriate.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e976aa971992453589fa7690 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Updated both references in the handoff file (lines 101 and 109). The cache-cleaning command now reads `GOCACHE="$PWD/.gocache" go clean -testcache`, matching the Makefile's local cache setup (`GOCACHE ?= $(CURDIR)/.gocache`). This ensures the handoff's guidance invalidates the same cache that `make verify` uses, as documented in the accompanying backlog.
