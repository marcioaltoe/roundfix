# Cockpit reads the Run Event Journal, never the sink

The interactive Live Run View consumes Run Events exclusively from the Run Event Journal — cursor-paged replay plus data-version follow polling — even while the owning process is producing those events live. The best-effort sink path never feeds the cockpit; it remains for non-TTY text output only.

We chose this because the journal is already a critical sink in every resolve and watch cycle, so the durable source exists whenever a cockpit exists, and a single cursor-ordered source makes live and replayed rendering identical by construction. The alternative — pushing the live tail through the sink and paging the journal only for scrollback — requires stitching two sources at a moving boundary, which invites duplicate and gap bugs that cursor-only consumption rules out by design.

Consequences: the live tail lags by up to one poll interval (~250ms), which is acceptable for this product; every cockpit holds a read-only database connection alongside the writer in the same process (supported under WAL and already covered by store tests); and the TUI's bounded event-buffer sink adapter is no longer the cockpit delivery path.
