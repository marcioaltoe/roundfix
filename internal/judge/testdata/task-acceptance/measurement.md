# Task acceptance fixture measurement

This is an offline fixture, not a live measurement. Eight implementation Tasks
from four fixture Specs receive predetermined fake-transport answers. The
record exercises ties, both labels, whole-Spec bootstrap draws and baselines.

Verdict: inconclusive

The decision rule requires at least 150 answered Tasks and 20 repaired Tasks.
AUROC is zero when no positive/negative pair exists; interval bounds are zero
when every bootstrap draw is discarded. Those cases remain inconclusive.
Both PCG streams use their question-file seed for the state and sequence.
Random scores follow canonical Spec/Task order.

The sibling sabotaged record changes only AUROC by 0.1. Task files under
`repo/` supply the title, criteria and state hashes for the consistency check.
