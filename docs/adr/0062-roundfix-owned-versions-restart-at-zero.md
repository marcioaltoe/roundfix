# ADR-0062: Roundfix-owned versions restart at zero

Status: Accepted

Roundfix restarts the CLI, Context-Driven Baseline, schemas, manifests, and every Roundfix-owned distributed skill at `0.0.1`, without compatibility obligations for prior setup identifiers or transitions. Existing operational configuration, Runs, and Run Database state remain intact; upstream-managed skill versions and Git history remain authoritative, while changelog, tags, and GitHub Releases restart only through the confirmation-gated Release Plan.
