# ADR-0046: Setup-owned agent instructions are declarative

Status: Accepted

The setup workflow records selected profiles, modules, decisions, and managed artifacts in a versioned manifest and wraps generated Markdown in stable ownership markers. Updates resolve these identifiers instead of inferring intent from prose, preserve repository-authored content outside managed boundaries, and request confirmation only when a decision identifier is missing or incompatible. This makes audit and safe correction deterministic, portable, and idempotent across template revisions.
