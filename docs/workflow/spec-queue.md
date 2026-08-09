# Spec queue

The one dependency-and-risk-ordered queue of approved Specs that
`docs/agents/autonomous-work.md` requires the Supervisor to maintain. The
Supervisor works this list top to bottom. A Spec leaves the list when it
archives.

Reordered on 2026-08-09, after Specs 0082, 0083, 0084, 0088 and 0089 archived
together in pull request #143. The previous order was written on 2026-08-03 and
every Spec on it has since shipped.

## Order

| # | Spec | State | Why here |
| --- | --- | --- | --- |
| 1 | `0093-a-spec-that-validates-itself` | 7 Tasks authored | Moved to the front on 2026-08-09 after Spec 0090's gate measured the cost of checking paperwork with an Agent: 461 tool calls and two context compactions, against 161 for the average implementation Task, to find a PRD citing an ADR that says something else. Eight of that gate's sixteen rows audited artifacts rather than product. `roundfix spec check` already runs nineteen of those rules in 0.04 seconds. This Spec adds the missing rule, moves the checks to authoring, and lets the gate read product. Every Spec after it is authored and gated more cheaply, including the ones still unwritten. |
| 2 | `0090-a-gate-that-could-have-failed` | 7 Tasks authored | Every Spec below is accepted by a mechanism that, on 2026-08-09, certified work that had not happened. Task 05 of Spec 0089 settled `completed` behind a `grep` that matched an unrelated line already in the file, and its Result described `awk` inventories and a byte-level diff that were never run. The same day the authoritative gate returned two different verdicts for one tree. Nothing below can be trusted before this. |
| 3 | `0091-a-proof-that-can-refuse` | 5 Tasks authored | The same defect one subsystem over, and live in this repository today. `roundfix profiles validate` passes `claude` / `opus-9-does-not-exist` / `high` with exit 0, measured 2026-08-09; `codex` refuses the equivalent with `model_not_advertised` and `opencode` refuses an unadvertised effort. The `frontend` category routes to `claude`, so it is currently unprotected. Carries the diagnostic cleanup that today appends an unactionable session-close error after the real message. |
| 4 | `0092-a-run-that-can-hand-back-its-work` | 7 Tasks authored | The largest measured waste in the loop after gate cost, and the only item here that stops work outright. Implementing Spec 0089 re-executed Task 01 in four separate Runs and Task 02 in three, all against unchanged inputs. Then three superseded Run Branches refused two consecutive `roundfix resolve` invocations until they were removed with `git branch -D`, because `reconcile` classifies them `unintegrated` and preserves them while Branch Integrity Preflight refuses every new Run. |
| 5 | `0080-cheap-detectors-run-before-the-gate` | 8 Tasks authored | The biggest measured round-cost win, and the only Spec here that needs no authoring. Measured on Spec 0079: 92 minutes for the first Run, then 29 and 30 for rounds whose only work was moving one declared constant, while the authoritative gate runs cold in about 90 seconds. Deliberately after 1–3: its subject is gate economics, and optimising the cost of a verdict before the verdict is trustworthy is the worst possible order. |
| 6 | `0085-what-an-agent-reads-before-it-decides` | 8 Tasks authored | Cuts context cost per turn for everything after it. `docs/adr/` holds 105 records with no structural separation between the 31 accepted, the 20 carrying only a legacy body line, and the 53 with no status at all. Also owns the archive conventions that `_reviews/` still lacks — 33 pull-request folders sit in the active read path — and the findings hygiene that leaves resolved findings at `status: pending`. |
| 7 | `0081-a-journal-cheap-to-write-and-keep` | 9 Tasks authored | Cost, not correctness, but growing: the Run store was 3.1 GB when this Spec was written and is 6.8 GB on 2026-08-09. Absorbs the backlog entry for recording what each Agent Session consumed, which is not a nice-to-have — the capture step of a production loop is defined to collect diff, output, error, **cost**, status and new state. |
| 8 | `0094-separated-authorship-for-acceptance` | to author | Promoted out of 0090 rather than folded into it, on the Secondbrain's reading. Separating maker and checker by model does not help when one author wrote the Spec, the Tasks, the acceptance criteria and the interpretation of the result — which is exactly what produced Task 05. Wants the checker to receive diff, state, logs and rubric without inheriting the maker's account of how the problem was solved. |
| 9 | `0095-baseline-adoption-that-can-adopt-itself` | to author | Absorbs `greenfield-adoption-cannot-satisfy-its-own-gate`, the two HTTP Contract findings, and the `baseline-and-derived-tooling` rollup. Carries the external skill restoration left over from Spec 0082's F-002 and F-003. |
| 10 | `0096-review-that-converges` | to author | The loop's last mile: Review Issue identity across Rounds, and a loop that can act on comments about its own artifacts. Absorbs the `review-and-delivery-convergence` rollup. |

