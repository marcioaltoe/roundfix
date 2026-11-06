package delivery

import "strings"

type headCheckSet struct {
	Pending, Failed      []string
	RunFailures, Members int
}

func evaluateHeadCheckSet(report CheckReport, rerunAttempts map[string]int) headCheckSet {
	var set headCheckSet
	checkRuns := make(map[string]bool)
	pendingRuns := make(map[string]bool)
	failedRuns := make(map[string]bool)
	for _, check := range report.Checks {
		set.Members++
		label := checkLabel(check)
		runID := checkRunID(check.Link)
		if runID != "" {
			checkRuns[runID] = true
		}
		switch bucket := strings.ToLower(strings.TrimSpace(check.Bucket)); bucket {
		case "pass", "skipping":
		case "fail", "cancel", "cancelled":
			if _, reran := rerunAttempts[runID]; reran {
				set.Pending = append(set.Pending, label+" (re-run requested)")
				pendingRuns[runID] = true
			} else {
				state := strings.ToLower(strings.TrimSpace(check.State))
				if state == "" {
					state = bucket
				}
				set.Failed = append(set.Failed, failedCheckLabel(label, state, check.Link))
				failedRuns[runID] = true
			}
		default:
			set.Pending = append(set.Pending, label)
			pendingRuns[runID] = true
		}
	}
	runPaths := make(map[string]bool)
	for _, run := range report.WorkflowRuns {
		runPaths[run.Path] = true
		if !checkRuns[run.ID] {
			set.Members++
		}
		if pendingRuns[run.ID] {
			continue
		}
		name := workflowRunLabel(run)
		if run.Status != "completed" {
			set.Pending = append(set.Pending, name+" ("+run.Status+")")
		} else if !workflowRunPassed(run) {
			if attempt, reran := rerunAttempts[run.ID]; reran && attempt >= run.Attempt {
				set.Pending = append(set.Pending, name+" (re-run requested)")
			} else if !failedRuns[run.ID] {
				set.RunFailures++
				set.Failed = append(set.Failed, failedCheckLabel(name, run.Conclusion, run.URL))
			}
		}
	}
	for _, workflow := range report.ExpectedWorkflows {
		if !runPaths[workflow] {
			set.Members++
			set.Pending = append(set.Pending, workflow+" (no run yet)")
		}
	}
	if set.Members == 0 {
		set.Pending = []string{"no check reported"}
	}
	return set
}

func checksFailedBlocker(report CheckReport) string {
	var failed []string
	failedRuns := make(map[string]bool)
	for _, check := range report.Checks {
		bucket := strings.ToLower(strings.TrimSpace(check.Bucket))
		if bucket != "fail" && bucket != "cancel" && bucket != "cancelled" {
			continue
		}
		state := strings.ToLower(strings.TrimSpace(check.State))
		if state == "" {
			state = bucket
		}
		failed = append(failed, failedCheckLabel(checkLabel(check), state, check.Link))
		if id := checkRunID(check.Link); id != "" {
			failedRuns[id] = true
		}
	}
	for _, run := range report.WorkflowRuns {
		if run.Status == "completed" && !workflowRunPassed(run) && !failedRuns[run.ID] {
			failed = append(failed, failedCheckLabel(workflowRunLabel(run), run.Conclusion, run.URL))
		}
	}
	if len(failed) == 0 {
		return BlockerChecksFailed
	}
	return BlockerChecksFailed + ": " + strings.Join(failed, "; ")
}

func checkLabel(check PullRequestCheck) string {
	if check.Workflow == "" {
		return check.Name
	}
	return check.Workflow + " / " + check.Name
}

func workflowRunLabel(run WorkflowRun) string {
	if run.Name == "" {
		return run.Path
	}
	return run.Name
}

func workflowRunPassed(run WorkflowRun) bool {
	return run.Conclusion == "success" || run.Conclusion == "neutral" || run.Conclusion == "skipped"
}

func failedCheckLabel(label, conclusion, link string) string {
	label += " (" + conclusion + ")"
	if link != "" {
		label += " " + link
	}
	return label
}
