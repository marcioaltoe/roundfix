---
status: pending
created_at: 2026-07-29
updated_at: 2026-07-29
---

# QA gate — discovery stops at the first hard gate, so one Spec pays for five expensive rounds (2026-07-29)

Evidence from implementing Spec `0010-contrato-de-health` in the `vortex` repository, a
Bun/TypeScript monorepo. The Spec shipped correct behavior, but it needed **8 Tasks and 7 QA
rounds**. Reading the reports back, most of that cost is not defect density — it is the order
in which the gate discovers things, plus an under-provisioned gate environment. The three
findings below are what a Supervisor cannot work around from the outside.

Two frictions from the same session are already filed and are **not** repeated here:
`2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md` and
`2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md`.

## 1. A failing static gate skips every behavioral journey, serializing discovery across rounds

- Symptom / evidence: in round 1, the repository Verification failed with two reproducible
  test failures. The gate recorded `SG-01 fail` and then marked **J-01 through J-08 as
  `skipped`** with the reason `Static-gate hard stop`. The report says so directly:

  > The code-caused SG-01 failure triggered the `qa-gate` hard stop, so J-01 through J-08 are
  > skipped.

  The eight skipped rows were the live-container journeys — exactly the checks that later
  found the real defects. They surfaced one round at a time:

  | Round | Found | Would round 1 have found it? |
  | --- | --- | --- |
  | 1 | 2 stale sibling tests (static) | — |
  | 2 | `/health/visio` returns a constant; first read after outage hangs 10 s | yes, both are J-rows |
  | 3 | process never recovers after the dependency returns | yes |
  | 4 | health query rejected in production; response leaks the database address | partly |
  | 5 | production bundle drops structured log fields | yes |

  Each round is a full cycle: Agent turn, repository Verification, then the gate rebuilding
  images and starting containers. In this Spec that was roughly 15 minutes of gate time per
  round, five times, plus the corrective Task each round generated.

- Root cause: the gate treats the repository Verification as a precondition for the whole
  matrix, not as one row among others. But a static failure does not always mean the artifact
  is unusable — in round 1 the two failures were assertions in sibling test files, and the
  application still built, booted and served traffic. The gate had everything it needed to run
  the behavioral rows and chose not to.

- Action / suggestion: separate "the artifact cannot be produced" from "the repository gate is
  red". When the artifact still builds and boots, continue to the behavioral rows and report
  the static failure as its own finding. Keep the hard stop only for failures that prevent
  producing or starting the artifact. A `partial` verdict already exists to express "some rows
  could not run", so the vocabulary is there. This alone would have collapsed rounds 2, 3 and 5
  of this Spec into round 1.

## 2. The QA Agent cannot reach the Docker socket or bind a loopback port

- Symptom / evidence: rounds 6 and 7 produced **no product finding at all** and still could
  not pass, because fifteen rows were blocked by the environment:

  ```
  Docker access: `rtk docker ps -a`, `rtk docker images`, and `rtk docker network ls` each
  failed with `permission denied while trying to connect to the docker API at
  unix:///Users/marcio/.docker/run/docker.sock`.
  Loopback listener: the freshly built application failed at `Bun.serve` with `EADDRINUSE`
  on both ports `18010` and `38101`.
  ```

  The same machine, from an ordinary shell outside the Agent sandbox, ran `docker ps`
  successfully and had both ports free. The Supervisor ended up building the image and running
  the outage journeys by hand to obtain the evidence the gate could not collect.

  Earlier rounds in the same session **did** run containers, so the capability is not uniformly
  denied — one Task's `## Result` notes that a Docker build "passou após liberar o acesso ao
  Docker local", implying an interactive grant that a detached Run cannot obtain.

- Root cause: `defaults.agent_full_access: false` keeps the runtime's normal sandbox, and that
  sandbox denies the Docker socket and loopback bind. The only documented escape is
  `--agent-full-access`, which is all-or-nothing for the whole Run and for every Agent
  category.

- Action / suggestion: the QA category needs container and loopback access to do its job; the
  implementation categories usually do not. Consider a per-category capability grant (for
  example `profiles.qa.capabilities: [docker, loopback]`) so a Supervisor does not have to
  choose between a blocked gate and full machine access for every Agent in the Run. At minimum,
  detect the denial and report it as a distinct terminal reason instead of fifteen blocked rows
  that read like an incomplete matrix.

## 3. `archive` refuses a `partial` verdict with no auditable override

- Symptom / evidence: `roundfix archive 0010-contrato-de-health` refused with
  `no passing QA verdict: newest QA Report verdict is "partial"; expected "pass"`. The verdict
  was `partial` for the environment reason in finding 2 — the report itself states
  *"No product finding was confirmed on this build"*. The Spec was archived by hand: editing
  the PRD frontmatter to `status: archived` plus `archived` and `source_slug`, writing a note
  explaining the decision, and moving the folder.

  This is the second Spec in that repository archived that way; `0008-pin-de-imagem-no-cd-do-axis`
  carries the same hand-written preamble from 2026-07-27.

- Root cause: the archive precondition models the verdict as a binary gate, while `partial`
  legitimately means "no defect found, some rows unreachable" — which is exactly the state an
  under-provisioned environment produces.

- Action / suggestion: accept `partial` behind an explicit, recorded override — for example
  `roundfix archive <slug> --accept-partial --reason "<text>"` — writing the reason into the
  archive metadata. That keeps the decision auditable in the Spec instead of living in a
  hand-typed blockquote, and removes the incentive to hand-edit frontmatter that the tool
  otherwise owns.

## What worked — keep

Worth naming, because the corrections should not weaken it:

- **The gate found what nothing else could.** Five defects reached the gate that the repository
  Verification, the scoped suites and the review all passed: a health check returning a
  hard-coded `healthy`, a 10-second hang, a process that never recovered, a production-only
  query rejection, and a bundle that dropped structured log fields. Every one of them was found
  by building the production artifact and pulling the plug on its dependency.
- **The gate builds its own instruments.** In one round it created a TCP listener that accepts
  connections and never answers, purely to exercise a timeout ceiling. That is the right
  instinct and it should survive any change to round economics.
- **The Verification Feedback retry is well calibrated.** One repair turn on a deterministic
  failure, then a full re-run, settled several Tasks without a second round.
- **Task `## Result` sections refused scope creep correctly.** One Task diagnosed a stub outside
  its slice, recorded it as a follow-up and did not touch it — which is what let the Supervisor
  route it deliberately instead of discovering it inside an unrelated diff.

## Routing — 2026-08-01

Routed to [Spec 0063](../specs/0063-qa-cycle-economics/_prd.md) on 2026-08-01.
