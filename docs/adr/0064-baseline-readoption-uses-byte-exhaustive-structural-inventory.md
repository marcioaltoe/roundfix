# ADR-0064: Baseline Readoption uses byte-exhaustive structural inventory

Status: Accepted

Baseline Readoption partitions every nonblank byte in bounded agent-instruction carriers into deterministic Source Baseline Entries, then requires an explicit classification and individual disposition for each entry before mutation. The scanner never infers normative meaning: exact source bytes, ordered entry digests, decisions, destinations, reasons, and proposed Repository-Specific Normative Rules all enter the Change Plan digest, trading a larger decision set for omission-proof, stale-safe adoption of unknown historical states.
