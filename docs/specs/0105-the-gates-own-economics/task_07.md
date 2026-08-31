---
status: pending
type: backend
---

# Task: The derived Verification passes the checker that demands it

Corrective, from this Spec's own first Run. `SC-QA-VERIFICATION-AUTHORED`
requires a `qa` Task to carry exactly the command Roundfix derives, and
`SC-VERIFY-NON-HERMETIC` refuses that command: it reads the awk regex
`/^.*\/qa-report-/` as an external path. The Spec cannot satisfy its own two
detectors at once, so no `qa` Task can be authored at all.

## Work

- Express the derived command without regex literals the hermeticity detector
  reads as paths. Stripping a directory prefix is the offending construct;
  splitting the path on its separator answers the same question without a
  regex that looks like one.
- Keep the rigour that made the derived command worth deriving. It must still
  select the newest report by parsed date and sequence rather than by
  lexicographic sort, because selecting an older report is the failure this Spec
  removes. A simpler command that reintroduces that failure is not a fix.
- Re-render `task_06.md` with the corrected command. It carries the refused one
  today, which is why the checker is red on this Spec right now; rendering the
  derived command into a `qa` Task file is what this feature does, so the fix is
  incomplete until the one file carrying its output is regenerated.
- Assert the property that was missing: the command `DerivedQAVerification`
  produces passes `roundfix spec check` on a Spec that carries it. A generator
  whose output the checker refuses is broken, and nothing tested that.
- Leave `SC-VERIFY-NON-HERMETIC` alone. Its false positive on regex literals is
  real and predates this Spec — it refused an authored command in Spec 0116 the
  same way — but widening a hermeticity check is not this Spec's scope, and a
  derived command that avoids the construct is correct regardless of whether the
  detector is later relaxed.

## References

- `_prd.md` → Goal 3, Core Feature 3
- `_techspec.md` → Interfaces: `DerivedQAVerification`
- The refusal reproduced on this Spec's own `task_06`

## Verification
- `grep -q "TestDerivedQAVerificationPassesTheChecker" internal/spec/task_test.go || exit 1; go test -count=1 ./internal/spec ./internal/speccheck && go test -count=1 -tags docscontract ./internal/docscontract`
