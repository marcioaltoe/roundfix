---
status: completed
type: backend
---

# Task: The command help states the outcomes it accepts

Corrective, from QA finding F-01. Every other surface — the glossary, the user
guide, the shipped skill, and the command's own refusal — now names Stopped and
Unresolved as the outcomes carry-forward accepts. The command's help text still
calls it a stopped-Run act, so the first place a Supervisor looks before
choosing a mutation switch is the one place that contradicts the delivered
behavior.

## Work

- The `--carry-forward` option line names the outcomes the act accepts rather
  than describing it as a stopped-Run operation.
- Check the sibling option lines in the same help block for the same class of
  drift, and correct any that describe behavior the command no longer has.
- Read the claim from the delivered command, not from the TechSpec draft.
- Change nothing else in the command surface: no flag, no default, no exit
  code, and no other command's help.

## References

- `_prd.md` → Goal 5, User Story 4, Core Feature 2
- `_techspec.md` → API Contracts
- QA finding F-01 in `qa/qa-report-2026-08-27.md`

## Verification
- `roundfix_help="$(go run -buildvcs=false ./cmd/roundfix reconcile --help 2>&1)" || exit 1; printf '%s' "$roundfix_help" | grep -q -- "--carry-forward" && printf '%s' "$roundfix_help" | grep -qi "Unresolved" && ! printf '%s' "$roundfix_help" | grep -q "stopped Run's settled Tasks" && go test -count=1 ./internal/cli`

## Result

- The Reconcile Command's `--carry-forward` help now states that it hands a
  Stopped or Unresolved Run's settled Tasks back to the checkout. The existing
  public help test asserts that exact option description.
- The sibling `--apply`, `--discard-superseded`, and `--format` descriptions
  were checked against the delivered command behavior and remain accurate, so
  they were left unchanged.
- The focused help subtest failed before the usage-string change because the
  accepted-outcome description was absent, then passed after the change:
  `go test -count=1 ./internal/cli -run '^TestRunCommandHelp$/^reconcile$'`.
- Diff inspection confirms no flag, default, output stream, exit code, or other
  command help changed. The declared `## Verification` command was not run; the
  Daemon owns that check and Task settlement.
