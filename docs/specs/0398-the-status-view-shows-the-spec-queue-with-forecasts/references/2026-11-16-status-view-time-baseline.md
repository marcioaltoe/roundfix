# Baseline: how long the status view takes, before Spec 0398 (recorded by the author, 2026-11-16, about 10:35 BRT)

The maintainer's standing condition, made for Spec 0397 on 2026-11-16, is no or almost no loss of performance. This file records what `roundfix follow-up --status` costs before this Spec, so that the release that carries it has a number to be compared with. Nothing was changed to take it.

## 1. The measurement

- Command: `roundfix follow-up --status --format json`, five runs in a row, each under `/usr/bin/time -p` and `taskpolicy -c utility`, the policy the maintainer's shell gives Roundfix on this machine.
- Binary: the installed Roundfix 0.145.1. Checkout: this repository's default branch at `f4580437`, clean.
- State the view read: a Run Database of 316 MB with a write-ahead log of 4 MB, one active Implement Run, one Delivery Queue item, one active Spec and no open pull request. Every run exited `0` and reported each of its six sources as `checked`.
- Machine: the Mac mini, 10 logical CPUs, while an Implement Run of this repository was in progress.

| Run | Wall time (s) | User (s) | System (s) |
| --- | --- | --- | --- |
| 1 | 1.11 | 0.18 | 0.24 |
| 2 | 0.90 | 0.18 | 0.18 |
| 3 | 0.79 | 0.17 | 0.18 |
| 4 | 0.83 | 0.16 | 0.16 |
| 5 | 0.71 | 0.16 | 0.16 |

Wall time: median 0.83 s, from 0.71 to 1.11. Processor time, user plus system: median 0.35 s, from 0.32 to 0.42.

## 2. The band used for a comparison

The five wall times have a standard deviation of 0.15 s. The standard error of a median of ten runs is taken as `1.2533 × 0.15 ÷ √10`, the normal approximation, which gives 0.06 s. A difference under about 0.1 s between two binaries, ten runs each, cannot be told from chance.

## Limits

- The wall time includes one `gh pr list` call over the network, and the processor time includes that child process.
- The first run was the slowest. Nothing was done to warm a cache or to clear one.
- It is one machine at one moment, with other work running on it.
- Five runs give a rough spread. The comparison at release uses ten runs of each binary, alternated, on the same checkout and Run Database.
