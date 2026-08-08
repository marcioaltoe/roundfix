---
date: 2026-08-08
supersedes: 2026-08-05-the-night-the-queue-moved.md
---

# Handoff — two quotas out, and a gate that said no

Written to end a session that ran 1 day 7 hours and spent $284 of Anthropic
allowance on its own context. Nothing is running: no Roundfix Run is active
anywhere on the machine, no background task, no scheduled agent job. Everything
below is committed on `ma/specs-0082-0083`, nineteen commits ahead of its
upstream, nothing pushed.

## The single most important fact

**Both implementation quotas are out or nearly out.**

- **Codex**: exhausted, returns 2026-08-12 17:43. Every Run attempted after
  17:43 on 2026-08-08 died in nineteen seconds; the usage-limit text arrives
  wrapped as `agent/protocol error`.
- **Anthropic**: 66% of the week used, resets **2026-08-10 14:00**
  (America/Sao_Paulo).
- **OpenCode Go**: freshly subscribed, all three windows at 0% — continuous
  (5h), weekly, monthly. Unusable from Roundfix today; see the next section.

The next session's whole purpose is to make the third one usable, because it is
the only runtime with capacity.

## The work that is queued and why it is the priority

Two problems stand between Roundfix and the OpenCode runtime. Both were measured
on 2026-08-08; neither is diagnosed.

**1. A configured `opencode` profile does not register.** Adding a `data`
category with `runtime: opencode`, `model: opencode-go/kimi-k3` to
`.roundfixrc.yml` and running `roundfix doctor` reported the same
`5 distinct tuples; 10 category references` as before the edit — as if the
profile did not exist. `opencode` is a supported runtime in the config parser
(`internal/config/profiles.go:420`) and in the adapter table
(`internal/agent/acpx_runner.go:60`), so the silence needs explaining before
anything else. The probe was reverted; the file is clean.

**2. The advertised catalog exceeds the capability cap.** `opencode models`
returns **417** entries:

| provider prefix | count |
| --- | --- |
| `openrouter` | 339 |
| `opencode` | 60 |
| `opencode-go` | 18 |

`internal/agent/selection_capabilities.go:35` sets `maxCapabilityValues = 64`,
and line 446 discards the entire model capability above it rather than
truncating. No capability means no Exact Agent Selection Proof, which means no
Run. Note that even the good case is fragile: the `opencode` provider alone is
60, four short of the cap.

This is filed as
`docs/backlog/2026-08-08-a-runtime-that-advertises-a-catalog-is-not-unusable.md`.
Its proposed shape is to bound what Roundfix *keeps* rather than what a runtime
may *advertise*.

Do not assume the fix is only the cap. Problem 1 may be the real blocker and may
have nothing to do with problem 2.

### What the OpenCode Go subscription actually contains

Eighteen models under the `opencode-go/` prefix. Two carry a **2x usage**
multiplier and are the wrong default: `gpt-5.6-luna` and
`deepseek-v4-flash`. `deepseek-v4-pro` carries no multiplier, which inverts the
obvious guess. The `OpenCode Zen` entries in the model picker — `gpt-5.6-sol`,
`gpt-5.5-pro`, `grok-build-0.1` — are a separate pay-per-use tier, not part of
the subscription.

Two account toggles are ON and were left ON deliberately, as the maintainer's
call: spending balance after limits are reached, and China-hosted providers. The
second matters for client repositories (fiscus, vortex handle fiscal data) more
than for this one.

OpenRouter is alive and cheap — Hermes ran GPT-5.6 Luna through it as recently
as 2026-08-08 at about $0.002 per call. An earlier note in this repository's
session history that the OpenRouter key might be stale is **wrong**; the logs
disprove it.

## Spec 0084 — where it stopped

`0084-an-update-that-can-run` replaces Spec 0082's managed-marker blocking
requirement, which made `roundfix baseline update` refuse six of the eight
repositories that had adopted a Baseline. Nine of its ten Tasks are completed
and integrated. The tenth is the authored QA gate.

**The QA gate ran and returned `fail`, with one finding, and it was right.**

F1: `make verify` was red because `docs/references/coverage-record.json` still
named `TestManagedRefreshBlocksHandEditedManagedMarker`, a test Task 01 removed
on purpose — it encoded the very requirement this Spec supersedes. The coverage
guard read a deliberate retirement as a regression.

**The repair is already applied**, using the flag the QA named:
`go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record`.
After `go clean -testcache`, `make verify` exits 0.

**What remains for 0084**: re-run the authored QA gate so its verdict flips.
ADR-0097 allows carrying rows forward on declared unmoved evidence, so only row
1 needs re-execution. Task 10's file still reads `status: failed` and the report
still reads `verdict: fail`; both are accurate until the gate runs again.

A caution for whoever re-runs it: **`make verify` reported exit 0 on a stale
test cache** while the underlying test failed. Run `go clean -testcache` first,
or run the specific test with `-count=1`. This is a variant of the defect Spec
0083 exists to close and it is not closed.

## What Spec 0084 delivered

- A managed region whose bytes are not the recorded ones is *classified*, not
  presumed damaged; only a duplicated managed identity still blocks.
