# Run Event payload stores raw producer JSON

Roundfix records Run activity as Run Events: one ordered record per meaningful occurrence, carrying Run identity, Batch when known, event source, event kind, a bounded text summary, and a structured payload. For Agent events, the payload is the raw ACP `session/update` JSON exactly as the ACP Runtime sent it, not a normalized Roundfix struct. Daemon-owned events define their own small JSON payloads.

We chose raw payloads because rendering fidelity must be able to improve retroactively. A journaled Run cannot be re-executed, so anything the conversion layer drops at write time is lost forever; the earlier normalized stream model already discarded diff old/new text, ACP tool kinds, and plan priorities that the Live Run View needs for rich rendering. With raw payloads, richer rendering is a reader-side change that also applies to Runs journaled before the change. Raw payloads also tolerate ACP protocol growth: readers skip unknown event kinds instead of failing, so journals written by newer SDK versions stay readable.

The trade-off is that SQL cannot query inside Agent payloads without JSON functions. Normalized columns exist only for what list views filter on: kind, source, Batch, Review Issue reference, tool ID, tool state, summary, and timestamp. The payload is write-once and read-as-blob.

Consequences: producers must never re-serialize or prune payload JSON; the bounded summary, not the payload, serves list rendering; and replay renderers must treat unknown kinds and unknown payload fields as skippable, never as errors.
