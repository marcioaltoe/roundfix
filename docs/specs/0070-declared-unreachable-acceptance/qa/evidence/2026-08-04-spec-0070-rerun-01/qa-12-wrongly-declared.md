# QA-12 — Wrongly declared reachable row

Status: pass

The disposable PRD declared `Archive Command help is reachable` unreachable.
The gate exercised `/private/tmp/roundfix-qa70-a10638b archive --help` anyway;
it exited 0 and described declared-only partial acceptance. The fixture QA
state therefore records a wrongly-declared-row finding and `verdict: fail`
rather than a declared block. Archive then exited 2 with the fail-verdict
diagnostic, and a fresh read retained the active Spec with no archive
destination.
