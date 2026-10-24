package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/releaseplan"
	"roundfix/internal/spec"
)

type followUpPendingSpec struct {
	Spec           string  `json:"spec"`
	QueuePosition  *int    `json:"queue_position"`
	Stage          *string `json:"stage"`
	TasksCompleted int     `json:"tasks_completed"`
	TasksTotal     int     `json:"tasks_total"`
}

type followUpReleaseChange struct {
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
}

type followUpRelease struct {
	State           string                  `json:"state"`
	Base            string                  `json:"base"`
	ProposedVersion string                  `json:"proposed_version"`
	Changes         []followUpReleaseChange `json:"changes"`
	Reason          *string                 `json:"reason"`
}

type followUpBlocker struct {
	Spec      string `json:"spec"`
	Blocker   string `json:"blocker"`
	ParkClass string `json:"park_class"`
	Next      string `json:"next"`
}

type followUpMaintainer struct {
	Kind   string `json:"kind"`
	Spec   string `json:"spec,omitempty"`
	Detail string `json:"detail"`
	Answer string `json:"answer,omitempty"`
}

type followUpStatusDocument struct {
	Schema       string                `json:"schema"`
	View         string                `json:"view"`
	State        string                `json:"state"`
	Repository   followUpRepository    `json:"repository"`
	Sources      []followUpSource      `json:"sources"`
	InProgress   followUpProgress      `json:"in_progress"`
	PendingSpecs []followUpPendingSpec `json:"pending_specs"`
	Release      followUpRelease       `json:"release"`
	Blockers     []followUpBlocker     `json:"blockers"`
	Maintainer   []followUpMaintainer  `json:"maintainer"`
}

func runFollowUpStatus(ctx context.Context, req followUpRequest, stdout, stderr io.Writer, loaded roundconfig.Loaded, environment commandEnvironment) int {
	progress, sources, err := readFollowUpProgress(ctx, loaded, environment)
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: follow-up failed: %v\nRun 'roundfix follow-up --help' for usage.\n", err)
		return exitRunFailed
	}
	doc := followUpStatusDocument{
		Schema: "roundfix/follow-up/v1", View: "status", State: "work_found",
		Repository: followUpRepository{loaded.GitRoot, followUpNamespace(ctx, loaded)},
		InProgress: progress, Sources: sources,
		Blockers: []followUpBlocker{}, Maintainer: []followUpMaintainer{},
	}
	pending, source := followUpPendingSpecs(loaded, progress)
	doc.PendingSpecs = pending
	doc.Sources = append(doc.Sources, source)
	for _, item := range progress.DeliveryQueue {
		if item.Stage == "parked" {
			doc.Blockers = append(doc.Blockers, followUpBlocker{item.Spec, item.Blocker, item.ParkClass, item.Next})
		}
	}
	if question, found := delivery.PendingQuestionFor(progress.queue); found {
		doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "pending_question", Spec: question.SpecSlug, Detail: question.Blocker, Answer: question.Answer})
	}
	plan, planErr := releaseplan.Build(ctx, releaseplan.Request{}, newReleasePlanGitSource(environment.workDir, environment.dependencies.releasePlanGitRunner))
	doc.Release = followUpRelease{Changes: []followUpReleaseChange{}}
	releaseSource := followUpSourceResult("release-plan", 1, planErr)
	switch {
	case errors.Is(planErr, releaseplan.ErrInvalidReleaseRange):
		doc.Release.State = "no_unreleased_changes"
		releaseSource = followUpSourceResult("release-plan", 0, nil)
	case planErr != nil:
		doc.Release.State = "unavailable"
		doc.Release.Reason = followUpString(planErr.Error())
	default:
		doc.Release.State = string(plan.State)
		doc.Release.Base = plan.Base.Tag
		doc.Release.ProposedVersion = plan.ProposedVersion
		for _, change := range plan.Changes {
			doc.Release.Changes = append(doc.Release.Changes, followUpReleaseChange{change.CommitSHA, change.Subject})
		}
		switch plan.State {
		case releaseplan.StateApprovalRequired:
			doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "release_approval", Detail: plan.Approval.Question})
		case releaseplan.StateManualClassificationRequired:
			doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "release_classification", Detail: "manual classification required"})
		}
	}
	doc.Sources = append(doc.Sources, releaseSource)
	if progress.count() == 0 && len(doc.PendingSpecs) == 0 && len(doc.Blockers) == 0 && len(doc.Maintainer) == 0 && (doc.Release.State == "no_unreleased_changes" || doc.Release.State == string(releaseplan.StateNoRelease)) {
		doc.State = "nothing_to_do"
	}
	if req.format == "json" {
		err = json.NewEncoder(stdout).Encode(doc)
	} else {
		_, err = fmt.Fprint(stdout, renderFollowUpStatusText(doc))
	}
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: follow-up failed: write status report: %v\nRun 'roundfix follow-up --help' for usage.\n", err)
		return exitRunFailed
	}
	return exitOK
}

