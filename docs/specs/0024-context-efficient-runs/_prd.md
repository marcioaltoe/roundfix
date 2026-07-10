---
spec: 0024-context-efficient-runs
status: active
created: 2026-07-10
surfaces: [cli, backend, data, infra, docs]
---

# Context-Efficient Runs

A measured phase-09 single-Task Run produced a 400 KB Console Log, approximately 100,000 tokens of text, including 31 full unified-diff JSON payloads and 330 echoed file reads. A replacement Agent Session spent approximately 300,000 tokens re-establishing Spec and repository context. Supervisors also had to grep that verbose Console Log for state changes, which caused a false failure report from an unrelated `Failed:` test-runner line. Roundfix must keep lossless evidence while giving Agents, Supervisors, and Detached Run callers only the information each needs.

## Goals

- Verification output reaches the Agent only when it can change the Agent's next action, and one failed verdict gets one bounded repair attempt.
- A Supervisor can replay or follow one Run through four compact JSONL categories without parsing Console Log text.
- The Console Log remains useful as the Detached Run record without duplicating file contents or unified diffs.
- Each Spec Task starts with its complete Task contract and bounded references to the larger context needed for focused work.

## User Stories

1. As a developer running Agents over Work Items, I want the Daemon to return only Verification failures, so that successful test output does not consume Agent context.
2. As a developer whose initial Verification fails, I want the same Agent Session to receive concise diagnostics and get one repair attempt, so that a fixable failure does not immediately strand the Work Item.
3. As a Supervisor or CI process, I want a filtered Run Event Stream for one explicit Run, so that I can observe status, Batch, Verification, and outcome changes without grepping logs.
4. As a developer inspecting a Detached Run, I want compact file-read and edit summaries in the Console Log, so that progress remains readable while full evidence stays available elsewhere.
5. As an Agent starting a Spec Task, I want the full assigned Task plus paths to relevant context, so that I can orient without receiving entire Specs, skill documents, or prior diffs in the initial prompt.
6. As a developer running independent Tasks, I want one Task's missing prerequisite to settle that Task failed without blocking other ready Tasks, so that context-efficiency changes preserve the current scheduler policy.
7. As a Supervisor or automation author, I want the Roundfix Skill and usage documentation to explain the compact surfaces and their evidence boundaries, so that I can choose the Run Event Stream, Console Log, or Run Event Journal correctly.

## Core Features

1. **P0 - Daemon-owned Verification cycle.** After initial Agent work, the Daemon runs the configured Verification. A pass settles under the existing policy and sends no successful command output to the Agent. A failure returns only Verification Feedback to the same Agent Session. See ADR-0014 and ADR-0038.
2. **P0 - One repair attempt.** Verification Feedback permits exactly one Agent repair for the Work Item, followed by one final Daemon Verification. The final verdict settles the Work Item through existing Task and Review Issue failure policies; there is no retry counter or unbounded loop. See ADR-0038.
3. **P0 - Failure-only feedback.** Verification Feedback contains the failed command, failure diagnostics, and enough location context to act. Successful command banners, progress dots, timings, and unrelated passing-package output never enter Agent context.
4. **P0 - Independent work continues.** A Task that still fails after the repair attempt settles failed, leaves dependents blocked, and does not block ready Tasks that do not depend on it. Missing task-specific credentials follow this same policy and are not reclassified as transport failures.
5. **P0 - Run Event Stream command.** `roundfix events <run-id>` replays a read-only JSONL projection of one explicit Run's Run Event Journal. `--follow` continues from replay until the Run reaches a terminal outcome; a terminal Run replays and exits.
6. **P0 - Stable Supervisor filters.** The stream's compact default contains `task-status`, `batch`, `verification`, and `outcome`. A filter option selects any subset of those categories. Internal Run Event kinds and raw Agent payloads are not part of this command's filtering vocabulary.
7. **P0 - Machine-readable stream contract.** Each JSONL record identifies the Run, journal position, category, time, and the category-specific Work Item or verdict fields. Stdout contains records only; diagnostics go to stderr; an unknown filter or missing Run ID fails argument validation.
8. **P1 - Compact Console Log.** Detached Run progress renders file reads as `read <path> (N lines)` and edits as `edit <path> (+N/-M)` without file bodies, raw ACP JSON, or unified diffs. The record points callers to the Run Event Journal and, after settlement, the commit for full bytes.
9. **P1 - Lossless journal preserved.** Agent Run Event payloads remain raw and lossless under ADR-0008. Console compaction does not change Run Event Journal replay, retention, or the opt-in Agent log policy.
10. **P1 - Spec Context Bundle.** Every Spec Task prompt embeds the complete assigned Task and a bounded manifest of paths for the PRD, TechSpec, Task Graph, required instruction documents, relevant interfaces, and files changed by prior Tasks. Larger referenced contents and prior diffs are never embedded in the initial prompt.
11. **P1 - Replacement-session orientation.** When acpx replaces or resumes an Agent Session, the next Task prompt carries the same bounded Spec Context Bundle, so recovery does not require an unstructured repository rediscovery pass.
12. **P0 - Shipped guidance.** The Roundfix Skill and user documentation ship with the behavior change. They include copy-paste replay, follow, and filter recipes; the four category meanings; JSONL/stdout/stderr and terminal behavior; the one-repair Verification Feedback lifecycle; Console Log summary examples; the lossless journal boundary; Spec Context Bundle behavior; and recovery guidance for failed Verification and missing task prerequisites.

