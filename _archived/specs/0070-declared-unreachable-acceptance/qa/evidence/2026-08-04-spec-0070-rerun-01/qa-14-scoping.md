# QA-14 — Static-failure scoping

Status: pass

The first full-gate attempt failed in one unrelated integration-test wait after
3,286 passing tests. While unclassified, it implicated QA-02 only: the QA
binary used by archive journeys had already built independently, and the
governance, documentation, and public scratch entry points did not depend on
that Agent-start timing assertion. Those rows continued and retained their own
results instead of becoming finding-blocked.

Unchanged-build reproduction classified the failure as transient gate load:
the exact test passed 1/1 and 10/10, then the controlled full gate passed all
3,287 tests. QA-02 therefore also closes pass. No row is finding-blocked and no
failure cause is folded into a blocked-row count.
