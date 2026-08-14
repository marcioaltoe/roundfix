---
spec: 0102-a-preflight-that-proves-what-the-run-needs
status: active
created: 2026-08-12
surfaces: [backend, cli]
---

# A preflight that proves what the Run needs

Preflight proves things the Run will never use and misses things it will. Every
configured Agent Selection must prove, including fallbacks in categories a Run
never reaches, so one intermittent adapter refuses Runs whose preferred selection
is perfect. Meanwhile a declared worktree source file that does not exist is
discovered per worktree, after the Run exists — a checkout observed on 2026-08-08
would have failed bootstrap in every Task Worktree and died without producing a
line of code, and it was caught only because a maintainer checked by hand. When a
proof does fail, the surfaced next action can point at a cleanup problem that does
not exist rather than at the model identifier that is actually wrong.

## Project Constraints

- Identifier strategy: applicable — Agent Selection, Preferred Selection, Fallback
  Selection, Agent Session, and Task Worktree are glossary terms this Spec changes
  the proof obligations of. The closing node checks whether the work introduced or
  changed a term the glossary should carry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read by this Spec. Adapter sessions are opened by the
  existing runtime through its own configuration, which this Spec does not change.
  Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0119 makes the refusal that fired
  first the refusal, which governs how a failed proof is reported and is exactly
  the principle a spurious cleanup error violates; ADR-0111 makes an unobserved
  Verification unknown rather than a verdict, which is the discipline a lazily
  proven fallback must not break — an unproven fallback is unknown, never proven.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the preflight, agent, and worktree
  packages plus their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Run is refused for a reason that would have stopped it, not for one that
   would not.
2. A declared file the Run needs is proven present before the Run exists.
3. A failed proof surfaces one next action, and it is the one that fixes the
   cause.
4. An unproven fallback is reported as unproven, never as proven.

## User Stories

1. As a Supervisor whose preferred selection is healthy, I want a Run to start
   when only an unreached fallback's adapter is intermittent, so that a runtime I
   am not using cannot block the work.
2. As a Supervisor with a declared worktree source file missing from my checkout,
   I want preflight to name it, so that I fix it before a Run instead of losing
   every Task Worktree to bootstrap.
3. As a maintainer diagnosing a failed selection proof, I want the printed next
   action to name the wrong model identifier, so that I am not sent to repair a
   session that was never opened.

## Core Features

1. **Declared worktree sources are proven.** Each declared repository-relative
   source copied into a Run Worktree is checked for existence during preflight,
   and a missing one refuses by name before any Run is created.
2. **The preferred selection always proves; fallbacks do not block.** The
   preferred Agent Selection of every used category is proven as today. A fallback
   is proven without blocking the Run when it fails, and its state is reported as
   unproven rather than assumed good.
3. **A cleanup failure never outranks its cause.** A disposable Agent Session
   cleanup error that only reports a session absent, when the setup error already
   explains why no session exists, is not joined into the surfaced error and never
   supplies the printed next action.
4. **A starved proof is unknown, not unavailable.** A proof that exceeded its
   deadline while the host was saturated is reported as undetermined rather than
   as a rejected Selection, so an operator reads a machine under load rather than
   a broken configuration. Measured across three occasions: at load average 23 the
   refusal widened from two categories to six, and at 4.9 every tuple passed on the
   first attempt with nothing else changed.
5. **One next action per refusal.** A failed readiness check surfaces exactly one
   deterministic next action, and it addresses the cause the check found.

## User Experience

A preflight refusal names the missing declared file, or the failed selection with
its category and the advertised set, and prints one next action. A Run whose
fallback could not be proven starts, and reports that fallback as unproven in the
same place it reports the selections it did prove.

## Non-Goals / Out of Scope

- Changing which selections a category may configure, or the fallback chain shape.
- Changing adapter capability discovery or the bounds that refuse a malformed
  capability payload.
- Bootstrap execution itself; this Spec proves the declared inputs exist and does
  not change what bootstrap does with them.
- Any change to how a Run selects among proven selections at dispatch.

## Success Metrics

- A Run starts with a healthy preferred selection while an unreached fallback's
  adapter is failing, measured against a configuration of five distinct selections
  in a repository this Spec did not build, where an intermittent adapter refused
  several Runs on 2026-08-08.
- A missing declared worktree source refuses before a Run is created and names the
  file.
- A failed selection proof prints one next action, and it names the model
  identifier rather than session cleanup.
- No refusal surfaces a cleanup error whose cause is the setup error beside it.

## Decisions

- An unproven fallback is reported as unproven rather than silently assumed, so
  the report never claims evidence the preflight did not gather.

## Open Questions

- Whether fallback proving is lazy by default, tolerant by default, or configured.
  Tolerant with an explicit unproven report is the default until answered, because
  it never blocks a Run on a selection the Run will not use, while keeping the
  absence of evidence visible.
- Whether a declared worktree source that exists but is unreadable is the same
  refusal as one that is absent. The default is to refuse both, naming which.
