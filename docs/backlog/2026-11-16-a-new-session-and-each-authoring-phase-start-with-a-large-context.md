---
type: perf
status: open
created: 2026-11-16
spec: null
reason: null
---

# A new Supervisor session and each Authoring Phase start with a large context

## Slow

Every call of a Supervisor session re-reads what the session loaded when it started, and every Authoring Phase reads the whole glossary and the agent guides before it reads anything about its Spec. Whoever supervises Spec work from an Agent harness with a usage limit pays for it on each call. Spec 0397 tells the Supervisor to end its session at each phase boundary, so a session starts more often and this floor is paid more often.

The maintainer sent this out of the Idea of Spec 0397 on the night of 2026-11-15 to 2026-11-16, choosing "Original + respostas (Recommended)": "Registro de fase, fronteira de sessão e adoção dos modelos numa Spec; o custo de nascer de cada sessão vira Backlog Entry."

## Measured

The author of Spec 0397 measured it on the Mac mini on 2026-11-15, from the calls the operator's usage script collected for 2026-11-08 to 2026-11-15. The record is [the measurement adopted by Spec 0397](../specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/references/2026-11-15-usage-by-context-and-session-identity-measurement.md), section 5.

- The first call of the seven main sessions of this repository carried 82,386, 85,884, 87,692, 89,155, 90,244, 90,246 and 150,494 tokens: input plus cache read plus cache write.
- The first call of the 38 main sessions opened in a repository checkout on the host carried a median of 75,822 tokens and a first quartile of 51,541.
- By file size at about 3.5 bytes per token, an estimate and not a count, an author reads before anything about its Spec: the glossary, 182,265 bytes, about 52 thousand tokens; the agent guides, 78,826 bytes, about 22 thousand; and its authoring skill, about 4 to 13 thousand.

Not measured: how the 82 to 90 thousand tokens divide between what Roundfix writes into a repository (the root instruction file, the guides it names and the installed skills a harness lists) and the user's own configuration (personal instructions, memory, and tool servers with their tool lists).

## Target

- The first call of a new Supervisor session of this repository carries under 60 thousand tokens. A quarter of the host's sessions already start below 52 thousand.
- The reads an Authoring Phase makes before anything about its Spec stay under 40 thousand tokens by the same estimate, against about 78 thousand today.

Both numbers are proposals. The Spec that adopts this entry fixes them after it measures which part of the floor Roundfix writes, and reads the effect by phase from the record Spec 0397 delivers.
