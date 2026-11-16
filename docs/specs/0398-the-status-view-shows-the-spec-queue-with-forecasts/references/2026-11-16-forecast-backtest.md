# Backtest of the queue forecast over this repository's history (recorded by the author, 2026-11-16, 10:50–11:40 BRT)

The maintainer accepted a forecast only if it holds on history: "No PRD, o método é repetido sobre ~100 Specs dos últimos 30 dias e só é aceito se o fim real cair dentro do limite em pelo menos 80%." This file records the replay: its inputs, the method, what held, what was refused, and the history it ran on. Nothing was changed to take it.

## 1. Inputs

- Runs: `roundfix runs list --state terminal --limit 0 --format plain` in this repository's checkout on the Mac mini, with the installed Roundfix 0.145.1, at 2026-11-16 13:52 UTC. It listed 165 terminal Implement Runs of 110 Specs, started from 2026-10-17 12:10 UTC to 2026-11-16 11:30 UTC: 108 Clean, 56 Unresolved and 1 Failed.
- Git history of the default branch at `f4580437`: the commit that added each Spec's `_prd.md` and the commit that added its Archive Record, each with its committer time.
- 109 of the 110 Specs are delivered. The other, Spec 0397, was in delivery when the listing was taken.

Three times describe a delivered Spec. Its plan is the committer time of the commit that added its PRD to the default branch. Its start is the start of its first Implement Run. Its end is the committer time of the commit that added its Archive Record, which is its delivery merge. Every one of the 109 Specs has the three, in that order.

Two durations are measured from them, in hours:

| Duration | Specs | Median | 85th percentile | Longest |
| --- | --- | --- | --- | --- |
| Start to merge: end minus start | 109 | 1.94 | 4.36 | 15.44 |
| Start gap: a Spec's start minus the start of the Spec that started before it, counted when its plan was already on the default branch at that earlier start | 74 | 2.05 | 3.88 | 21.59 |
| Plan to merge: end minus plan (refused below) | 109 | 6.28 | 21.19 | 61.94 |

## 2. The method

- A percentile is the nearest rank: of `n` sorted values, the value at position `ceil(q × n)`. The likely time is the 50th percentile and the limit the 85th.
- A forecast needs at least 11 values. A forecast re-made after time has passed uses only the values greater than that time and needs at least 5 of them.
- The sum of two durations is the set of every pairwise sum. A set of more than 512 values is cut to 512: of `n` sorted values, the values at positions `ceil(k × n ÷ 512)` for `k` from 1 to 512.
- A Spec in progress: its end is its start plus the start-to-merge time.
- A planned Spec that has not started, one, two or three rows behind a Spec in progress: its start is the start of that Spec plus one start gap for each row, and its end is its start plus the start-to-merge time.

The replay is walk-forward: each forecast uses only the Specs whose delivery merge came before the instant it was made, and the gaps of the Specs that had started by then.

## 3. What held

Each line is one kind of forecast, made at the instant named, for every Spec with enough history before it.

| Forecast | Made | Cases | Limit held | Likely time held | Median of limit ÷ actual |
| --- | --- | --- | --- | --- | --- |
| End of a Spec in progress | at its start | 98 | 80 (81.6%) | 46 (46.9%) | 1.91 |
| End of a Spec in progress | 1 hour after its start | 95 | 77 (81.1%) | 43 (45.3%) | 1.91 |
| End of a Spec in progress | 2 hours after its start | 48 | 41 (85.4%) | 22 (45.8%) | 1.72 |
| End of a Spec in progress | 3 hours after its start | 25 | 22 (88.0%) | 12 (48.0%) | 1.65 |
| End of a Spec in progress | 4 hours after its start | 15 | 12 (80.0%) | 6 (40.0%) | 1.86 |
| Start of a Spec one row behind | when the Spec ahead starts | 63 | 55 (87.3%) | 35 (55.6%) | 2.29 |
| End of a Spec one row behind | when the Spec ahead starts | 63 | 53 (84.1%) | 27 (42.9%) | 1.67 |
| Start of a Spec two rows behind | when the first Spec starts | 39 | 36 (92.3%) | 21 (53.8%) | 2.23 |
| End of a Spec two rows behind | when the first Spec starts | 39 | 36 (92.3%) | 20 (51.3%) | 1.84 |
| Start of a Spec three rows behind | when the first Spec starts | 19 | 19 (100%) | 8 (42.1%) | 2.00 |
| End of a Spec three rows behind | when the first Spec starts | 19 | 17 (89.5%) | 10 (52.6%) | 1.88 |

