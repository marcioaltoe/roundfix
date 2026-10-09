# Published passage analysis

Source: Candido, Melo and d’Amorim, Test Suite Parallelization in Open-Source Projects, ASE 2017, https://damorim.github.io/publications/candido-etal-ase17.pdf.

Requirement 7 names the passage quoted in the PRD: the published results report 97.5% null dereferences, 1.6% concurrent access on unsynchronized structures, and 0.8% likely broken dependencies among observed failures. This is an offline analysis of the supplied published quotation, not a fresh PDF retrieval. The grant allows optional network reads of public CI logs only; no other network request was sent.

The passage supports checking shared state and ordering when parallelizing tests. It does not establish a failure rate for Go or show that race detection catches null dereferences. Those percentages cannot prove this conversion safe; repository-specific race and shuffle runs supply that proof.
