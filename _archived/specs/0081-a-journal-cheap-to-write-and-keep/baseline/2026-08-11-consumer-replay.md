# Pre-change Run Event Journal consumer replay — 2026-08-11

The recorded SQLite journal was produced by a build of pre-Spec commit
`a2a4c86b7570e4ef782ccc2ff390033f586d47fd` through that build's production
Journal Sink. It contains one terminal Implement Run and six ordered Run
Events spanning Agent payload rendering, Task state, Verification verdict,
Daemon status, and terminal outcome. No production source was changed to
record or replay it.

## Recorded artifacts

| Artifact | Purpose | SHA-256 |
| --- | --- | --- |
| `2026-08-11-prechange-roundfix.db` | SQLite Run Event Journal written by the pre-change build | `1c38ed391769fb2c3ba5e48875ed868737bf0bc5659c12bdbf7ee8b19f5c5dec` |
| `2026-08-11-prechange-consumer-expectations.json` | Pre-change `events`, Attach, reconcile, and `gc` observations | `3c4eb9511ee29dda3116ac8863e8b7f0b5a74fcaf594470ce6f0788b7f5ab04b` |
| `2026-08-11-prechange-cockpit.golden` | Full color-free 120x40 Cockpit rendering | `16eb1d01eb12a55728e0343c12768f04e40809359acc8214c368a42e5a8d9700` |

The generator source is `task10-corpus-generator_test.go.txt`. Installed as
`internal/store/task10_corpus_generator_test.go` in an archive of the named
commit, this exact command wrote the journal:

```sh
rtk proxy env GOCACHE=/private/tmp/roundfix-task10-gocache-before TASK10_CORPUS_OUTPUT=/private/tmp/roundfix-task10/prechange-roundfix.db go test -count=1 ./internal/store -run '^TestTask10RecordPrechangeJournal$' -v
```

The two observation harnesses are
`task10-cli-consumer-harness_test.go.txt` and
`task10-cockpit-consumer-harness_test.go.txt`. The pre-change build recorded
their expectations; `TestJournalConsumerCorpusReplaysEveryConsumer` installs
the same sources through Go overlays and compares the current packages with
those artifacts.

## Per-consumer comparison

| Consumer | Pre-change observation | Current-build comparison |
| --- | --- | --- |
| `events` | Four Supervisor records: Task started, Verification passed, Task settled completed, and Integration Pending outcome; exact schema fields and values are in the JSON artifact. | Decoded records are identical. |
| Attach | Five timeline rows, in cursor order, ending with `Run reached IntegrationPending.` | Timeline rows are identical. |
| Cockpit rendering | Complete color-free 120x40 Attach-mode render from the recorded journal. | Render is byte-identical to the golden file. |
| Reconcile replay detection | `task_02` is completed, with both passed Verification and settled-completed evidence decoded from journal payloads. | Completed Task and both evidence flags are identical. |
| `gc` | At fixed cutoff `2021-01-01T00:00:00Z`, the terminal Run is selected and all six journal rows are pruned from a disposable copy. | Run selection and six-row prune result are identical. |

The comparison covers observable consumer output and mutation on disposable
copies of the recorded journal. The source fixture remains byte-identical for
every replay.
