# QA-14 — Static-failure scoping

Status: pass

The required static command, `rtk make verify`, exited 0. Therefore it
implicated no matrix row and caused zero `blocked (finding: ...)` statuses.
QA-03 through QA-13 and QA-15 through QA-16 all remained runnable and are
being executed independently.

This is the exact conditional behavior Task 06 requires this report to state:
had a named static check failed, only rows whose entry point, observable, or
evidence depended on that check would wait on it. No failed check exists on
this build, so naming an implicated row would be dishonest.
