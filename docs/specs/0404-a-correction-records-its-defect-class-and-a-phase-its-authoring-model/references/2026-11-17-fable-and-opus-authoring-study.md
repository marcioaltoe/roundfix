# Study of Fable and Opus as authoring models, Specs 0370 to 0403 (2026-11-17)

Made by the operator on the Mac mini on 2026-11-17, read-only, after the maintainer asked whether Fable is effective for Ideas, PRDs and TechSpecs compared with Opus. The author of the Idea of Spec 0404 copied the study's per-Spec table and its notes here from the operator's files, which are outside the repository, and added the totals below. The table is the operator's; a `?` is the operator's own mark for a value nobody could establish.

## How the table was built

The operator's notes on the sources, as written:

- The per-Spec table merges sizes, authors, authoring cost, delivery outcome and defect classification.
- Authoring cost per authoring agent: active minutes (gaps under 5 minutes summed), wall minutes and deduplicated tokens, from the harness transcripts.
- Every subagent transcript that touches Specs 036x to 040x was scanned for its model, the Specs it references and its first prompt.
- Write and Edit tool calls to `docs/specs/<slug>/_*.md` were scanned per transcript; the scan misses the files Fable wrote through a shell.
- Slug, plan commit, TechSpec, PRD and Idea words and the count of Task files (QA included) come from Git.
- `roundfix runs list --state all` and `roundfix runs causes --since 2026-11-09` were read.
- Sources for the defect classes: the operator's memory notes of the delivery cycles, the subjects of the correction commits, the pre-PR review artifacts, and line 340 of the TechSpec of Spec 0398 at `71244939`.

Caveats, as written: Specs 0365 to 0379 were authored on another machine and have no transcripts on this one; the Fable sample is 4 Specs, all `feat`, all with an Idea, all larger.

## The per-Spec table

Columns: `tasks_excl_qa` is the number of implementation Tasks; `authoring_out_tokens_k` is thousands of output tokens and `authoring_cache_read_M` millions of cache-read tokens of the authoring agents; `authoring_active_min` is their active minutes; `corrections` is the operator's count of the corrections of the delivery; the three columns after it count causes by kind, and can add up to more than `corrections`; `review_dismissed` counts review findings dismissed as false; `queue_to_merge_h` is hours from the queue to the merge where the operator read it.