## User Experience

A normal successful Verification appears to the caller as a verdict and Run Event, but contributes no output to Agent context. On the first failure, the Agent receives only the actionable failure diagnostics and continues in the same Agent Session; after its repair, the Daemon reports the final verdict and settles the Work Item.

Supervisors use `roundfix events <run-id> --follow` for the four-category default stream or pass a filter subset. Each stdout line is one complete JSON object suitable for incremental parsing. Detached Run callers can still open the Console Log, but repetitive reads and edits occupy one summary line each rather than carrying their contents.

## Non-Goals / Out of Scope

- Pruning, normalizing, or re-serializing raw Agent payloads in the Run Event Journal. See ADR-0008.
- A global multi-Run event stream, implicit newest-Run selection, or raw Run Event kind filters.
- More than one Verification repair, a configurable verification-attempt count, or expanding Run Budget semantics.
- Suppressing failure diagnostics that the Agent needs to repair the Work Item.
- Embedding generated summaries or selected excerpts from PRDs, TechSpecs, skills, or source files.
- Transport retry, preflight sweep resilience, force-stop integration hints, or changes to acpx's buffer implementation.
- Provisioning LLM API keys or other task-specific credentials.
- Changing Run Event Journal retention or the opt-in Agent log policy.

## Success Metrics

- Successful Verification contributes zero command-output messages to Agent context while still producing one authoritative verdict.
- A failed Verification produces at most one Verification Feedback repair and exactly one final Verification verdict.
- The default Run Event Stream emits only the four stable Supervisor categories and zero raw Agent payloads.
- A Console Log fixture containing 31 edits and 330 reads contains 31 edit summaries and 330 read summaries, with zero unified diffs and zero echoed file bodies.
- A Spec Task prompt embeds exactly one complete Task and zero complete PRDs, TechSpecs, skill documents, source files, or prior diffs.
- A Task that fails for a missing credential does not prevent an independent ready Task from starting.
- The canonical Roundfix Skill documents every new command, flag, filter, output contract, Verification transition, and evidence location; the embedded skill copy has zero sync drift.
- Documented Run Event Stream examples parse as JSONL and demonstrate replay, follow, and a filtered subset without grepping Console Log text.

## Decisions

- The Daemon owns Verification and allows one failure-only repair through the same Agent Session. See ADR-0014 and ADR-0038.
- The Run Event Stream requires an explicit Run ID and exposes stable Supervisor categories rather than raw Run Event kinds.
- Console compaction preserves the lossless Run Event Journal required by ADR-0008.
- The Spec Context Bundle embeds the Task and references all larger context by path.
- Existing independent-Task scheduling and clean failure settlement remain unchanged.
- The Roundfix Skill and usage documentation are part of the feature's completion contract, not follow-up work.

## Open Questions

None.