Every kind holds its limit for at least 80% of its cases, its likely time for 40% to 60%, and its limit stays under 2.5 times the actual duration at the median. The re-made forecast of a Spec in progress gave no time for 3 Specs at 3 hours and 4 at 4 hours, where fewer than 5 earlier Specs had taken longer.

A Spec planned after the Spec ahead had started was forecast at its plan, from the gaps longer than the time already passed: 6 cases, the limit held for 5 on the start and 5 on the end. Six cases decide nothing.

## 4. What was refused

| Basis | Cases | Limit held | Likely time held | Why it is refused |
| --- | --- | --- | --- | --- |
| Plan to merge, forecast when the plan reaches the default branch | 96 | 74 (77.1%) | 40 (41.7%) | Under 80%, and the limit is 2.82 times the actual at the median. With another Spec ahead at the plan it held for 55 of 75 (73.3%); with none ahead it held for 19 of 21 with the likely time holding for 17, which is no median. |
| Start of a planned Spec with nothing in progress ahead, from the time it became ready | 78 | 59 (75.6%) | 35 (44.9%) | Under 80%: this wait ends when somebody queues the Spec. |
| End of the same Spec, from the time it became ready | 78 | 60 (76.9%) | 28 (35.9%) | Under 80%, and the likely time is early. |
| End of a Spec behind a Spec in progress, adding the medians and adding the limits | 63, 39 and 19 | 55, 36 and 17 | 25 (39.7%), 14 (35.9%) and 6 (31.6%) | The limit holds and the likely time comes earlier with each row: the median of a sum is later than the sum of the medians. The pairwise sums of section 3 keep it between 42.9% and 52.6%. |
| End of a Spec in progress, from Specs of the same size (three Tasks or fewer, or more) | 98 | 81 (82.7%) | 45 (45.9%) | It holds, with a wider limit (2.10 against 1.91); the plainer basis is kept. |

A Spec becomes ready when its plan is on the default branch and the Spec that started before it is merged. The refused bases share one cause: each of the first three spans a wait that ends when a person acts, and those waits ran from one minute to twenty hours in this history.

The window of the history changed little. With the last 20, 30, 40 or 60 delivered Specs, or the last 7 or 14 days, the end of a Spec in progress held for 78 to 81 of 98. The same windows left the plan-to-merge basis between 72 and 78 of 96.

## Limits

- One repository on one machine, over 30 days. Specs 0365 to 0375 ran on another machine and are not in it.
- The Spec ahead is taken as the Spec that started before, because a Delivery Queue's order is not kept once a new queue replaces it.
- The replay tests a forecast at fixed instants. The status view makes one whenever it is asked.
- A row more than three behind a Spec in progress was not tested: 19 cases reach three.
- The same history chose the method and tested it. The replay is walk-forward, so no forecast saw its own outcome, and the operator repeats the count on the first 20 Specs delivered after the release.

## 5. The history

Times are UTC.

