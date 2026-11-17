# Boundary Notices in the operator's Phase Record, 2026-11-16 to 2026-11-17

Measured by the author on 2026-11-17 on the Mac mini, read-only, from the two day files of the operator's Phase Record (`~/.roundfix/phase-record/2026-11-16.jsonl` and `2026-11-17.jsonl`). The record started with the release of Spec 0397 (Roundfix 0.146.0) on 2026-11-16, so these are all of its phase entries up to the measurement. The files also hold 20 prompt entries, which no Boundary Notice reads and this measurement leaves out.

## Method

The author listed every entry whose `kind` is `phase`, in the order of its `time`, with its phase, the number of its Spec, the last segment of its `repository`, its `boundary` mark and its `previous_phase` and `previous_spec`. Every entry carries the same Claude Code session: the operator's Supervisor session, which started each Authoring Phase itself or through a subagent that shares its session identifier. The session identifier and the full repository paths are left out of this record. The last column applies the rule this Spec adopts: an entry is marked only when the same session's latest earlier entry is `release` and the new entry is not a `release` of the same Spec.

## Table

| # | time (UTC) | phase | Spec | repository | marked by Roundfix 0.147.0 | previous | marked under this Spec's rule |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 2026-11-16T21:17:48Z | follow-up | - | roundfix | no | - | no |
| 2 | 2026-11-16T23:15:41Z | techspec | 0399 | roundfix | no | - | no |
| 3 | 2026-11-16T23:24:11Z | techspec | 0399 | roundfix | no | - | no |
| 4 | 2026-11-16T23:49:13Z | follow-up | - | roundfix | yes | techspec of 0399 | no |
| 5 | 2026-11-16T23:56:12Z | tasks | 0399 | roundfix | no | - | no |
| 6 | 2026-11-17T00:16:21Z | techspec | 0400 | plan0400 | yes | tasks of 0399 | no |
| 7 | 2026-11-17T00:24:41Z | tasks | 0400 | plan0400 | yes | techspec of 0400 | no |
| 8 | 2026-11-17T09:38:43Z | techspec | 0401 | plan0401 | yes | tasks of 0400 | no |

## Results

- Eight phase entries, one session, three Specs authored, no delivery and no release recorded yet.
- Roundfix 0.147.0 printed the Boundary Notice four times: entries 4, 6, 7 and 8. Each followed an Authoring Phase; none followed a release.
- Entry 5 printed none only because entry 4 was a `follow-up`, which the 0.147.0 rule exempts as an earlier phase.
- Under this Spec's rule none of the eight entries is marked, because no earlier entry is a `release`.
- Entries 2 and 3 repeat the same phase of the same Spec, which no rule marks.

The maintainer's message of the same evening, "Esse negocio de criar sessões a cada etapa está atrasando demais.", came while entries 1 to 4 were being written.
