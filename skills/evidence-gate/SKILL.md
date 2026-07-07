---
name: evidence-gate
description: Require fresh verification evidence before claiming work is complete, fixed, passing, committed, PR-ready, or ready for handoff. Use before completion claims, commits, pushes, PR creation, issue status changes, delivery summaries, or agent handoffs; do not use for brainstorming or planning that makes no delivery claim.
metadata:
  category: qa
  tags: [qa, testing, code-review, workflow]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Evidence Gate

Evidence comes before status. Do not say work is complete, fixed, passing, ready, committed, pushed, PR-ready, or handoff-ready until the current turn includes evidence for that exact claim.

## When to use

- Before saying a task is done or fixed
- Before committing or pushing
- Before opening or updating a PR
- Before marking an issue ready for review or done
- Before handing work to another agent or human
- After an agent, tool, or CI job reports success

## Gate

1. State the claim you are about to make.
2. Identify the smallest command or inspection that proves that claim.
3. Run the full command, or inspect the required artifact.
4. Read the output and exit code.
5. Compare the result to the claim.
6. Report the exact evidence, or state the blocker.

## Evidence matrix

| Claim                  | Evidence required                                  |
| ---------------------- | -------------------------------------------------- |
| Tests pass             | Test command exits 0 and reports no failures       |
| Lint clean             | Lint command exits 0 and reports no warnings       |
| Typecheck/build passes | Typecheck/build command exits 0                    |
| Bug fixed              | Original reproduction or regression check passes   |
| Requirements met       | Checklist against requested scope, not tests alone |
| Commit ready           | `git diff --cached` reviewed plus relevant checks  |
| Push ready             | Branch, commit, and remote target verified         |
| PR ready               | Verification commands, PR title/body, diff scope   |
| Handoff ready          | Current status, changed files, verification result |

## Reporting format

Use this structure in the final response or handoff:

```text
Verification:
- `<command>`: passed
- `<command>`: failed — <reason>

Not run:
- `<command>` — <reason>
```

If a check fails, do not soften the result. State the failure and the next action.

## Anti-patterns

- Saying "should pass" without running the check.
- Treating one passing check as proof of all checks.
- Trusting a tool or subagent success message without inspecting the diff or output.
- Opening a PR with missing local verification notes.
- Marking scope complete because tests pass.
