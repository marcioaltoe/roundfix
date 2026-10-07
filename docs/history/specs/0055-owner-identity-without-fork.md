---
schema: roundfix/archive-record/v1
spec: 0055-owner-identity-without-fork
title: Owner identity without fork
status: archived
created: "2026-07-28"
archived: "2026-07-31"
disposition: pass
source: docs/history/specs/0055-owner-identity-without-fork
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-31.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "55"
delivery_commit: 37d57b2266e746f78fa70a5acccad210e0e29e1f
---

# Owner identity without fork

Spec 0037 gave Force Stop a real ownership proof: Runs record an opaque owner start-time identity, and Force Stop refuses to signal a PID whose live identity does not match. The proof is obtained by forking `/usr/bin/ps` on every read, so it fails exactly when the host cannot fork — which is precisely when a runaway Run has exhausted the machine and the escape hatch is needed most. Under the same pressure, a failed capture at Run creation silently records no identity, so a Run can run without reuse protection with no warning. Separately, `roundfix stop <run-id> --force` rejects its trailing flag — the same argument-ordering defect Spec 0042 fixed for the Attach Command. Evidence: [owner identity forks ps and fails closed under load](../../findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md); the Stop Command flag defect was observed during the Spec 0042 QA recovery on 2026-07-28, alongside the Attach fix.
