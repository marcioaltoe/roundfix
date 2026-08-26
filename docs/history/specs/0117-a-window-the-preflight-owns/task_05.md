---
status: completed
type: docs
---

# Task: Document the term and the two bounds

This Spec coins **Run Window** and adds a second time bound beside
`budget.max_run_duration`. ADR-0137 requires the Run budget to be explained
where it is configured; a reader who meets either bound must learn how the two
relate without reading source.

## Work

- `CONTEXT.md`: add **Run Window** in the shape of its neighbours — the window
  during which Runs may be created, repository-scoped, durable, governing the
  start and never the finish. Record what to avoid: `Session` is taken by
  **Agent Session**, and `curfew` and `deadline` invite the wrong reading.
- The configuration surface that prints `budget.max_run_duration` gains the
  converse pointer: that key bounds how long a Run may run, the Run Window
  bounds when one may start.
- The user guide documents the command, the refusal, and the crossing report.
- Claims are read from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.

## References

- User Story 2: The refusal names the cutoff and how to act
- Core Feature 5: The two time bounds explain each other

## Verification
- `grep -q "Run Window" CONTEXT.md && grep -q "Agent Session" CONTEXT.md && grep -qi "run window" docs/user-guide/commands.md && go test -count=1 -tags docscontract ./internal/docscontract 2>&1 | grep -q "^ok"`

## Result

Implementation:

- Added **Run Window** to `CONTEXT.md` as a durable, repository-scoped bound on
  Run creation that governs the start and never the finish. The entry reserves
  `Agent Session` and rejects `curfew` and `deadline` as misleading alternatives.
- Updated the generated config surface and the configuration reference so
  `budget.max_run_duration` explains that it bounds how long a Run may run and
  the Run Window bounds when one may start.
- Added the `window` command reference, including next-occurrence and
  repository scope, the closed-window Preflight refusal with its literal times
  and recovery commands, and the post-creation crossing report. The wording
  follows the delivered CLI output; no TechSpec-only behavior was added.
- Added the delivered implementation spellings `RunWindow` and `run_window` to
  the Run Window glossary entry so the Spec's Vocabulary Contract is satisfied
  without changing the user-facing term.

Focused-check evidence:

- Pre-change signal: `rtk rg -n "Run Window|run window|window bounds|budget\\.max_run_duration" CONTEXT.md docs/user-guide/commands.md docs/user-guide/configuration.md internal/config/config.go` exited `0` only for the existing budget references in the configuration table and validation error; it found no Run Window documentation or generated-config pointer.
- `rtk go test -count=1 ./internal/config -run 'Test(DefaultConfigYAML|RenderedConfig)'` first failed before running tests because the sandbox denied the default Go build-cache path. Rerunning with `rtk env GOCACHE=/tmp/roundfix-task05-go-cache go test -count=1 ./internal/config -run 'Test(DefaultConfigYAML|RenderedConfig)'` passed (`ok roundfix/internal/config 0.391s`).
- `rtk git -c core.fsmonitor=false diff --check` exited `0` after the final edits.
- `rtk rg -n "Run Window|Agent Session|curfew|deadline|budget\\.max_run_duration|closed at|may run past" CONTEXT.md docs/user-guide/commands.md docs/user-guide/configuration.md internal/config/config.go` found the glossary, configuration, refusal, and crossing wording.
- The Verification Feedback artifact contained no diagnostic output. The focused
  docs-contract check first reproduced `SC-VOCABULARY-UNDOCUMENTED` for the
  delivered `RunWindow` and `run_window` spellings; after documenting them,
  `rtk env GOCACHE=/tmp/roundfix-task05-feedback-go-cache go test -count=1 -tags docscontract ./internal/docscontract -run '^TestCheckActiveCorpusHasNoErrors$'` passed (`ok roundfix/internal/docscontract 0.428s`). The focused `TestCheckCorpusGolden` and `TestCheckCorpusBudget` checks also passed.
- The Daemon-owned `## Verification` command was not run in this Agent turn.
