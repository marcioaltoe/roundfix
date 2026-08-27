---
status: completed
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running command.

## Work
- A set-and-passed window refuses `implement` with no Run created
- A 07:00 cutoff set at 23:00 leaves the window open
- A second set without `--force` changes nothing; with `--force` replaces
- A Run created before the cutoff reaches its terminal outcome after it passes
- A crossing Run is created and reports
- A schema-12 database migrates and behaves unchanged without a window
- The glossary carries the term and the vocabulary detector runs

## References

- All user stories and core features

## Verification
- `roundfix spec check 0117-a-window-the-preflight-owns --strict && go test -count=1 ./internal/store ./internal/cli ./internal/watch 2>&1 | grep -q "^ok"`

## Result
QA gate passes all acceptance criteria.