func followUpPendingSpecs(loaded roundconfig.Loaded, progress followUpProgress) ([]followUpPendingSpec, followUpSource) {
	rows := []followUpPendingSpec{}
	root, active, err := readFollowUpSpecs(loaded)
	queued := make(map[string]bool)
	for _, item := range progress.queue.Items {
		queued[item.SpecSlug] = true
	}
	queue := append([]followUpQueueItem{}, progress.DeliveryQueue...)
	sort.Slice(queue, func(i, j int) bool { return queue[i].Position < queue[j].Position })
	for _, item := range queue {
		position, stage := item.Position, item.Stage
		rows = append(rows, followUpPendingSpec{Spec: item.Spec, QueuePosition: &position, Stage: &stage})
	}
	for _, item := range active {
		if !queued[item.Slug] {
			rows = append(rows, followUpPendingSpec{Spec: item.Slug})
		}
	}
	if err == nil {
		for index := range rows {
			graph, loadErr := spec.Load(root, rows[index].Spec)
			if loadErr != nil {
				err = errors.Join(err, loadErr)
				continue
			}
			rows[index].TasksTotal = len(graph.Tasks)
			for _, task := range graph.Tasks {
				if task.Status == spec.StatusCompleted {
					rows[index].TasksCompleted++
				}
			}
		}
	}
	return rows, followUpSourceResult("specs", len(rows), err)
}

func renderFollowUpStatusText(doc followUpStatusDocument) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Status for %s in %s\n", doc.Repository.Namespace, doc.Repository.Root)
	fmt.Fprintf(&out, "In progress: %d\n", doc.InProgress.count())
	for _, run := range doc.InProgress.Runs {
		fmt.Fprintf(&out, "  run %s %s %s\n", run.ID, run.State, followUpDash(run.Spec))
	}
	for _, item := range doc.InProgress.DeliveryQueue {
		fmt.Fprintf(&out, "  queue %d %s %s %s\n", item.Position, item.Spec, item.Stage, followUpDash(item.Blocker))
	}
	for _, pr := range doc.InProgress.PullRequests {
		fmt.Fprintf(&out, "  pull request #%d %s: %s\n", pr.Number, pr.Branch, pr.Title)
	}
	git := doc.InProgress.Git
	if git.hasChanges() {
		fmt.Fprintf(&out, "  git %s %s: %d changed path(s), %d unpushed commit(s)\n", git.Branch, git.Head, len(git.ChangedPaths), git.UnpushedCommits)
	}
	fmt.Fprintf(&out, "Pending Specs: %d\n", len(doc.PendingSpecs))
	for _, item := range doc.PendingSpecs {
		position, stage := "-", "not queued"
		if item.QueuePosition != nil {
			position = strconv.Itoa(*item.QueuePosition)
		}
		if item.Stage != nil {
			stage = *item.Stage
		}
		fmt.Fprintf(&out, "  %s %s %s %d/%d Tasks\n", position, item.Spec, stage, item.TasksCompleted, item.TasksTotal)
	}
	switch doc.Release.State {
	case "no_unreleased_changes":
		out.WriteString("Next release: no unreleased changes\n")
	case string(releaseplan.StateNoRelease):
		out.WriteString("Next release: none (no_release)\n")
	case "unavailable":
		fmt.Fprintf(&out, "Next release: unavailable: %s\n", *doc.Release.Reason)
	default:
		fmt.Fprintf(&out, "Next release: %s from %s (%s)\n", doc.Release.ProposedVersion, doc.Release.Base, doc.Release.State)
	}
	for _, change := range doc.Release.Changes {
		fmt.Fprintf(&out, "  %s %s\n", shortReleasePlanSHA(change.Commit), change.Subject)
	}
	fmt.Fprintf(&out, "Blockers: %d\n", len(doc.Blockers))
	for _, blocker := range doc.Blockers {
		fmt.Fprintf(&out, "  %s %s: %s\n", blocker.Spec, blocker.ParkClass, blocker.Next)
	}
	fmt.Fprintf(&out, "Waiting on the maintainer: %d\n", len(doc.Maintainer))
	for _, item := range doc.Maintainer {
		if item.Kind == "pending_question" {
			fmt.Fprintf(&out, "  question %s parked %s: %s\n", item.Spec, item.Detail, item.Answer)
		} else {
			fmt.Fprintf(&out, "  release %s: %s\n", doc.Release.ProposedVersion, item.Detail)
		}
	}
	out.WriteString("Sources:\n")
	for _, source := range doc.Sources {
		fmt.Fprintf(&out, "  %s: %s", source.Name, source.State)
		if source.State == "checked" {
			fmt.Fprintf(&out, " %d", source.Count)
		}
		if source.Detail != "" {
			fmt.Fprintf(&out, ": %s", source.Detail)
		}
		out.WriteByte('\n')
	}
	return out.String()
}
