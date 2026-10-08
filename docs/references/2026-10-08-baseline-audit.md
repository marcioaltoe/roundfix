# Baseline audit — 2026-10-08

Read-only audit of `/Users/marcio/dev/roundfix` at `main` `c73a92e0` (after v0.59.0). Nothing in the repository was changed. No adopter or Secondbrain content was sent to Jev.

Scope: Baseline catalog `internal/baseline/assets/` (16 modules, 177 clauses, about 57 KB of guidance text), the guides rendered into `docs/agents/*.md` (about 67 KB, read at the start of every Run session), the 41 required skills, ADR-0228..0253, `docs/user-guide/**`, the Run Database `~/.roundfix/roundfix.db` (Specs 0225–0248: 41 Runs, 104 Task sessions, 34 QA sessions, 28.8 h of Run wall time), `queue-interventions.md` (entries 145–240), the 230 Archive Records, and `BRIEFING-authoring.md`.

## Headline evidence

| Measure (Specs 0225–0248 unless noted) | Value | Source |
|---|---|---|
| Task-agent runs of `make verify-incremental` (full `go test`) inside Runs | 74 runs in 20 of 23 Specs, mean 231 s, **284 agent-minutes total** | `run_events` `agent.tool_started`/`tool_updated` |
| Daemon `make verify-changed` per Task settlement | 125 runs, mean 350 s; it grew from about 150–200 s (0200–0205) to 400–450 s (0244–0248) | `run_events` `daemon.verification` |
| QA verdicts | 11 fail, 14 partial, 9 pass (about 1.5 QA attempts per Spec; mean 16–19 min each) | `daemon.qa` |
| Archive dispositions 0180–0248 | 26 `qa-override`, 2 `partial`, the rest `pass`. Since ADR-0240 (0236+), 4 of 11 Specs still needed an override: 0241, 0242, 0243, 0245 | `docs/history/specs/*.md` |
| QA reruns caused by Surface Transcript drift where the code was correct | 0219, 0241, 0242 (F-05), 0243, 0247 | interventions 132, 213, 220, 227, 240 |
| QA failures caused by stale wording in another guide | 0235 F2, 0242 F-03, 0244 (twice: settle.md, then reopen.md), and the CI failure on 0240 (profiles.md) | interventions 193, 218, 231, 232, 208 |
| Runs or CI lost to non-hermetic or host-specific tests | 0235 (test read the host's `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`), 0241 (macOS socket path over 104 bytes, so the test SKIPs), EPERM flakes on 0212/0215/0228 | interventions 188, 211, 99, 106, 161 |
| CI failures on suiteguard installation after a green Run | 0231 (`internal/delivery`), 0245 (`internal/config`), plus the ADR-0252 incident | interventions 175, 238; ADR-0252 |
| Operator fixes for `SC-ADR-RELATED` on Specs authored in parallel | 0222 (parked revalidation-failed), 0228 (flaky-check), 0234, 0236, 0237 | interventions 147, 164, 184, 197, 199 |
| Required skills vs installed `.agents/skills` | 41 = 41, no drift; the retired `council`/`the-fool` are gone from dispatch | modules `requiredSkills`; `.agents/skills/` |

## Findings

Kinds: **stale** = contradicts current behavior or an ADR; **contradiction** = two live rules disagree; **gap** = missing rule that evidence shows would have prevented rework; **cost** = rule that spends Run time without benefit.

| id | file:line | kind | evidence | proposed change |
|---|---|---|---|---|
| B01 | `internal/baseline/assets/modules/spec-workflow.json:105` (`clause.spec.verification-two-tiers`) → `docs/agents/spec-routing.md:23`; also `core.json:422,427` → `docs/agents/agent-instructions.md:31,33` | contradiction + cost | The clause says "For each Task, run the selected incremental Verification … before handoff". `implement-task` §7 says that in a Daemon-assigned turn the agent runs focused checks and the Daemon runs the declared Verification and the Run gate. Agents follow the Baseline: 74 full-suite runs, 284 agent-minutes, competing for CPU with the Daemon's serialized `verify-changed`. The Daemon's `verify-changed` failed only 6 of 125 times, so the agent-side run caught almost nothing the Daemon would not have caught. | Scope both tier clauses: outside a Run, run the incremental tier before handoff. Inside a Daemon-assigned turn, run focused tests of the changed packages only, because the Daemon runs the Task Verification and the Run Verification (`defaults.verification`) at settlement. Each reworded clause gets a new id that `replaces` the old one (ADR-0222), plus a retention disposition. |
| B02 | `docs/agents/setup-context.json:54-58,255-266` (`verification.gate` `rtk make verify`, `verification.incremental` `rtk make verify-incremental`); `internal/baseline/assets/decisions.json:116,135` | stale | `rtk` is a local output filter on this machine. CI runs `make verify`, `make verify-docs` and `make verify-contracts` (`.github/workflows/ci-verify.yml:75-86`), and the clause says "CI must run the selected repository Verification". `verify-incremental` is `fmt-check vet test …`, the full suite, so it is not incremental. The user guide (`context-driven-development.md:262,630`) repeats the `rtk` values. | Record `verification.gate: make verify` and `verification.incremental: make verify-changed`, the same selective gate as the Daemon (a repository decision through `baseline update`, not a module edit). Make `--adopt-suggested` stop proposing an `rtk`-prefixed command, and update the user guide example. |
| B03 | `docs/agents/autonomous-work.md:5-7`; `decisions.json:261,280` (`runtime.backend` / `runtime.design` defaults) | stale + contradiction | The guide says backend work uses `codex gpt-5.6-sol` and design work uses `claude opus 5 xhigh`. `.roundfixrc.yml` profiles use `gpt-6.1-sol/high`, `opus/high` and `gpt-5.6-luna/max`, plus the light tier (ADR-0238). ADR-0037 and ADR-0049 make Agent Selection Profiles the only source of model choice. Every adopter gets the stale default. | Retire both decisions (with a retention disposition) and render "Runs use the Agent Selection Profiles of Project Config (`roundfix profiles`); `complexity: low` Tasks without a Governed Path run on the light tier". If retiring is refused, at least stop rendering a model name. |
| B04 | `core.json:408` (`request-pull-request-review`) and `core.json:618` (`request-review-explicitly`) → `agent-instructions.md:27,77`; `autonomous-work.json:137` (`loop-05`) → `autonomous-work.md:25` | cost (duplicate) | The review-policy rules appear three times, about 2.3 KB. The two core clauses overlap by most of their text: the provider set, the meaning of `none` and blocking on stale review. | Merge into one core clause and keep only "Clean is a claim" in loop-05 (with `replaces`). |
| B05 | `spec-workflow.json:139,144,149,183` → `issue-tracker.md:5,7,9` and `spec-routing.md:33` | cost (duplicate) | "Status lives only in the Task file; dependencies only in `_tasks.md`" is stated four times. | Keep `tracker-artifacts`; retire `status-only-in-task` and the tail of `project-constraints-05`. |
| B06 | `core.json:560,570` → `agent-instructions.md:55,59`; `secondbrain.json` → `secondbrain.md:46`; `secondbrain.md:7,9` | cost (duplicate) | "Use local tools for local code; never external research" appears three times. Secondbrain lines 7 and 9 both restate the index-first query order. | One core clause; drop the secondbrain repetitions. |
| B07 | `AGENTS.md`/`CLAUDE.md` root blocks ("Domain and documentation rules are mandatory: … docs-layout.md"; "Optional cross-project knowledge follows … secondbrain.md") | cost | Run sessions read all of `docs/agents` at start (504 tool calls in 23 Specs, about 17k tokens per session). `docs-layout.md` (16.6 KB) is mostly the Backlog, Findings and Inbox record templates. No Run Task writes those records, and the Run sandbox cannot reach the Secondbrain or Exa that `secondbrain.md` (7 KB) mandates. | Root block: read `docs-layout.md` before creating, moving or retiring a document under `docs/`, and `secondbrain.md` only when authoring or researching outside a Run. Optional later step: split the record templates into their own guide. |
| B08 | `spec-workflow.json:188` (`project-constraints-06-outside-evidence`) → `spec-routing.md:35` | gap | Four of the last five overrides (0241 R8b, 0242 row 12, 0243 R08, 0245 Q08c) were outside-evidence rows pointing at "the original authoring-ablation logs" or at a lookup the authorization forbids. Three more rows needed github.com, a fake judge transport the CLI does not expose, or a TLS certificate macOS rejects (0241 R4, 0242 row 07, 0243 R05). Each one meant an override, a cherry-pick and a retry, and sometimes a review finding against the override. | New clause (replaces the old one). The outside source must be committed (`docs/references/` in the planning PR) or published, and reachable by the QA agent without network before the Run starts. Never cite an artifact that exists only on the operator's machine or a lookup the Spec's authorization forbids. A row that needs the network or an open PR is declared under Unreachable Acceptance when authored, and is not left to become an environment block. |
| B09 | `spec-workflow.json` (no Surface Transcript clause); `.agents/skills/write-tasks/SKILL.md:260-267`; `write-techspec/references/concrete-contracts.md:29` | gap | Five QA reruns in 0219–0247 were transcript lines the implementing test never asserted, such as a final Plan Digest line, `exit status 1`, indentation or `Usage` lines. The code was correct each time. A rerun costs about 22–25 min (repository-Verification precondition plus QA) and an operator amendment. | The implementing Task's test asserts every line of the transcript's stdout and stderr plus its exit code, copied whole, with only the declared controlled variations. Its Verification runs that test by name. Add it to `write-tasks` (owned skill, version record) and as one spec-workflow clause. |
| B10 | `autonomous-work.json:132` (`loop-04-verify-the-class`) → `autonomous-work.md:23` | gap | Stale wording in a second guide failed QA on 0235, 0242 and 0244 (twice in a row, because the first fix did not grep everywhere) and failed CI on 0240. | Extend loop-04: a Task that changes a message, field, flag, refusal or rule greps the old wording across `docs/`, `skills/` and `.agents/skills/` and declares each file it must change. QA re-checks with the same sweep. |
| B11 | `go.json:173` (`test-observable-behavior`) → `docs/agents/go.md:22` | gap | 0235 was parked because a test read the host's real key. 0241 needed a corrective Task because a socket path over 104 bytes made a required test SKIP on macOS. There were EPERM flakes in daemon verification on 0212, 0215 and 0228. | New Go clause: tests set or clear every environment variable the code reads (`t.Setenv`), never read the host's keys, HOME or Roundfix Home, keep unix-socket paths short (`os.MkdirTemp("", …)`, not a deep `t.TempDir`), and never let a Verification depend on a test that can SKIP on the host. |
| B12 | `docs/agents/specific-repository.md` (repository-owned; no suiteguard rule) | gap | Suiteguard installation was missed three times: 0231 and 0245 in CI, and the ADR-0252 incident. Since ADR-0252 the contract only runs when it is relevant. | Add a rule: every test package that spawns a process installs `suiteguard.Main` in `TestMain`, and `TestEverySpawningPackageInstallsTheSuiteGuard` enforces it. |
| B13 | `autonomous-work.json:117` (`loop-01-qa-once`) → `autonomous-work.md:17` | gap | 0222, 0228, 0234, 0236 and 0237 each needed an `SC-ADR-RELATED` fix by hand after a sibling Spec merged a new ADR. 0222 parked `revalidation-failed` and 0228 `flaky-check`. | Add: before a Spec authored in parallel enters a Run or the queue, rebase it on the default branch and run `roundfix spec check <slug> --strict --run-verification`. |
| B14 | `autonomous-work.json:117` (`loop-01`) | gap (ADR-0249) | Reopen proves a Late Dependency only if the corrective Task is in the QA Task's dependency closure. 0245 needed "QA needs task_05" by hand (intervention 233). | Add: a corrective Task is added to the QA Task's `needs`, then `roundfix reopen --spec <slug>`. |
| B15 | `spec-workflow.json:125` (`verification-fails-before-the-change`) → `spec-routing.md:21` | stale | The clause asks for `spec check --run-verification`. The strict mode is what reports `SC-GLOSSARY-*` and `SC-ADR-RELATED` as blocking. The briefing and every Finish step use `--strict`. | Add `--strict` to the clause. |
| B16 | `context-workflow.json:278` (`docs-one-job-per-directory`) → `docs-layout.md:119`; `findings-09` (`:308`) → `docs-layout.md:173` | stale (needs maintainer decision) | ADR-0248 removes retired Review Artifacts and handoffs and reduces retired Findings and Backlog Entries. `docs/history/` now holds only `adr/ backlog/ findings/ specs/`. The clause still tells new retirements to move whole into `docs/history/reviews/` and `docs/history/handoffs/`, which the next sanitize deletes. That contradicts the 2026-10-06 rule that history holds only Secondbrain-relevant material. | Decide whether new retirements are written in reduced form directly (Reviews and handoffs removed; Findings and Backlog reduced at retirement), or whether a sanitize step is run after each retirement. Then reword the clause. |
| B17 | `context-workflow.json:234` (`read-domain-contract`) → `domain.md:7` | cost (low) | "Read both [glossary and accepted ADRs] before … changing behavior" points at an 80 KB glossary and about 250 ADRs. Agents mostly grep (106 grep reads against 146 whole-file reads of CONTEXT.md; 206 ADR reads). | Change to: read the terms the Spec's Glossary Declaration and requirements name, and the ADRs its Active ADR obligations row lists. Search the rest. |
| B18 | `spec-workflow.json:202` (`keep-artifacts-in-spec-folder`) → `docs-layout.md:198` | cost (low) | 2.1 KB, mostly the QA Archive Override procedure, which runs only through `roundfix archive --qa-override` and the `archive-spec` skill. | Keep the prohibition on hand-editing and the meaning of the Archive Record. Move the override procedure out to the skill and the user guide. |
| B19 | `cli-surface.json` → `cli.md:9` and `specific-repository.md:30-32` | cost (duplicate) | The skill-sync rule is stated twice, once in the Baseline and once in repository text. | Delete the repository HARD RULE. The Baseline clause covers it. |
| B20 | `core.json:642` (`ask-user-answerable-decisions`, 1.5 KB) → `agent-instructions.md:79` | cost (low) | Read in every non-interactive Run session, where no user-interaction tool exists. | Shorten. Keep the full protocol in a supporting reference, or render it only for interactive sessions. Defer unless Spec A has room. |

Verified **with no change needed**: the Archive Record (ADR-0247), the glossary duty (ADR-0244 via `clause.domain.glossary-currency`), the QA partial policy and network-denied rows (ADR-0240 in loop-01 and project-constraints-06), the retired skills (ADR-0246), and Lost Rollout and the light tier (runtime-only; the light tier is documented in `write-tasks`). Skill dispatch matches the installed skills exactly.

Out of Baseline scope, but measured:
- `verify-changed` duration doubled (about 200 s to 420 s) between 0200 and 0248.
- Review still flags authorized `qa_override` archives (0215, 0218, 0241).
- 0241 and 0242 still needed a manual PR after a review fix.

These belong in Backlog Entries.

## (c) BRIEFING-authoring.md rules

| Rule | Disposition |
|---|---|
| 1, 4, 5, 7, 10, 17, "consult Secondbrain + Exa" | Already in `write-tasks` or the Baseline (spec-workflow, go, secondbrain). Delete from the briefing. |
| 6 (do not raise the Roundfix skill version until Spec 0195) | Stale: Spec 0195 has landed and ADR-0233 makes the record command choose the version. Delete. |
| 16 (backticks; "spec check still calls it honest") | Stale: ADR-0226 makes such a command never honest. Keep one line in `write-tasks` only. |
| 15 (glossary) | Now `clause.domain.glossary-currency` plus `SC-GLOSSARY-*`. Reduce to a pointer. |
| 8 (outside evidence: never another repository's manifest) | Move to the Baseline as B08, extended to operator-only artifacts and network. |
| 13 (cross-Spec Verification dependency → Prerequisites + Build Order) | Move to spec-workflow next to B13. |
| 9 (characterize before changing; declare each break) | Move to `loop-04` (B10 Task), since the maintainer's "specs evolve, never regress" rule is not yet in any guide. |
| 3 (`### QA settlement` byte-identical), 11 (clause removal needs a retention disposition and a `baseline_update_test` fixture), 12 (canonical `.agents/skills` plus `make skills-sync`), 14 (derived Baseline files only via `make baseline-digests`), 18 (Sanctioned regeneration with explicit `outputs:` for `-record-module-versions`) | Specific to Roundfix. Move to `docs/agents/specific-repository.md` (repository-owned text, no module bump). 18 explains 0245 F2. |
| Authorization section (`paths:` measured with a `GovernedPath` probe, `paths: []`) | Move the measuring technique to `specific-repository.md`. The record shape is already Baseline (`project-constraints-02`). |
| Maintainer decisions 2026-09-30 and the upstream HEAD | Program state, not a rule. Keep it in operator memory. |

## Proposed Specs (at most 4 implementation Tasks plus QA each)

### Spec A — One gate per Run and a smaller Baseline (performance)

Findings B01, B02, B03, B04, B05, B06, B07, B15, B19.

1. **Verification tiers inside a Run.** Change `core` and `spec-workflow` with new clause ids that `replace` the old ones (B01, B15). Align `implement-task` §7 wording if needed (owned skill, version record). Bump **core** and **spec-workflow**.
2. **Decision values and runtime header.** Record `verification.gate`/`verification.incremental` without `rtk` (B02; setup-context through `baseline update`), and fix the `--adopt-suggested` suggestion and `context-driven-development.md`. Retire `runtime.backend`/`runtime.design` or stop rendering a model (B03): `decisions.json`, the autonomous-work guide template, and a retention disposition. Bump **autonomous-work**; there are decision-version effects on every adopter.
3. **De-duplicate.** Merge the review-policy clauses (B04), the tracker clauses (B05) and the local-research clauses (B06) with `replaces`, and delete the repository HARD RULE duplicate (B19). Bump **core**, **spec-workflow**, **secondbrain**.
4. **Root reading set.** Change the root-block wording for `docs-layout.md` and `secondbrain.md` (B07). Bump **context-workflow** and **secondbrain** (root blocks). Optionally add B17, which also bumps context-workflow.
5. **QA.**

Tasks 1, 3 and 4 share `internal/baseline/module-versions.json`, the digests and `docs/agents/*`, so they run in series (Wave collision). Each declares the two sanctioned regeneration commands with `outputs:`.

### Spec B — Authoring rules that stop QA reruns (results)

Findings B08, B09, B10, B11, B12, B13, B14, plus the briefing migration (c).

1. **Outside evidence and transcripts in spec-workflow.** B08 (new clause replacing `project-constraints-06`), B09 (transcript clause), and briefing rule 13. Bump **spec-workflow**.
2. **Authoring skills.** The `write-tasks` transcript assertion rule and outside-evidence guidance, with `write-techspec` `concrete-contracts.md` aligned (owned skills, both version fields, the `owned-skill-versions.json` record). No module bump.
3. **Loop and Go rules.** `loop-01` gets the parallel-Spec re-check (B13) and the Late Dependency `needs` rule (B14). `loop-04` gets the wording sweep and characterize-before-change (B10, rule 9). A new Go hermetic-test clause (B11). Bump **autonomous-work** and **go**.
4. **Repository rules.** `specific-repository.md`: suiteguard (B12), plus briefing rules 3, 11, 12, 14 and 18 and the GovernedPath probe. Glossary: add any new term, such as "Run reading set" if Spec A coins it. No module bump. After merge, the operator trims `BRIEFING-authoring.md` outside the repository.
5. **QA.**

Run Spec B after Spec A, or rebase it: both bump `spec-workflow` and `autonomous-work`, so the second one takes "the next version". B16 (the history retirement form) needs a maintainer decision first. Once decided, it fits Spec B Task 1 or a third Spec.

## Run-time estimate

Baseline: Specs 0225–0248 averaged about 75 min of Run wall time per Spec (28.8 h / 23), plus about 1.5 QA attempts and several operator interventions per Spec.

| Change | Expected effect |
|---|---|
| B01 + B02 (no agent-side full suite inside Runs) | About 12 agent-minutes per Spec removed (284 min / 23). Wall-clock saving is estimated at **5–10 min per Spec**: Task sessions finish sooner, and the serialized Daemon `verify-changed` (350 s mean) stops competing with up to three concurrent `go test ./...` runs. |
| B09 (complete transcript assertion) | Removes about 0.4 QA reruns per Spec (5 in 0219–0247, 4 in the last 9 Specs). About 22–25 min each, so roughly **9–11 min per Spec** plus the operator's amendment. |
| B08 (reachable outside evidence) | Most of the remaining overrides disappear (4 of the last 11 Specs). Each saves about 10–20 operator-minutes (carry-forward, cherry-pick, `archive --qa-override`, retry) and one probable review false positive. |
| B10 (wording sweep) | About 4 QA or CI failures in 23 Specs, roughly **4 min per Spec**. |
| B11 + B12 + B13 + B14 | Prevents the park, fix-PR and CI-rerun classes behind about 10 interventions in entries 145–240. These are operator-time savings rather than Run time. |
| B04–B07, B17–B20 (smaller guides) | About 15–25 KB less guide text per session, which is a few thousand tokens and seconds per session across about 140 sessions per 23 Specs. This is a small time gain; the main value is less contradictory text. |

Combined estimate: **about 20–25 min less wall time per Spec (about 25–30 %)**, and about a third fewer operator interventions in the classes seen from 2026-10-04 to 2026-10-07. The estimate assumes the defect mix of 0225–0248 continues. The `verify-changed` growth is a separate cost that this audit does not address.
