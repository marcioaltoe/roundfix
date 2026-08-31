---
status: pending
type: test
---

# Task: The CLI fixtures build a qa seed the way Roundfix now renders it

Corrective, from QA finding F-001, and consequent rather than independent: the
derived Verification made these fixtures stale. `implementTaskContent` assigns a
per-seed fixture command to every seed including `qa` ones, so every assembled
CLI journey now stops at the mechanical stage with
`SC-QA-VERIFICATION-AUTHORED`, and `rtk make verify` exits 2.

## Work

- A `qa` seed with no explicit command renders the command Roundfix derives for
  its Spec, exactly as a real graph now must. Reference the derivation rather
  than transcribing its text: a fixture that copies the string stops testing the
  day the derivation legitimately changes.
- A seed that names its own command keeps it, so a fixture can still exercise
  the refusal deliberately.
- Leave the non-`qa` seeds alone. Their fixture command is what makes those
  journeys meaningful.
- The whole `internal/cli` suite passes on a fresh cache, which is where the
  failure appeared.

## References

- `_prd.md` → Goal 3, Core Feature 3
- QA finding F-001 in `qa/qa-report-2026-08-31.md`
- The consequent-fix ordering in `docs/agents/agent-instructions.md`: this lands
  after the change that made it necessary, never folded into it

## Verification
- `grep -q "DerivedQAVerification" internal/cli/implement_test.go || exit 1; go test -count=1 ./internal/cli`
