# Model recommendation ranking

Status: recommendation input, not routing policy  
Source snapshot: 2026-07-16

## Purpose and hard limit

This document curates the initial five Agent Selection recommendations that the Roundfix CLI will display for `general`, Task Type, `qa`, and `review` categories. It does not authorize automatic routing. Roundfix uses only configured Agent Selection Profiles; a ranking can help a user fill a profile but can never select or modify one.

The supplied benchmarks measure general coding-agent work. They do not publish statistically valid frontend-, backend-, QA-, or review-specific slices. Category ordering below combines those general results with an explicit qualitative task-fit judgment. Every CLI row must therefore state `category_specific: false` until a category-specific evaluation exists.

## Sources

Local source artifacts retained with this Spec:

- [Artificial Analysis coding-agent leaderboard](artificial-analysis-coding-agent-benchmarks-leaderboard.pdf) — Coding Agent Index v1.1, result, cost, and duration snapshots.
- [Artificial Analysis leaderboard image](artificial-analysis-coding-agent.jpeg) — captured source view supplied with the request.
- [DeepSWE models](deep-swe-models.pdf) — DeepSWE v1.1, 113-task result, cost, token, and step snapshot, updated 2026-07-13.

External source pages:

- Artificial Analysis Coding Agents: https://artificialanalysis.ai/agents/coding-agents
- DeepSWE leaderboard: https://deepswe.datacurve.ai/
- OpenAI GPT-5.6 Sol: https://developers.openai.com/api/docs/models/gpt-5.6-sol
- OpenAI GPT-5.6 Terra: https://developers.openai.com/api/docs/models/gpt-5.6-terra
- OpenAI GPT-5.6 Luna: https://developers.openai.com/api/docs/models/gpt-5.6-luna
- OpenClaw autoreview: https://github.com/openclaw/agent-skills/blob/main/skills/autoreview/SKILL.md

## Method

1. Use official model identifiers accepted by the relevant ACP Runtime.
2. Limit every category to five unique Agent Models; multiple effort variants of the same model do not occupy multiple positions.
3. Use the DeepSWE row for the exact recommended effort when result and average cost are shown.
4. Consider the broader Artificial Analysis Coding Agent Index as a cross-check because it combines DeepSWE, Terminal-Bench v2, and SWE-Atlas-QnA.
5. Prefer a quality/cost balance for routine work; reserve maximum-effort recommendations for models where the benchmark gain justifies the added cost or where the model serves as a quality-preserving fallback.
6. Apply task-fit judgment only to order the five candidates. Never present that judgment as category-specific benchmark evidence.

The benchmark values are snapshots, not promises. Provider price, availability, adapter behavior, and benchmark methodology may change. A ranking update requires new source artifacts, review, and a Roundfix release; it never rewrites existing configuration.

## Cross-check snapshot

The Artificial Analysis Coding Agent Index v1.1 snapshot reports these broad scores for relevant models at the listed leaderboard setting:

| Agent / model setting | Index |
| --- | ---: |
| Codex / GPT-5.6 Sol max | 80 |
| Codex / GPT-5.6 Terra max | 77 |
| Claude Code / Fable 5 max with fallback | 77 |
| Codex / GPT-5.5 xhigh | 76 |
| Codex / GPT-5.6 Luna max | 75 |
| Claude Code / Opus 4.8 max | 73 |

The “with fallback” Artificial Analysis Fable row is not evidence for Roundfix's fallback semantics. Roundfix separately configures and proves its own Fallback Chain.

## Initial recommendations

### General

| Rank | Agent Selection | DeepSWE result | Average cost/task | Recommendation rationale |
| ---: | --- | ---: | ---: | --- |
| 1 | `codex / gpt-5.6-sol / high` | 69% | $3.47 | Built-in quality/cost default; Sol leads the broader composite at max without requiring max for routine work. |
| 2 | `codex / gpt-5.6-terra / max` | 70% | $4.95 | Strong quality-preserving fallback with lower average cost than the largest Claude settings. |
| 3 | `claude / claude-fable-5 / high` | 69% | $9.18 | Cross-runtime high-quality option when provider diversity matters. |
| 4 | `codex / gpt-5.6-luna / max` | 67% | $3.03 | Cost-conscious option that remains competitive at maximum effort. |
| 5 | `codex / gpt-5.5 / xhigh` | 67% | $7.23 | Retained proven-generation option, but less attractive than GPT-5.6 on this snapshot. |

