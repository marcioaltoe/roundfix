# Dogfood findings — round 2 (specs 0003–0005 execution)

Running log for the second supervised dogfood cycle: executing
`0003-dogfood-polish`, the 0002 re-QA, `0004-watch-merge-readiness`, and
`0005-tui-cockpit` through `roundfix implement` (acpx layer, codex
gpt-5.5 xhigh). Round 1 lives frozen in `dogfood-findings.md` (27 findings).
Sizes: `cosmetic` / `small` / `spec`.

1. **acpx 10 MiB message buffer killed a finished turn.** (0003 Run
   `run_20260705T131519Z_d158ba8f114de398`, task_09.) The Agent completed and
   verified all work; at end of turn acpx emitted `{"error":{"code":-32603,
   "message":"Message buffer exceeded 10485760 bytes","data":{"acpxCode":
   "RUNTIME"}}}` and exited 1 — a single adapter message (large skill-file
   tooling output) blew acpx's per-message buffer. The runner mapped exit 1 →
   Batch failed; the Daemon settled the Task failed without verifying, and
   the Run ended Unresolved with done work preserved uncommitted. Recovery
   was manual daemon-role settlement (fresh verification + status + commit
   b756331). Follow-ups: report upstream (openclaw/acpx) and check for a
   buffer-size config; roundfix-side, see finding 2. Size: small (upstream +
   observe).

2. **Exit-code-only Batch classification can fail completed work.** The
   stream had already delivered the Agent's full completion report when acpx
   died; the runner classifies solely on process exit. Candidate: when the
   NDJSON stream carried a parsed `session/prompt` result (terminal stop
   reason) before a nonzero exit, classify the exit as teardown noise —
   journal it loudly, let the Daemon's own verbatim verification be the
   gate (it is the real arbiter anyway; ADR-0014). That keeps honesty (the
   commit still only happens on passing verification) while not discarding
   finished work over a transport hiccup. Size: small.

3. **Failed-but-done Tasks need a cheaper recovery than agent re-run.** With
   preserved work in the tree, today's options are: manual daemon-role
   settlement (hand verification + status + commit) or clean the tree and
   re-run the whole Task through the Agent. A `roundfix settle` (or
   resume-with-verify-first) that re-runs a failed Task's Verification over
   the preserved worktree and, on pass, settles + commits without invoking
   the Agent would turn this recovery into one command. Pairs with finding 2
   (which prevents the false failure) and round-1 finding 24 (graceful
   stop). Size: small/medium.

   **Recurred at 0005 task_07** (same `-32603` buffer error, again on a
   docs task touching large skill files — 2 for 2 on docs tasks since the
   cutover; manual settlement again). Findings 1–3 graduate to must-fix in
   the next spec cycle: the parsed-result-over-exit-code classification plus
   an upstream buffer report/mitigation.

4. **RESOLVED 2026-07-05** — 0004's docs task flagged that `merge-ready`
   (ADR-0019 vocabulary) had no `CONTEXT.md` entry; the **Merge-Ready**
   glossary term was added at the Run boundary in the same commit as this
   note.