| spec | type | tasks_excl_qa | techspec_words | prd_words | idea | author_idea | author_prd_techspec | author_tasks | authoring_out_tokens_k | authoring_cache_read_M | authoring_active_min | corrections | plan_defects | impl_defects | env_or_operator | review_dismissed | queue_to_merge_h | notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 0370 | fix | 2 | 2311 | 2555 | no | - | opus (MacBook; no transcript) | opus | na | na | na | 0 | 0 | 0 | 0 | 0 | na | merged without stops |
| 0376 | fix | 2 | 3312 | 2507 | no | - | opus (MacBook) | opus | na | na | na | 1 | 0 | 1 | 1 | 0 | na | review-blocked by anchor parse (tool defect) + 2 real findings fixed |
| 0377 | fix | 1 | 1642 | 2261 | no | - | opus (MacBook) | opus | na | na | na | 1? | ? | ? | 0 | 1 | na | QA passed on 2nd round (cause not recorded); false review finding dismissed |
| 0378 | fix | 1 | 1430 | 1794 | no | - | opus (MacBook) | opus | na | na | na | 0 | 0 | 0 | 0 | 0 | na |  |
| 0379 | feat | 2 | 1290 | 1906 | no | - | opus (MacBook) | opus | na | na | na | 0 | 0 | 0 | 0 | 0 | na |  |
| 0380 | feat | 3 | 5507 | 2737 | no | - | opus-5-5 | opus-5-5 | 131 | 33.4 | 34 | 1 | 0 | 0 | 1 | 0 | na | gate race test failing on main too (pre-existing) |
| 0381 | fix | 2 | 2155 | 2271 | no | - | opus-5-5 | opus-5-5 | 139 | 41.9 | 42 | 1 | 0 | 1 | 0 | 0 | na | review: unquoted conflicted paths |
| 0382 | feat | 3 | 4025 | 2616 | no | - | opus-5-5 | opus-5-5 | 217 | 85.9 | 58 | 1 | 0 | 1 | 1 | 0 | na | review: impl violated TechSpec invariant 4; owner panic 8.5h (env) |
| 0383 | fix | 3 | 2101 | 2319 | no | - | opus-5-5 | opus-5-5 | 139 | 50.5 | 46 | 1 | 0? | 1 | 0 | 1 | ~2.8 | QA selector in task_04 did not match subtest name; fixed by renaming subtest (ambiguous plan/impl) |
| 0384 | fix | 2 | 2107 | 2502 | no | - | opus-5-5 | opus-5-5 | 158 | 43.8 | 42 | 2 | 0 | 2 | 0 | 0 | na | review Makefile quoting + later main regression (#851) |
| 0385 | fix | 4 | 2836 | 2297 | no | - | opus-5-5 | opus-5-5 | 160 | 72.1 | 45 | 1 | 0 | 1 | 0 | 0 | na | review TOCTOU |
| 0386 | feat | 2 | 4335 | 2879 | no | - | opus-5-5 | opus-5-5 | 112 | 28.5 | 37 | 3 | 1? | 2? | 1 | 0 | na | 3 review rounds on skill rule gaps (duplicate notices etc.) - rule incompleteness arguably plan; also blocked by 0384 regression |
| 0387 | fix | 1 | 1925 | 2053 | no | - | opus-5-5 | opus-5-5 | 73 | 22.3 | 19 | 0 | 0 | 0 | 0 | 0 | na |  |
| 0388 | fix | 1 | 3238 | 2156 | no | - | opus-5-5 | opus-5-5 | 84 | 15.6 | 17 | 3 | 0 | 3 | 0 | 0 | ~0.5 | 3 review rounds on test classifier/fixtures |
| 0389 | fix | 4 | 2795 | 2200 | no | - | opus-5-5 | opus-5-5 | 125 | 49.9 | 26 | 0 | 0 | 0 | 1 | 1 | na | OOM in Linux lane (env) |
| 0390 | fix | 2 | 2293 | 1700 | no | - | opus-5-5 | opus-5-5 | 117 | 34.7 | 39 | 0 | 0 | 0 | 0 | 0 | na |  |
| 0391 | fix | 2 | 4033 | 2612 | no | - | opus-5-5 | opus-5-5 | 101 | 39.5 | 30 | 0 | 0 | 0 | 0 | 0 | 1.6 |  |
| 0392 | feat | 4 | 9070 | 5818 | yes | fable-5-2 xhigh (main session) | fable-5-2 xhigh | fable-5-2 high (+4 opus prototype helpers) | 758 | 196 | 234 | 1 | 0 | 1 | 0 | 1 | 2.7 | review F1 symlink write outside repo / F2 read error; idea cost = half of shared 0392+0393 idea turns |
| 0393 | feat | 3 | 10443 | 5707 | yes | fable-5-2 xhigh (main session) | fable-5-2 xhigh | fable-5-2 (same agent) | 555 | 152 | 132 | 1 | 0 | 1 | 0 | 1 | 2.6 | QA sabotage survived (test only used relative path) |
| 0396 | fix | 1 | 3120 | 2159 | no | - | opus-5-5 | opus-5-5 | 109 | 33.8 | 28 | 0 | 0 | 0 | 0 | 0 | 1.2 |  |
| 0397 | feat | 4 | 5369 | 5613 | yes | fable-5-2 xhigh | fable-5-2 xhigh | fable-5-2 high | 560 | 128 | 110 | 2 | 0 | 2 | 0 | 0 | 4.7 | review F1 sweep error skips retention / F2 unbounded read; fixed-clock test broke after 11:00 BRT |
| 0398 | feat | 4 | 5454 | 5376 | yes | fable-5-2 xhigh | fable-5-2 xhigh | fable-5-2 high | 520 | 114 | 108 | 2 | 1 | 1 | 0 | 0 | 3.8 | review F1 = PLAN: TechSpec line 340 prescribed reading archive records from working tree; F2 control chars; perf O(n) introduced by the F1 fix |
| 0399 | fix | 3 | 3094 | 1571 | no | - | opus-5-5 | opus-5-5 | 139 | 33.6 | 32 | 3 | 1 | 2 | 2 | 1 | ~11 (overnight) | QA F-02 (QA contract in task_04) = plan; F-01/F-03 impl; Q04 env (rustup); Q06 failed by supervisor's own correction text; review UTF-8 panic |
| 0400 | fix | 3 | 2500 | 1440 | no | - | opus-5-5 | opus-5-5 | 95 | 19.9 | 22 | 1 (running) | 1 | 0 | 1 | 0 | running | task_01: 3 existing skill-coverage tests missed by the plan (explicit plan defect); QA partial env (rustup) |
| 0401 | fix | 3 | 2690 | 1627 | no | - | opus-5-5 | opus-5-5 | 96 | 20.6 | 29 | - | - | - | - | - | not delivered | amended once (ADR-0413 -> 0414) by maintainer decision not defect |
| 0402 | feat | 2 | na | na | no | - | opus-5-5 | opus-5-5 | 93 | 22.2 | 29 | - | - | - | - | - | not delivered |  |
| 0403 | feat | 4 | na | na | no | - | opus-5-5 | opus-5-5 | 90 | 21.7 | 21 | - | - | - | - | - | not delivered |  |

## Totals, computed by the author from the table

A cell with a `?` is counted apart as doubtful. Spec 0400 was still running and Specs 0401 to 0403 were not delivered, so none of them is in a total.

| group | delivered Specs | `feat` / `fix` | Tasks | corrections | plan defects | implementation defects | environment or operator | dismissed findings | TechSpec words | output tokens, thousands (mean) | active minutes (mean) |
| --- | ---: | --- | ---: | ---: | --- | ---: | ---: | ---: | --- | --- | --- |
| Fable (0392, 0393, 0397, 0398) | 4 | 4 / 0 | 15 | 6 | 1 | 5 | 0 | 2 | 5,369 to 10,443 | 520 to 758 (598) | 108 to 234 (146) |
| Opus, with transcripts (0380 to 0391, 0396, 0399) | 14 | 3 / 11 | 33 | 16 | 1, and 2 doubtful | 13 | 6 | 3 | 1,925 to 5,507 | 73 to 217 (129) | 17 to 58 (35) |
| Opus, authored on another machine (0370, 0376 to 0379) | 5 | 1 / 4 | 8 | 2 | 1 doubtful | 1 | 1 | 1 | 1,290 to 3,312 | not available | not available |

- Delivered Specs: 23. Corrections: 24. Specs with a `?` in a correction or defect cell: 3 (0377, 0383 and 0386).
- Corrections per Spec: Fable 6 / 4 = 1.50; Opus with transcripts 16 / 14 = 1.14. Corrections per Task: Fable 6 / 15 = 0.40; Opus with transcripts 16 / 33 = 0.48.
- Mean output tokens of the authoring agents: 598 thousand for the 4 Fable Specs and 121 thousand for the 18 Opus Specs with transcripts (the 14 delivered, 0400 and 0401 to 0403), a ratio of 4.9.
- Spec 0400, still running, had one plan defect: its plan missed three existing tests.
