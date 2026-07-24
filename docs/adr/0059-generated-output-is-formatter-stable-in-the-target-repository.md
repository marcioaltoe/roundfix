# ADR-0059: Generated output is formatter-stable in the target repository

Status: Accepted

Immediately after a confirmed apply and clean audit, the target repository's formatter rewrote managed Markdown, leaving setup audit and the repository's Verification mutually dirty. Formatter compatibility is therefore part of the generated-output contract: managed Markdown must survive the repository's selected formatter unchanged, proven by fixtures that compose apply, formatter check, selected Verification, audit, and reapply with no delta. Apply idempotency that does not compose with the repository's own gate would force maintainers to choose between a clean audit and a passing format check.
