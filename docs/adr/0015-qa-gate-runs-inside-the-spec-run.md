# QA gate runs inside the spec run behind an opt-in flag

With the opt-in `--qa` flag, a spec Run ends with a qa-gate step that runs only when every Task in the Task Graph is completed. The Daemon reads the verdict from the QA Report frontmatter — a missing or unreadable verdict counts as fail — ends the Run as Unresolved on a failing verdict, and commits the QA Report in its own commit either way. The verdict frontmatter is a contract addition the upstream qa-gate skill must adopt.
