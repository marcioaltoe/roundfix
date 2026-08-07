---
status: pending
created_at: 2026-08-07
updated_at: 2026-08-07
kind: finding
spec: 0082-the-manifest-already-answered-that
---

# The setup refresh interviews a repository that already answered (2026-08-07)

Refreshing an adopted repository's Context-Driven Baseline costs a full
interview plus a supervised ACP turn, on every repository, on every release.
Measured across the fleet on 2026-08-07: nine adopted repositories carry six
distinct catalog digests, and only four sit on the digest the current binary
generates.

## What was observed

**The interactive command re-asks what the manifest holds.** Driven on a
throwaway clone of this repository, `roundfix baseline` announced
`Baseline workflow: update`, named the current profile, and then issued twelve
prompts — every one shaped `Keep <id>=<stored value> (default) / Change <id>`.
The manifest is read only to seed prompt defaults; `promptBaselineDecisions` is
called unconditionally for every profile decision.

**The classification step spawns an ACP runtime.** After profile alignment the
run went silent. A `SIGQUIT` goroutine dump located it exactly:

```
cli.promptBaselineClassification
  → baselineacp.(*Analyzer).Segment
    → baselineacp.(*Analyzer).attemptSegmentation
      → agent.(*ACPXRunner).RunSealedPrompt
        → agent.(*ACPXRunner).runSealedPromptCommand
          → syscall.wait4
```

`classificationRuntime` hard-codes runtime `codex`.

**The stall's cause is not established.** Two facts sit in tension and neither
was resolved here. On this machine `which -a codex` resolves first to a cmux
shim — a bash script under `/var/folders/.../T/cmux-cli-shims/<uuid>/codex` —
ahead of the real binaries at `~/.local/bin/codex` and `/opt/homebrew/bin/codex`.
But `roundfix doctor` exits `0` and reports `codex: ok`, naming
`/Users/marcio/.local/bin/codex`, and reports the `codex-acp` adapter healthy at
1.1.13. So readiness diagnostics and the sealed-prompt runner may resolve the
same runtime differently, or the stall may have nothing to do with resolution at
all. Do not act on the shim as a cause without first proving which binary the
sealed prompt actually spawns.

**The 5-minute bound did not bound it.** `agent.SealedPromptTimeout` is
`5 * time.Minute` and reaches the subprocess through `exec.CommandContext`. A
run left alone was still alive past ten minutes. The mechanism is not proven; a
pipe write-end inherited by a descendant of the killed child, which would block
`Wait` after the direct child dies, is the leading hypothesis and needs a
dedicated repro before anyone acts on it. This one is independent of the
resolution question above: whatever the sealed prompt spawns, its declared
timeout did not end the run.

**The non-interactive path never opens the manifest.** With the manifest
present and complete, `roundfix baseline plan --profile go-cli-tui` exits `3`
naming all eight required decisions. Supplying all ten decisions from that same
manifest exits `3` again: `instruction-preservation mode is required`. There is
no preservation flag on `baseline plan`; the only route is a strict Decision
Document that, under Preservation, must bind every Source Baseline Entry.

**A moved profile digest downgrades an update to an adoption.**
`inspectBaselineHumanState` assigns the resolved profile and every stored
decision *before* it compares the profile digest, then returns on a mismatch
without setting update mode. A repository observed on 2026-08-07 therefore
announced `Baseline workflow: adoption / Existing state: incompatible — the
existing Setup Manifest references an unavailable or changed Baseline Profile`
and immediately offered `Reuse existing profile <stored> (default)` followed by
every stored value as a prompt default. It asked twenty-one prompts. The same
message covers two different causes — an unresolvable profile and a merely
moved digest — so neither can be treated differently until they are separated.

**A profile change restarts the interview.** In that same run, selecting a
replacement profile at prompt 18 re-asked at prompts 19 onward the decisions
already answered at prompts 3 onward, including ones both profiles share.

## Why it matters

The cost is why the fleet drifts. A refresh that should be mechanical costs
twelve to twenty-one confirmations and a model turn that can stall past its own
timeout, so repositories are left stale and agents run against guidance that no
longer matches the binary.

## Route

Spec 0082 covers the manifest-driven update, the managed-refresh preservation
mode that removes the analyzer from the refresh path, and the digest-drift
correction. It does **not** cover two of the observations above, which remain
live and unowned:

- the sealed prompt outliving `SealedPromptTimeout`, which also affects first
  adoption and every other sealed-prompt caller;
- which binary the sealed prompt actually spawns for runtime `codex`, given that
  `doctor` and `which -a` disagree on this machine;
- the profile-change interview restart, recorded as the open question in Spec
  0082's PRD.