Rows 1 through 7 are authored and executable: every one has a task graph, and
`roundfix spec check` reports no findings on any of them. Numbers 8 through 10 are
intent, not reservations — each takes the next free prefix when it is minted, per
the numbering rule that scans `docs/specs/` and `docs/specs/_archived/` for the
highest.

## Why this order

The arc is trust, then unblock, then cheapen, then tidy. It was validated on
2026-08-09 against the Secondbrain, whose
`wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` states the
principle directly: the central problem is not testing the generated code, it is
testing the mechanism that decides the code is correct. The same page settles the
one ordering question worth arguing about — 0090 before 0080 — in a sentence: a
cheap gate that lies is worse than an expensive gate that declares `unknown`.

That page also lists the attacks a QA package for an agent must include. Four of
them are defects this repository hit on a single day: a command that returns exit
0 without executing assertions, a test that is absent or matches nothing, a
timeout after dispatch, and evidence that depends on context unavailable after a
restart.

Three things changed in the Specs themselves as a result of that validation.

**A verdict that can say `unknown`.** Task Verification is binary, and so is
Batch settlement. Both the verification concept and
`wiki/concepts/convergencia-observavel-em-sistemas-operacionais.md` require a
third state: a timeout or a partial execution is an unknown outcome, not a
rejection, and a claim whose sustaining content is unavailable becomes
unprovable rather than probably-still-valid. All three of the day's failures are
binary-verdict failures — a pass that should have been unknown, a batch marked
failed that should have preserved its outcomes, a gate that flipped between two
answers. This is a shared primitive of 0090 and 0092, not a detail of either.

**Three controls, not one probe.** 0090 was drafted around a pre-work probe that
refuses a Task whose Verification already passes. That is one of three controls
the literature prescribes: a positive control that must pass, a negative control
that must fail, and an observability control that must fail when the command
never reached the surface it claims. Gate health becomes a recorded metric —
how many negative controls it rejected, and when the gate's own test was last
updated.

**Authorship separation earns its own Spec.** It was going to be a bullet inside
0090. The reading is that it is a distinct mechanism: freeze the rubric before
the implementation runs, record its authorship and version, and let a reviewer
return `unproven` without being forced to approve or reject. The cheap,
mechanical half — the rubric exists and is hashed before implementation — stays
in 0090; the rest is Spec 7 above.

`convergencia-observavel-em-sistemas-operacionais` supplies 0092's frame and a
cross-project precedent: in `fluxus`, ADR-0043 already establishes that a
scheduled run which cannot start is evidence that should retry, not a silent
gap. That is the same shape as this repository's backlog entry for a Session
that never opened being recorded as a failed Batch.

## Closed by maintainer override on 2026-08-09

Four rows of Spec 0089's QA gate stayed blocked because no Roundfix Run could be
started in the gate session to read the effort receipt end to end. The Spec was
archived under `qa_override`, and the maintainer closed the rows rather than hold
delivery for them. They are recorded in the archived report as unverified by a
real Run, and the first Run that routes to the `opencode` runtime will exercise
that path incidentally.

Round 001 of pull request #143 closed the same way; the record is at
`docs/specs/_reviews/pr-143/_closure.md`.
