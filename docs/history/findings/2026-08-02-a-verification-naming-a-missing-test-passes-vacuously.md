---
status: done
created_at: 2026-08-02
updated_at: 2026-09-08
absorbed_by: 0095-a-verification-that-ran-before-anyone-believed-it
---

# A Verification naming a missing test passes vacuously (2026-08-02)

`go test ./pkg -run TestThatWasNeverWritten -count=1` exits **0**. The Daemon runs each Task's declared Verification verbatim, reads exit 0, and settles the Task `completed`. A Task whose Agent implemented nothing therefore settles `completed` as long as its Verification names tests that do not exist.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-02-a-verification-naming-a-missing-test-passes-vacuously.md`.
