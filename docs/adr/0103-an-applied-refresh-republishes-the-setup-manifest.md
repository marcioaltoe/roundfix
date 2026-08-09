---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An applied refresh republishes the Setup Manifest so the update converges

An update that leaves the Setup Manifest describing bytes other than the ones on
disk cannot tell a finished sweep from an unfinished one, and re-proposes the
same work forever. The Baseline Command therefore republishes the Setup Manifest
as part of every applied refresh, and a second run against an unchanged catalog
reports the repository current with no proposed change. This makes the
divergence self-healing: one applied refresh moves a repository whose manifest
predates its managed regions into a state where the recorded digests describe
the bytes, and every later run classifies against a record that is true. The
convergence is the acceptance criterion, not an implementation detail — the
Secondbrain's `convergencia-observavel-em-sistemas-operacionais` states that a
job ends clean only when the system confirms source, desired state, and
destination are aligned, and reporting `current` on the second run is that
confirmation. Treating the manifest as an append-only adoption record was
rejected because it makes the recorded state permanently unfalsifiable.
