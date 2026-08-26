---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The gate accepts a manifest that names a catalog that is gone

## Opportunity

`docs/agents/setup-context.json` records the `catalogDigest` that produced the
current managed bytes. Nothing verifies that the recorded digest matches the
catalog the repository actually ships, so the two drift silently.

Measured on 2026-08-10, immediately after merging `#147`: the manifest recorded
`sha256:968c473a…` while the embedded catalog hashed to `sha256:eb010990…`. The
drift arrived because a commit changed `internal/baseline/assets/` after the
last `roundfix baseline update`, and `make baseline-digests` regenerates derived
pins without touching the manifest. Running `roundfix baseline update --yes`
corrected it in one line.

`make verify` passes in that state. Confirmed by stashing the correction and
rerunning the full gate: exit 0.

## Value

The recorded digest is the manifest's only claim about provenance — it answers
"which catalog produced these bytes". When it is stale, that answer is wrong,
and it is wrong in the direction that hides work: a reader concludes the guides
are current for the shipped catalog when they were rendered from an older one.

The Baseline is the mechanism this repository uses to keep the fleet's agent
instructions honest, so a provenance field the gate does not check is the one
field most likely to be believed without evidence. The window is every commit
between a catalog edit and the next `baseline update`, and that window is
routine: the two commands are separate, and only one of them runs in `verify`.

## Shape

The gate could compare the manifest's `catalogDigest` against the embedded
catalog and fail when they differ, the way it already fails on a drifted skill
digest. A softer variant reports drift as a warning from `roundfix baseline
update --format json` without failing, which does not close the window but does
make it visible.

Worth settling in the same work: whether `make baseline-digests` should refresh
the manifest itself, since the drift originates there — the command already owns
every other derived pin, and leaving one out is what splits the two commands.
This shape is non-binding.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
