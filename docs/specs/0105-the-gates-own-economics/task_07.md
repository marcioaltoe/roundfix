---
status: completed
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

## Result

Implementation:

- `DerivedQAVerification` now obtains the report basename by splitting on `/`
  and validates filenames and verdict whitespace with string operations. The
  generated command contains no slash-leading awk regex token for the
  hermeticity checker to misclassify.
- The command still parses calendar dates, leap years, and numeric rerun
  sequences before sorting. The existing selection cases now also cover a
  malformed rerun suffix, a non-digit date, and surrounding verdict whitespace.
- `task_06.md` carries the corrected derived command.
- `TestDerivedQAVerificationPassesTheChecker` exercises the real Task-stage
  checker through its existing `internal/speccheck` suite. The checker-side
  fixture verifies both that the command is independently hermetic and that a
  `qa` Task carrying it has neither `SC-QA-VERIFICATION-AUTHORED` nor
  `SC-VERIFY-NON-HERMETIC` findings.
- `SC-VERIFY-NON-HERMETIC` implementation code was not changed. This slice
  introduced, changed, or retired no glossary term, so `CONTEXT.md` needs no
  update.

Focused checks:

- Before the generator change,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go test -count=1 ./internal/spec -run '^TestDerivedQAVerificationPassesTheChecker$'`
  failed because the checker classified `/^.*\\/qa-report-/` as an external
  path. This reproduced the missing contract at the named seam.
- After the change,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go test -count=1 ./internal/spec -run '^TestDerivedQAVerification(RequiresTheNewestReportToPass|PassesTheChecker)$'`
  passed. This covers checker acceptance, parsed date and sequence ordering,
  verdict-domain enforcement, malformed names, and missing reports.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go test -count=1 -tags docscontract ./internal/docscontract -run '^TestCheckActiveCorpusHasNoErrors$'`
  passed after `task_06.md` was re-rendered. This checks the corrected command
  against the active Spec corpus without running the Task's declared
  Verification sequence.

The Daemon-owned Verification command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260831T183328Z_c381ee928ef8acd5`
- Source commit: `75339d0f6c469f8c892a802f382de9c23a82aabf`
