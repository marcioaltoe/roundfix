# Fetch creates a run

Roundfix treats `roundfix fetch` as a tracked Fetch Run in the global Run Database rather than an untracked helper command. This gives fetched markdown artifacts a durable history and keeps Review Source access, Artifact Directory validation, and status inspection on the same path as the watch loop without starting an Agent, committing, or pushing.
