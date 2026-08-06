# The QA gate proves machine facts before it spends an agent turn

The QA gate runs a Daemon-owned mechanical stage before its Agent Session: a
deterministic, hermetic, citation-only pass that computes the facts the gate
performs by hand today — the tooling authorization's bounded files against each
Task commit's actual changed paths, consequent-fix commit shape, report
structural shape, evidence-path resolution, and the repository verification
gate — reports every failure it finds in one pass rather than the first, and
withholds the Agent Session when a blocking fact is present. The stage
computes; ADR-0080 keeps sole ownership of verdict semantics and the three
typed blocked-cause counts, and no stage may make a verdict more permissive.

The gate stays one terminal `qa` Task node under ADR-0091: the stage is a step
inside it, never a command outside the graph, because a stage that could run
independently would reintroduce the per-run choice ADR-0088 removed. Placing
it with the Daemon follows ADR-0014, which already gives the Daemon ownership
of Verification execution, and keeps its cost and failures visible without an
Agent turn. Measured across Spec 0079's first three Runs, the first took 92
minutes for seven Tasks plus the gate, and the two corrective Runs took 29 and
30 minutes. Those Runs exposed three distinct defects; two were facts a Git
read or an exit status already knew.
