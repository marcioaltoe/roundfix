---
spec: 0055-owner-identity-without-fork
status: archived
created: 2026-07-28
surfaces: [backend, cli]
archived: "2026-07-31"
source_slug: 0055-owner-identity-without-fork
---


# Owner identity without fork

Spec 0037 gave Force Stop a real ownership proof: Runs record an opaque owner
start-time identity, and Force Stop refuses to signal a PID whose live
identity does not match. The proof is obtained by forking `/usr/bin/ps` on
every read, so it fails exactly when the host cannot fork — which is
precisely when a runaway Run has exhausted the machine and the escape hatch
is needed most. Under the same pressure, a failed capture at Run creation
silently records no identity, so a Run can run without reuse protection with
no warning. Separately, `roundfix stop <run-id> --force` rejects its trailing
flag — the same argument-ordering defect Spec 0042 fixed for the Attach
Command. Evidence:
[owner identity forks ps and fails closed under load](../../findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md);
the Stop Command flag defect was observed during the Spec 0042 QA recovery on
2026-07-28, alongside the Attach fix.

## Project Constraints

- Identifier strategy: not applicable — the owner identity remains an opaque
  equality-compared token; no project-owned Internal Identifier is created.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all behavior is local process
  inspection and CLI argument parsing; no authentication or HTTP surface.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0044 requires orphaned Run locks
  to be reclaimed only on proven owner death, which the new identity source
  must keep proving; ADR-0052 protects compare-and-set terminal completion;
  the Spec 0037 ownership-proof semantics (refuse on proven mismatch) are
  preserved. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Owner identity is captured and compared without spawning any subprocess,
  so Force Stop works when the host cannot fork.
- An unreadable identity and a proven identity mismatch are distinct
  conditions with distinct diagnostics; only the mismatch keeps refusing.
- A Run running without reuse protection is never silent.
- `roundfix stop <run-id> --force` parses like every other command.

## User Stories

1. As an operator force-stopping a runaway Run on a loaded host, I want the
   ownership proof to read process identity without forking, so that the
   escape hatch does not fail with the failure it exists to escape.
2. As an operator told a Run's owner cannot be proven, I want the diagnostic
   to say whether the identity was unreadable or genuinely mismatched, so
   that I am not sent to investigate PID reuse when the machine simply could
   not answer.
3. As a user whose Run was created while identity capture failed, I want a
   startup warning and a durable marker, so that a Run running with
   PID-only protection is observable.
4. As an operator with a dead host and an unreadable identity, I want a
   documented path to stop the Run anyway, so that failing closed on
   ignorance has a supervised exit.
5. As a user typing `roundfix stop <run-id> --force`, I want the trailing
   flag accepted, so that the documented escape hatch works as printed.

## Core Features

1. Owner start-time identity comes from a direct kernel read on both
   supported platforms — procfs on Linux, sysctl on macOS — behind the
   existing platform structure; the token stays opaque and
   equality-compared, and non-Unix stubs keep their behavior.
2. Failure classification separates `identity mismatch` (proven different
   owner — keeps failing closed exactly as today) from
   `identity unreadable` (host could not answer — its own diagnostic naming
   the resource failure and next action).
3. A failed identity capture at Run creation emits one warning and records
   a durable marker on the Run; the silent NULL degradation remains only for
   genuinely legacy rows predating the identity column.
4. A documented, supervised operator path exists for stopping a Run whose
   identity is unreadable; it is explicit, never the default, and preserves
   the proven-mismatch refusal.
5. The Stop Command accepts its flags in any position relative to the Run
   ID, matching the Attach Command's fixed parsing.
6. All existing proofs keep their behavior: reuse refusal, matching
   identity, absent-owner-as-proof, and legacy no-token degradation.

## User Experience

- A Force Stop refusal names which condition it hit: proven mismatch, or
  unreadable identity with the host error and the documented next action.
- A Run created without reuse protection prints one startup warning and is
  visibly marked in Run inspection output.
- `roundfix stop run_x --force` and `roundfix stop --force run_x` behave
  identically.

## Non-Goals / Out of Scope

- Weakening the proven-mismatch refusal or the compare-and-set terminal
  contract.
- Windows-native identity capture beyond the existing stub behavior.
- Changing Stop Request semantics, force-stop signaling order, or the Agent
  Session registry owned by Spec 0037.
- Removing the legacy PID-only degradation for rows that predate the
  identity column.

## Success Metrics

- Owner identity capture and comparison spawn no subprocess, proven by a
  test that fails if the implementation execs anything.
- Under simulated fork exhaustion, Force Stop distinguishes unreadable from
  mismatched, each with its own diagnostic; the full suite passes under
  parallel load with no `ps`-related failures.
- A Run whose identity capture failed at creation emits the warning and the
  marker is queryable.
- `roundfix stop <run-id> --force` succeeds with the flag in trailing
  position.

## Decisions

- Failing closed on a proven mismatch is preserved; unreadable identity
  fails closed by default but gains a documented supervised path, because an
  escape hatch that cannot be opened under the exact failure it exists for
  is not an escape hatch.
- The token stays opaque: no parsing, no format guarantee, equality only.

## Open Questions

None.
