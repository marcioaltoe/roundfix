### resolve

```bash
roundfix resolve --pr <number> [--spec <slug>]
```

Works only over downloaded Compatible Artifacts — it does not fetch. It
assigns bounded Batches to the Agent, verifies each Batch, and commits
successful Batches directly on the PR Head Branch when auto-commit is enabled.
At Batch settlement each Review Issue propagates to the Review Source
individually: `resolved` threads resolve; `invalid` and `duplicated` threads
get an explanatory Outcome Comment and then resolve; `failed` threads get the
failure reason and stay open; issues still unresolved at Run end receive a
closing comment. Comments carry an idempotency marker, so retries never
duplicate. The Run's review artifacts are committed in one separate docs
commit (`docs: review round NNN for pr <n>`, ADR-0036) and Final Push runs
only when no Unresolved Review Issues remain. Artifact roots outside the
repository, or reached through a symbolic link, are reported and never staged.

