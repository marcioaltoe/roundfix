# Reproduce the focused evidence

Run from the repository root of audited head 95133cf3c91c7885bc27484b47759e76fe95376b. The `.py.txt` files retain the local harness without adding compiled repository sources. Copy runner, sabotage, positive, cli and catalog to `/tmp/qa0333/` with `.py` extensions; copy scope there for the history audit. The helper uses this repository and its existing Go caches, with network resolution disabled and fresh HOME/TMPDIR for every CLI command.

Run `python3 /tmp/qa0333/positive.py`, then `python3 /tmp/qa0333/sabotage.py`, then `python3 /tmp/qa0333/cli.py`. The sabotage harness initializes and stages Git metadata before each mutation. Run `python3 /tmp/qa0333/catalog.py` and `python3 /tmp/qa0333/scope.py` for the independent file/history assertions. Every subprocess is waited on; fixtures are removed after restoring write permission. Keep output in temporary storage until the suites exit before copying evidence into the checkout.

The seven fixture-setup transcripts document the initial archive-without-Git refusal. Their exits are not credited as mutation detections. The subsequent sabotage transcripts show the named assertions rejecting the mutations. Row 14 reached its assertion in its initial copy and was not repeated.

The documentation gate uses the retained temporary go-timeout wrapper: it delegates normal Go commands unchanged and adds `-timeout 20m` to every test invocation. Repository configuration was not edited. The Daemon-supplied repository Verification must not be rerun by a QA child Agent.
