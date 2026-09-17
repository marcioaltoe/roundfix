# Tracked executable replay

Command:

```text
GOCACHE=/private/tmp/roundfix-0141-qa-replay-cache go run docs/specs/0141-a-commit-that-carries-the-work-and-nothing-else/qa/evidence/2026-09-17-qa/tracked_executable_replay.go .
```

Result: exit 0.

```text
kept: [.agents/skills/systematic-debugging/find-polluter.sh .githooks/commit-msg .githooks/pre-commit]
dropped: []
settlement commit paths:
.agents/skills/systematic-debugging/find-polluter.sh
.githooks/commit-msg
.githooks/pre-commit
```

The first attempt inherited `commit.gpgsign=true` and could not reach the
product assertion because the sandbox denied access to the user's GPG state.
The disposable fixture now sets `commit.gpgsign=false`; the successful replay
above is the measured result.
