# Cleanup regeneration evidence

QA cloned the audited commit without hard links to
`/private/tmp/roundfix-0130-qa.MsBR7M/repo`; `rtk git rev-parse HEAD` returned
`c6a5ea156c70d33d59e9e5e1d98fcc2fe24cb318`.

To reproduce the PR #180 post-cleanup precondition without changing the Run
Worktree, QA ran `rtk git rm -r docs/workflow` in that disposable clone. Git
staged 45 deletions and no commit was created.

With `GOCACHE=/private/tmp/roundfix-0130-qa.MsBR7M/gocache`, two consecutive
`rtk make baseline-digests` invocations each exited 0 and ended with:

```text
baseline-digests: no changes; derived artifacts already match their canonical sources
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}
```

The independent rereads `test ! -e docs/workflow` and
`rtk git diff --exit-code` both exited 0. The first proves regeneration did not
restore the removed directory; the second proves it made no unstaged change over
the staged cleanup deletion.

