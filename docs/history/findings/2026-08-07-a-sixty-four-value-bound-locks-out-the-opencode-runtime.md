---
status: done
absorbed_by: 0088-a-third-runtime-that-can-run
created_at: 2026-08-07
updated_at: 2026-09-08
kind: finding
---

# A sixty-four-value bound locks out the opencode runtime (2026-08-07)

`opencode` is one of Roundfix's three supported runtimes and no profile can use it. Its adapter advertises 431 models; Roundfix's capability validator caps a select option at 64 values and refuses the evidence.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-07-a-sixty-four-value-bound-locks-out-the-opencode-runtime.md`.
