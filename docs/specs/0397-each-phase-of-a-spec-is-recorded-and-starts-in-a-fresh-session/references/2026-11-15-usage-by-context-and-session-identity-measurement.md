# Measurement: where Claude and Codex usage falls on one host, and what a command run by a session can know (recorded by the author, 2026-11-15, 20:18–20:30 BRT)

Machine: Mac mini, Claude Code 2.1.325, Roundfix 0.145.1. The events come from the operator's usage script, which is outside this repository and was only read and run: `ai_usage.py collect --since 2026-11-08`, run once at 20:18 BRT. It returned 19,508 events from the transcripts and session stores of this host: 10,452 Claude Code calls, 7,170 Codex turns, 1,852 OpenCode messages and 34 status line readings. The days with Claude calls are 2026-11-08, 09 and 11 to 15.

Definitions used in every table:

- **Cost** is the script's API-equivalent cost of a call, in USD-eq: the model's price times `input + 5 × output + 1.25 × cache write (5 min) + 2 × cache write (1 h) + 0.1 × cache read`. The script's price file says the weights were fitted on 2026-11-14 against ccusage 20.0.51 and serve "only to split the real plan percentage across projects". It is a proportion, never a bill.
- **Context of a call** is `input + cache read + cache write` tokens: what the call carried.
- **Where a call ran** is what the script resolves from the call's working directory: `repo` when it is a repository checkout, `roundfix-run` when it is under `/.roundfix/worktrees/`, `other` otherwise. The script finds a Run id only when the directory name holds one.
- **This repository** is the calls the script resolves to the Roundfix project. A subagent's calls carry its parent session's working directory.

## 1. Almost no Claude usage happens inside a Run

| Where the call ran | Claude cost (USD-eq) | Share |
| --- | --- | --- |
| A repository checkout (`repo`) | 1,029.0 | 94.2% |
| Another directory (`other`) | 57.7 | 5.3% |
| A Roundfix worktree (`roundfix-run`) | 5.2 | 0.5% |
| Total | 1,092.0 | 100% |

Of the 5.2 USD-eq in a Roundfix worktree, 3.9 (0.4% of the total) carry a Run id and 1.3 do not: those calls ran in a Delivery Queue item worktree, where the pre-PR review runs outside a Run.

Main sessions hold 568.2 USD-eq (52.0%) and subagents 523.7 (48.0%).

## 2. Runs spend Codex, and the directory names the Run but rarely the Task

| Where the Codex turn ran | Tokens (millions) | Share |
| --- | --- | --- |
| A Roundfix worktree whose name holds a Run id | 907.9 | 98.1% |
| A Roundfix worktree with no Run id (Delivery Queue item worktrees) | 14.0 | 1.5% |
| A repository checkout or another directory | 3.3 | 0.4% |
| Total | 925.2 | 100% |

Of the 907.9 million tokens with a Run id, 885.6 ran in the Run Worktree, whose name holds only the Run, and 22.3 (2.5%) ran in a Task Worktree, whose name also holds the Task. Every Codex turn came through acpx.

## 3. In this repository the subagents hold three quarters of the Claude cost

| This repository, by agent | Cost (USD-eq) | Share |
| --- | --- | --- |
| `general-purpose` subagents | 346.0 | 74.0% |
| The main session | 120.1 | 25.7% |
| `Explore` subagents | 1.6 | 0.3% |
| Total | 467.8 | 100% |

The total is 42.8% of the host's Claude cost, all of it in a repository checkout. By model: Opus 301.1 (64.4%), Fable 164.3 (35.1%), Sonnet 2.4 (0.5%). The agent type says `general-purpose` for every subagent the Supervisor briefed, so nothing in a transcript says which Spec or which phase a subagent worked on.

On 2026-11-14 alone, the day one Supervisor session had the PRD, TechSpec and Tasks of two Specs written: subagents on Fable 145.2, subagents on Opus 39.1, the main session on Opus 23.8 and on Fable 16.3, subagents on Sonnet 2.4.

## 4. Most calls carry 200 thousand tokens or more

| Calls | Number | Median context | 90th percentile | Share of calls at 200k or more | Share of cost at 200k or more | Cost that is cache read | Cost that is cache read above 150k per call |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Main sessions, every project on the host | 5,823 | 230,629 | 546,916 | 59.2% | 73.6% | 75% | 38% |
| Subagents, every project on the host | 4,629 | 282,604 | 509,177 | 69.2% | 84.4% | 74% | 40% |
| Main sessions, this repository | 1,416 | 197,885 | 388,588 | 48.5% | 58.1% | 70% | 25% |
| Subagents, this repository | 3,044 | 272,232 | 500,149 | 67.9% | 83.4% | 77% | 42% |

The largest subagent call in this repository carried 966,104 tokens. The two costliest main sessions on the host belong to another repository of the fleet: 123.9 USD-eq over 1,039 calls and 25.7 hours, reaching 861,377 tokens, and 99.4 USD-eq over 594 calls and 25.0 hours, reaching 961,203 tokens. The costliest main session of this repository cost 48.2 USD-eq over 575 calls and 20.7 hours and reached 477,701 tokens. Of 57 main sessions on the host, 18 reached 200 thousand tokens.

"Cache read above 150k per call" is a counterfactual ceiling: the cost that would disappear if no call re-read more than 150 thousand cached tokens. It is not a forecast of what a rule would save.

## 5. A new session already starts large

The first call of each of the 38 main sessions opened in a repository checkout carried a median of 75,822 tokens (first quartile 51,541, third quartile 85,817). The seven sessions of this repository started at 82,386, 85,884, 87,692, 89,155, 90,244, 90,246 and 150,494 tokens.

By file size at about 3.5 bytes per token, an estimate and not a count, an author of this repository reads before anything about its Spec: `CONTEXT.md`, 182,265 bytes, about 52 thousand tokens; the guides under `docs/agents/`, 78,826 bytes, about 22 thousand; and its authoring skill, 13 to 45 thousand bytes, about 4 to 13 thousand.

Inference, not a measurement: by the script's weights a new session that writes 85 thousand tokens to the one-hour cache pays 170 thousand input-equivalent tokens once, and a call that re-reads 140 thousand cached tokens fewer saves 14 thousand, so a boundary pays for itself after about 12 calls.

## 6. A command the session runs can read the session's identity

In the Claude Code 2.1.325 session that recorded this file, the environment of a shell command held `CLAUDE_CODE_SESSION_ID`, and inside a subagent of that session the variable held the parent session's id, the same id that names the session's transcript. The variable `CLAUDE_CODE_SESSION_ATTENDED` was also set. Codex and OpenCode were not measured.

## Limits

- One host and one week. The MacBook was not read; this repository has been developed only on the Mac mini since 2026-11-11.
- This repository's Agent Selection Profiles prefer Codex for every Task Type but `frontend`, which is why its Runs spend almost no Claude usage. A repository whose Runs prefer Claude would show a larger share inside Runs.
- The cost is the script's proportion. The real limit percentage comes from the status line readings, which began on 2026-11-14 at 10:42 BRT and are too few (34) to split by phase.
- A headless Run has no status line, as the script's own notes say.
- The script and its price file may change; the numbers are those of 2026-11-15.
