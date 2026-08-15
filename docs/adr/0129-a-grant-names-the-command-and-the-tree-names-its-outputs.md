---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T08:10:00Z
updated_at: 2026-08-15T08:10:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A grant names the command; the tree names its outputs

ADR-0081 keeps an authorization drawn around the cause, with the computable
effects of the authorized edit following it. The changed-path audit implements
that by asking each grant to enumerate every output of its regeneration command,
which is the enumeration of consequences that decision rejected: the list grows
with the artifact set, drifts the first time a command gains an output, and must
be repeated in every grant that touches the same command. The repository already
records ownership where the artifacts live — `_ownership.yml` names each derived
tree's owner and the command that regenerates it — so a grant names the command
and the audit resolves that command's outputs from the tree. A grant stays one
sentence about a cause; the effects are read from the place that already has to
be true for the Baseline to work at all.