### Backend

| Rank | Agent Selection | DeepSWE result | Average cost/task | Recommendation rationale |
| ---: | --- | ---: | ---: | --- |
| 1 | `codex / gpt-5.6-sol / high` | 69% | $3.47 | Best default for complex repository changes, debugging, and multi-file implementation. |
| 2 | `codex / gpt-5.6-terra / max` | 70% | $4.95 | Quality-preserving alternative with a strong long-horizon result. |
| 3 | `claude / claude-fable-5 / high` | 69% | $9.18 | Cross-runtime option for architecture-heavy or broad refactoring work. |
| 4 | `codex / gpt-5.5 / xhigh` | 67% | $7.23 | Stable legacy choice when a repository has already validated it. |
| 5 | `codex / gpt-5.6-luna / max` | 67% | $3.03 | Lower-cost option for bounded backend slices. |

### Frontend

| Rank | Agent Selection | DeepSWE result | Average cost/task | Recommendation rationale |
| ---: | --- | ---: | ---: | --- |
| 1 | `claude / claude-fable-5 / medium` | 65% | $6.09 | Preferred balanced frontend profile; task-fit is qualitative because the benchmark has no UI-specific slice. |
| 2 | `codex / gpt-5.6-sol / high` | 69% | $3.47 | Strong cross-runtime fallback for implementation and integration correctness. |
| 3 | `claude / claude-opus-4-8 / high` | 52% | $4.28 | Design-heavy alternative retained for repositories that have validated Opus; general benchmark result is weaker. |
| 4 | `codex / gpt-5.6-terra / max` | 70% | $4.95 | Quality-focused Codex alternative when visual specialization is less important than repository completion. |
| 5 | `codex / gpt-5.6-luna / max` | 67% | $3.03 | Cost-conscious option for small, well-specified UI slices. |

### QA

| Rank | Agent Selection | DeepSWE result | Average cost/task | Recommendation rationale |
| ---: | --- | ---: | ---: | --- |
| 1 | `codex / gpt-5.6-sol / high` | 69% | $3.47 | Built-in QA default for broad behavior checks and evidence synthesis. |
| 2 | `codex / gpt-5.6-terra / max` | 70% | $4.95 | Quality-preserving fallback for long validation flows. |
| 3 | `codex / gpt-5.6-luna / max` | 67% | $3.03 | Cost-conscious option for deterministic, well-bounded QA matrices. |
| 4 | `claude / claude-fable-5 / high` | 69% | $9.18 | Cross-runtime option for exploratory behavior analysis. |
| 5 | `codex / gpt-5.5 / xhigh` | 67% | $7.23 | Legacy option for established QA setups already proven on GPT-5.5. |

### Review

| Rank | Agent Selection | DeepSWE result | Average cost/task | Recommendation rationale |
| ---: | --- | ---: | ---: | --- |
| 1 | `codex / gpt-5.6-sol / high` | 69% | $3.47 | Matches the OpenClaw autoreview default and balances review depth with cost. |
| 2 | `codex / gpt-5.6-terra / max` | 70% | $4.95 | Strong alternative when Sol cannot start or a quality-preserving fallback is required. |
| 3 | `claude / claude-fable-5 / high` | 69% | $9.18 | Cross-runtime review option with a competitive general coding result. |
| 4 | `codex / gpt-5.5 / xhigh` | 67% | $7.23 | Established predecessor for repositories that prefer a validated older generation. |
| 5 | `claude / claude-opus-4-8 / high` | 52% | $4.28 | Design- and architecture-oriented alternative; lower general benchmark result keeps it fifth. |

### Optional Task Types

`data`, `infra`, `docs`, `test`, and `chore` use the General recommendation list in the initial release. The CLI still labels the requested category, shows exactly five entries, and reports that the recommendation source is `general`. A category-specific list should be added only with documented evidence and a release-reviewed ranking change.

## Interpretation rules for the CLI

- Display benchmark result and cost as source snapshot data, not current price or guaranteed quality.
- Display the source date and `category_specific: false` with every result set.
- Never preselect rank 1 in a way that bypasses explicit user confirmation during configuration.
- Never change an existing profile when ranking data changes.
- Never hide the configured fallback behind recommendation output; the effective profile remains the primary display.
- If a recommended tuple fails disposable proof, report it as unavailable on the machine without reordering or substituting the ranking.