- The presented plan names every unrecorded region and every line the refresh
  removes, inside the Plan Digest.
- The update converges: an applied refresh republishes the Setup Manifest, and
  the second run reports the repository current with zero changes.
- The fourteen structural clauses the catalog stopped emitting are emitted
  again.
- An unresolved Baseline Profile reports its identity, the searched locations,
  and the repair action instead of a raw filesystem error.
- Three Normative Clauses seated: the Secondbrain is a consultation source when
  a decision is formed, a Spec's acceptance rests on at least one row of
  evidence it did not author, and the glossary is checked when work closes.
- `CONTEXT.md` gained Managed Region, Managed Refresh, Unrecorded Managed
  Region, and External Acceptance Evidence.
- ADR-0101 through ADR-0104.

## The pattern worth changing

Four Runs were needed. Three failed, and **every failure was a task-authoring
defect, not an Agent defect**:

1. Task 03 required a change to `.agents/skills/roundfix/SKILL.md` that no
   authorization named. The Agent recorded the blocker and changed nothing —
   correct behavior, missing authority.
2. Task 06 broke a maintained fixture the Task never enumerated.
3. Told to move that fixture, the Agent moved it, and the frozen parity corpus
   refused: `repointing its blob and fixtures was tried on 2026-07-30 and
   reverted`. **The instruction was wrong.** The sanctioned escape was written
   inside the test all along — declare the grown guide in the
   `evolvedPastFrozenCorpus` set and let `TestFormatterComposition` cover its
   bytes. Task 07 then hit the same wall for `docs/agents/spec-routing.md`.

The shape repeats: the boundary was described instead of enumerated. ADR-0104
and the outside-evidence clause exist because of exactly this, and they were
seated by this Spec — after the fact, by three of its own Runs.

There is also a self-inflicted diagnostic hole. Every Task's Verification uses
`go test … > log 2>&1 && grep -q … log`, adopted to stop pipelines from hiding
exit status. It works, and it leaves the Daemon's captured diagnostics **empty**,
so the failure feedback the Agent receives carries no cause. Two Agents burned
their retry deducing what a one-line message would have told them. Captured in
the Secondbrain inbox, still pending triage.

## Filed on 2026-08-08

Backlog entries, all `open`:

- `2026-08-08-record-usage-per-agent-session.md` (`feat`) — Roundfix persists
  which tuple ran each Task and records nothing about what it consumed, so
  concurrency and reasoning effort cannot be told apart. The question that
  motivated it could only be answered by inference over billing aggregates.
- `2026-08-08-a-session-that-never-opened-is-a-selection-failure.md` (`fix`) —
  the Fallback Chain did not activate on quota exhaustion, because Roundfix
  emitted `AGENT_WORK_STARTED` before the adapter exited and then classified the
  exit as a failure after work began. A one-Run override also cannot express a
  fallback, so working around it required an artificially distinct tuple.
- `2026-08-08-a-runtime-that-advertises-a-catalog-is-not-unusable.md` (`fix`) —
  the capability cap described above.

Secondbrain inbox, pending triage:
`inbox/roundfix/2026-08-08-verification-que-redireciona-esconde-o-diagnostico-do-agente.md`.

Adopted by Spec 0084 as its outside evidence:
`docs/specs/0084-an-update-that-can-run/references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`,
the read-only fleet measurement across the eight repositories that carry a Setup
Manifest.

## Open decisions the maintainer has not answered

1. **Re-run the QA gate now or after the 2026-08-10 reset.** Nothing else in
   0084 depends on it; the code is integrated and the local gate is green.
2. **Whether to merge PR #143 before or after 0084 closes.** The maintainer
   chose "fix before merge" on 2026-08-08, before it was established that
   `release.yml` fires only on a `v*` tag or manual dispatch — merging to main
   publishes nothing. That fact may or may not change the choice; it is theirs.
3. **Whether to open the OpenCode Spec now.** Asked twice, unanswered.

## Queue after this

| # | Spec | State |
| --- | --- | --- |
| — | OpenCode runtime reachable | not yet minted; the priority |
| 1 | 0084 | 9/10 done, QA gate to re-run |
| 2 | 0080 | 8 Tasks authored, pending; needs one outside-evidence row added |
| 3 | 0081 | 9 Tasks authored, pending; same |
| 4 | 0085 | baseline reform, remainder |
| 5 | 0086 | external skill restoration |
| 6 | 0087 | separate acceptance authorship; may be unnecessary after ADR-0104 |

0080 and 0081 were authored before the outside-evidence clause existed. They are
mechanically clean under `roundfix spec check`; the clause is not enforced by any
`SC-*` code, which is itself a gap worth naming in 0085.

## State to reproduce

```
branch: ma/specs-0082-0083   (19 commits ahead, unpushed)
PR:     #143 OPEN — "feat: a Baseline update that reads the manifest, and a gate that can say no"
gate:   make verify exits 0 after go clean -testcache
spec:   roundfix spec check 0084-an-update-that-can-run → no findings
runs:   none active
```

The Secondbrain has two commits from this session that are pushed to it, not to
this repository: the inbox captures and their triage.
