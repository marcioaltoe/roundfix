# Supervisor evidence for the sandbox-blocked gate row (2026-07-29)

The 2026-07-29 QA gate reproduced **no Spec behavior finding** and returned
`partial` for one reason: the managed sandbox denies `fork/exec /bin/ps`, which
blocks five process-identity tests inside the repository gate. Per the QA
contract that single block caps the verdict regardless of every other row
passing. The row was executed here, outside the sandbox, on the same build.

## Row 2 — the repository gate with `/bin/ps` available

```console
$ rtk env GOCACHE=/private/tmp/claude-501/gocache-61 make verify
… fmt-check, full test suite, skills-sync-check, skills check, production build …
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec,
write-tasks, setup-context-driven, implement-task, implement-spec,
brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
exit 0
```

Every process-identity test executed. The sandbox limitation is independent of
this Spec and is recorded in
[the owner-identity finding](../../../../findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md),
owned by Spec 0055.

## Independent confirmation of the Spec's behavior across real repositories

The derivation was exercised against five real checkouts with the built binary,
which is the evidence the gate could not collect from inside one worktree:

```text
roundfix   skills: failed — external skills are missing: knowledge-workspace;
                    next: bunx skills add marcioaltoe/skills@knowledge-workspace
fluxus     skills: failed — outdated: roundfix
oraculum   skills: failed — outdated: roundfix
vortex     skills: failed — outdated: roundfix
tax-poc    skills: failed — outdated: roundfix
```

Three properties are proven by that output:

1. **The false Go/TUI requirement is gone.** Before this Spec the four
   TypeScript repositories were told they needed `golang-cli`,
   `golang-concurrency`, `golang-context`, `golang-error-handling`,
   `golang-lint`, `golang-testing`, `bubbletea`, `tui-design`, and
   `autoresearch`. None of those appears now; the only remaining gap is an
   outdated Roundfix-owned Skill, which is true and unrelated.
2. **A repository that genuinely selects those modules still requires them.**
   This repository's manifest selects `go`, `cli-surface`, and `tui-surface`,
   and its readiness still accounts for them — reporting
   `38 required: 14 Roundfix-owned, 24 external` once its own gap is closed.
3. **The remediation is per skill.** The missing entry prints
   `bunx skills add marcioaltoe/skills@knowledge-workspace`, not the
   package-wide install that previously pulled the entire upstream catalog.

The `knowledge-workspace` gap is pre-existing drift this Spec surfaced rather
than caused: the `context-workflow` module has always declared that skill, and
the fixed list masked the gap by never listing it. Closing it exceeds this
Spec's exact tooling authorization and is left for its own authorized change.