| Spec | Plan on the default branch | First Implement Run | Delivery merge |
| --- | --- | --- | --- |
| 0275 | 2026-10-16T17:51:18Z | 2026-10-17T12:10:21Z | 2026-10-17T15:02:53Z |
| 0276 | 2026-10-17T15:14:12Z | 2026-10-17T15:33:14Z | 2026-10-17T19:55:05Z |
| 0277 | 2026-10-17T21:06:03Z | 2026-10-17T21:06:59Z | 2026-10-18T12:33:38Z |
| 0278 | 2026-10-18T14:16:39Z | 2026-10-18T14:17:53Z | 2026-10-18T15:51:02Z |
| 0279 | 2026-10-18T17:59:53Z | 2026-10-18T18:00:46Z | 2026-10-18T19:57:31Z |
| 0280 | 2026-10-18T20:36:41Z | 2026-10-18T20:37:33Z | 2026-10-18T22:50:17Z |
| 0281 | 2026-10-19T20:20:21Z | 2026-10-19T20:31:12Z | 2026-10-19T22:10:37Z |
| 0282 | 2026-10-19T20:20:21Z | 2026-10-19T22:40:06Z | 2026-10-20T00:40:18Z |
| 0283 | 2026-10-19T20:20:21Z | 2026-10-20T01:08:30Z | 2026-10-20T02:21:26Z |
| 0284 | 2026-10-20T02:45:10Z | 2026-10-20T02:46:02Z | 2026-10-20T03:29:08Z |
| 0285 | 2026-10-20T10:26:06Z | 2026-10-20T10:27:58Z | 2026-10-20T11:44:28Z |
| 0286 | 2026-10-20T10:26:06Z | 2026-10-20T12:12:02Z | 2026-10-20T14:02:25Z |
| 0287 | 2026-10-20T10:26:06Z | 2026-10-20T14:41:17Z | 2026-10-20T16:43:07Z |
| 0288 | 2026-10-20T17:09:12Z | 2026-10-20T17:12:24Z | 2026-10-20T18:25:46Z |
| 0289 | 2026-10-20T17:09:12Z | 2026-10-20T19:05:25Z | 2026-10-20T20:24:09Z |
| 0290 | 2026-10-21T00:48:31Z | 2026-10-21T00:50:25Z | 2026-10-21T02:34:43Z |
| 0291 | 2026-10-21T00:48:31Z | 2026-10-21T03:36:12Z | 2026-10-21T05:43:13Z |
| 0292 | 2026-10-21T00:48:31Z | 2026-10-21T14:38:08Z | 2026-10-21T18:13:54Z |
| 0293 | 2026-10-21T00:48:31Z | 2026-10-21T19:10:41Z | 2026-10-21T20:55:40Z |
| 0294 | 2026-10-21T15:16:05Z | 2026-10-21T21:51:59Z | 2026-10-21T23:07:33Z |
| 0295 | 2026-10-21T15:16:05Z | 2026-10-21T23:47:48Z | 2026-10-22T01:10:02Z |
| 0296 | 2026-10-21T15:16:05Z | 2026-10-22T01:42:06Z | 2026-10-22T03:16:22Z |
| 0297 | 2026-10-22T13:58:10Z | 2026-10-22T14:14:24Z | 2026-10-22T15:48:07Z |
| 0298 | 2026-10-22T13:58:10Z | 2026-10-22T16:31:07Z | 2026-10-22T18:11:49Z |
| 0299 | 2026-10-22T13:58:10Z | 2026-10-22T18:58:27Z | 2026-10-22T23:25:35Z |
| 0308 | 2026-10-22T23:33:29Z | 2026-10-22T23:34:45Z | 2026-10-23T00:51:55Z |
| 0300 | 2026-10-22T13:58:10Z | 2026-10-23T01:49:35Z | 2026-10-23T04:06:39Z |
| 0301 | 2026-10-22T23:54:39Z | 2026-10-23T04:56:15Z | 2026-10-23T06:28:08Z |
| 0302 | 2026-10-22T23:54:39Z | 2026-10-23T07:03:17Z | 2026-10-23T08:59:47Z |
| 0303 | 2026-10-22T23:54:39Z | 2026-10-23T13:28:24Z | 2026-10-23T19:51:47Z |
| 0304 | 2026-10-23T13:26:26Z | 2026-10-23T20:52:40Z | 2026-10-23T22:17:52Z |
| 0305 | 2026-10-23T13:26:26Z | 2026-10-23T23:03:27Z | 2026-10-24T00:11:01Z |
| 0306 | 2026-10-23T13:26:26Z | 2026-10-24T00:40:42Z | 2026-10-24T03:23:35Z |
| 0311 | 2026-10-24T11:55:11Z | 2026-10-24T11:56:10Z | 2026-10-24T13:44:11Z |
| 0312 | 2026-10-24T11:55:11Z | 2026-10-24T14:45:41Z | 2026-10-24T17:06:46Z |
| 0313 | 2026-10-24T11:55:11Z | 2026-10-25T12:21:15Z | 2026-10-25T22:38:10Z |
| 0317 | 2026-10-25T12:05:23Z | 2026-10-25T13:32:32Z | 2026-10-25T17:43:21Z |
| 0318 | 2026-10-25T13:00:55Z | 2026-10-25T18:11:08Z | 2026-10-25T19:06:03Z |
| 0314 | 2026-10-24T11:55:11Z | 2026-10-25T22:40:26Z | 2026-10-26T00:51:57Z |
| 0315 | 2026-10-24T15:14:07Z | 2026-10-25T23:19:38Z | 2026-10-26T02:36:54Z |
| 0316 | 2026-10-24T15:14:07Z | 2026-10-25T23:32:10Z | 2026-10-26T00:09:57Z |
| 0307 | 2026-10-23T13:26:26Z | 2026-10-26T00:11:44Z | 2026-10-26T03:22:52Z |
| 0309 | 2026-10-24T12:32:01Z | 2026-10-26T09:59:05Z | 2026-10-26T11:12:00Z |
| 0310 | 2026-10-24T12:32:01Z | 2026-10-26T11:13:02Z | 2026-10-26T12:40:09Z |
| 0319 | 2026-10-26T18:31:10Z | 2026-10-26T18:32:55Z | 2026-10-26T19:49:29Z |
| 0320 | 2026-10-26T18:31:10Z | 2026-10-26T19:50:26Z | 2026-10-26T21:16:17Z |
| 0321 | 2026-10-26T19:26:53Z | 2026-10-26T21:53:11Z | 2026-10-27T02:46:00Z |
| 0322 | 2026-10-26T19:26:53Z | 2026-10-26T22:51:06Z | 2026-10-27T04:22:01Z |
| 0325 | 2026-10-26T21:12:42Z | 2026-10-26T23:24:17Z | 2026-10-27T05:07:57Z |
| 0323 | 2026-10-26T19:59:30Z | 2026-10-27T00:00:26Z | 2026-10-27T01:11:27Z |
| 0324 | 2026-10-26T19:59:30Z | 2026-10-27T01:12:26Z | 2026-10-27T05:24:56Z |
| 0326 | 2026-10-27T01:18:10Z | 2026-10-27T05:28:00Z | 2026-10-27T09:15:24Z |
| 0327 | 2026-10-27T02:41:48Z | 2026-10-27T05:51:39Z | 2026-10-27T07:05:05Z |
| 0328 | 2026-10-27T02:41:48Z | 2026-10-27T07:06:49Z | 2026-10-27T08:36:54Z |
| 0329 | 2026-10-27T16:38:41Z | 2026-10-27T16:40:33Z | 2026-10-27T21:31:33Z |
| 0330 | 2026-10-27T16:38:41Z | 2026-10-27T16:44:37Z | 2026-10-27T18:04:29Z |
| 0331 | 2026-10-27T16:38:41Z | 2026-10-27T18:06:14Z | 2026-10-27T20:42:01Z |
| 0332 | 2026-10-30T18:19:35Z | 2026-10-31T14:26:38Z | 2026-10-31T15:31:21Z |
| 0333 | 2026-10-30T18:19:35Z | 2026-10-31T15:32:26Z | 2026-10-31T17:33:54Z |
| 0334 | 2026-10-30T18:19:35Z | 2026-10-31T17:35:20Z | 2026-10-31T20:55:15Z |
| 0335 | 2026-10-30T18:19:35Z | 2026-10-31T18:30:03Z | 2026-10-31T19:56:06Z |
| 0336 | 2026-10-31T17:16:09Z | 2026-10-31T20:59:46Z | 2026-10-31T21:53:24Z |
| 0337 | 2026-11-01T22:57:05Z | 2026-11-01T22:59:10Z | 2026-11-02T00:34:51Z |
| 0338 | 2026-11-02T02:19:13Z | 2026-11-02T02:21:04Z | 2026-11-02T07:51:40Z |
| 0339 | 2026-11-02T02:19:13Z | 2026-11-02T04:13:05Z | 2026-11-02T07:28:12Z |
| 0340 | 2026-11-02T02:19:13Z | 2026-11-02T05:14:21Z | 2026-11-02T06:38:09Z |
| 0341 | 2026-11-02T10:29:15Z | 2026-11-02T10:31:06Z | 2026-11-02T11:41:15Z |
| 0342 | 2026-11-03T19:19:24Z | 2026-11-03T19:21:39Z | 2026-11-04T00:39:14Z |
| 0343 | 2026-11-03T19:19:24Z | 2026-11-03T19:35:01Z | 2026-11-03T21:10:57Z |
| 0344 | 2026-11-03T19:19:24Z | 2026-11-03T21:12:41Z | 2026-11-03T22:45:02Z |
| 0345 | 2026-11-04T10:32:42Z | 2026-11-04T10:34:45Z | 2026-11-04T16:16:45Z |
| 0348 | 2026-11-04T10:32:42Z | 2026-11-04T11:07:51Z | 2026-11-04T12:17:19Z |
| 0347 | 2026-11-04T10:32:42Z | 2026-11-04T12:20:07Z | 2026-11-04T14:15:17Z |
| 0346 | 2026-11-04T10:32:42Z | 2026-11-04T14:16:19Z | 2026-11-04T18:28:59Z |
| 0349 | 2026-11-04T12:34:39Z | 2026-11-04T18:33:54Z | 2026-11-04T20:25:26Z |
| 0350 | 2026-11-04T12:34:39Z | 2026-11-04T20:27:18Z | 2026-11-04T22:59:20Z |
| 0351 | 2026-11-05T01:56:36Z | 2026-11-05T01:58:47Z | 2026-11-05T04:18:52Z |
| 0352 | 2026-11-05T10:56:09Z | 2026-11-05T10:58:03Z | 2026-11-05T15:18:08Z |
| 0353 | 2026-11-05T10:56:09Z | 2026-11-05T13:11:14Z | 2026-11-05T16:05:24Z |
| 0354 | 2026-11-05T13:02:58Z | 2026-11-05T16:36:30Z | 2026-11-06T02:51:03Z |
| 0356 | 2026-11-05T13:02:58Z | 2026-11-05T18:29:00Z | 2026-11-05T20:33:31Z |
| 0355 | 2026-11-05T13:02:58Z | 2026-11-05T20:35:27Z | 2026-11-06T00:03:52Z |
| 0357 | 2026-11-05T13:02:58Z | 2026-11-05T21:14:01Z | 2026-11-06T04:26:53Z |
| 0358 | 2026-11-05T13:02:58Z | 2026-11-06T00:45:39Z | 2026-11-06T03:29:43Z |
| 0359 | 2026-11-06T10:47:45Z | 2026-11-06T10:49:48Z | 2026-11-06T12:40:03Z |
| 0360 | 2026-11-06T10:47:45Z | 2026-11-06T12:41:09Z | 2026-11-06T17:13:32Z |
| 0362 | 2026-11-06T16:05:49Z | 2026-11-06T17:43:50Z | 2026-11-06T21:18:38Z |
| 0363 | 2026-11-06T16:05:49Z | 2026-11-06T18:24:09Z | 2026-11-06T20:53:53Z |
| 0361 | 2026-11-06T16:05:49Z | 2026-11-06T21:20:39Z | 2026-11-06T23:54:37Z |
| 0364 | 2026-11-06T23:46:02Z | 2026-11-07T00:27:18Z | 2026-11-07T02:02:16Z |
| 0377 | 2026-11-11T01:50:36Z | 2026-11-11T17:50:59Z | 2026-11-11T21:01:36Z |
| 0376 | 2026-11-11T01:29:58Z | 2026-11-11T21:43:42Z | 2026-11-11T23:26:07Z |
| 0378 | 2026-11-11T02:35:57Z | 2026-11-11T23:56:01Z | 2026-11-12T00:57:58Z |
| 0379 | 2026-11-11T02:25:48Z | 2026-11-12T01:31:35Z | 2026-11-12T02:46:28Z |
| 0380 | 2026-11-11T19:03:32Z | 2026-11-12T03:26:52Z | 2026-11-12T05:55:35Z |
| 0381 | 2026-11-11T19:03:32Z | 2026-11-12T06:27:51Z | 2026-11-12T08:08:36Z |
| 0382 | 2026-11-11T19:03:32Z | 2026-11-12T08:41:28Z | 2026-11-12T19:27:33Z |
| 0383 | 2026-11-12T21:26:45Z | 2026-11-12T21:29:07Z | 2026-11-13T00:06:29Z |
| 0385 | 2026-11-12T21:26:45Z | 2026-11-13T00:38:42Z | 2026-11-13T02:30:51Z |
| 0384 | 2026-11-12T21:26:45Z | 2026-11-13T02:35:10Z | 2026-11-13T04:07:07Z |
| 0386 | 2026-11-12T22:45:04Z | 2026-11-13T04:50:11Z | 2026-11-13T09:58:51Z |
| 0387 | 2026-11-12T23:39:55Z | 2026-11-13T07:02:47Z | 2026-11-13T08:07:51Z |
| 0388 | 2026-11-13T11:36:37Z | 2026-11-13T11:39:13Z | 2026-11-13T13:02:57Z |
| 0389 | 2026-11-13T11:48:46Z | 2026-11-13T13:40:40Z | 2026-11-13T15:58:55Z |
| 0390 | 2026-11-13T14:46:11Z | 2026-11-14T01:45:02Z | 2026-11-14T03:10:54Z |
| 0391 | 2026-11-14T10:57:12Z | 2026-11-14T12:25:12Z | 2026-11-14T13:57:51Z |
| 0392 | 2026-11-14T13:15:22Z | 2026-11-14T14:35:19Z | 2026-11-14T17:16:03Z |
| 0393 | 2026-11-14T13:30:48Z | 2026-11-14T21:07:33Z | 2026-11-14T23:43:43Z |
| 0396 | 2026-11-14T15:23:18Z | 2026-11-15T11:49:30Z | 2026-11-15T12:56:52Z |
