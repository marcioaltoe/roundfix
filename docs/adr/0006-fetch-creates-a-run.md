# Fetch creates a run

Roundfix treats `roundfix fetch` as a tracked Fetch Run in the global Run Database rather than an untracked helper command. This gives fetched markdown artifacts a durable history and keeps Review Source access, Artifact Directory validation, and status inspection on the same path as the watch loop without starting an Agent, committing, or pushing.

Each successful fetch writes a Round artifact directory. With automatic Round selection, Roundfix writes the next available Round instead of mutating an earlier one. If the user selects an explicit Round that already exists, fetch fails instead of overwriting local artifacts.

CodeRabbit fetch stores the Review Source identity in `source_ref`, currently as `thread:<id>,comment:<id>` for inline review threads. That identity lets later resolve runs deduplicate repeated findings across Rounds while preserving the historical artifacts from each fetch.
