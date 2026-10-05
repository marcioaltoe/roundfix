# Outside source observations — 2026-10-05

Observed through the web reader; initial curl-cffi request failed CONNECT proxy 403. These sources are independent primary publications. All support the design; none proves deployment of the new workflow before its own PR.

- [GitHub rerun documentation](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/re-run-workflows-and-jobs), line 43: reruns retain the original event SHA and ref. This supports explicitly constructing a current-base merge.
- [GitHub workflow commands](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands), lines 176–185: notice emits a log message and annotation, with an optional title parameter. Supports the fixed tested-base label.
- [GitHub check-run annotations API](https://docs.github.com/en/rest/checks/runs), lines 244–302: GET check-runs/{check_run_id}/annotations; maximum per_page 100; response includes title and message. Supports the scripted adapter request and decoding. Real job annotation behavior is not deployed by this QA run.
- [actions/checkout README](https://github.com/actions/checkout), lines 210–214 and 307–309: depth zero fetches all branch and tag history. Supports resolving the base remote ref in the checked-out job.
- [XNU kern_exit.c](https://raw.githubusercontent.com/apple-oss-distributions/xnu/main/bsd/kern/kern_exit.c), lines 1944–1947: P_REF_DEAD makes the process unavailable to proc_find before the later zombie-state transition. Supports an exiting interval before a zombie reading.
- [XNU kern_sig.c](https://raw.githubusercontent.com/apple-oss-distributions/xnu/main/bsd/kern/kern_sig.c), lines 1549–1624: killpg1 counts eligible members and returns EPERM in the POSIX path when none is counted. Supports EPERM as compatible with the observed exiting-only group.
- [XNU kern_sysctl.c](https://raw.githubusercontent.com/apple-oss-distributions/xnu/main/bsd/kern/kern_sysctl.c), lines 1083–1086 and 1143–1146: P_LEXIT is exposed as P_WEXIT alongside the process state. Supports the process-table flag used by the fixture reader. Combined with the real Darwin observations, this is evidence for the flag filter, not an assumption from the Spec.

Historical CI run source commands: `gh run view 37306411237 --repo marcioaltoe/roundfix --attempt <1|2|3> --log`; `gh run view 37308972464 --repo marcioaltoe/roundfix --log`; `gh run view 37309659532 --repo marcioaltoe/roundfix --log`. All exit 0 with permitted network access. Their focused checkout observations are in pr391-checkouts.txt.

The explicitly authorized old Daemon log observation is in prior-daemon-failure.txt; no other real Roundfix state was read.
