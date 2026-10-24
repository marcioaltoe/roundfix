package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func followUpNamespace(ctx context.Context, loaded roundconfig.Loaded) string {
	if ns := loaded.Config.Secondbrain.InboxNamespace; ns != "" {
		return ns
	}
	root, err := roundconfig.RepositoryRoot(loaded.GitRoot)
	if err != nil {
		root = loaded.GitRoot
	}
	raw, err := (preflight.ExecGitRunner{}).RunGit(ctx, root, "remote", "get-url", "origin")
	if err == nil {
		if _, repo, ok := parseReadinessRemote(raw); ok {
			return filepath.Base(repo)
		}
	}
	return filepath.Base(root)
}

func followUpSourceResult(name string, count int, err error) followUpSource {
	source := followUpSource{Name: name, State: "checked", Count: count}
	if err != nil {
		source.State = "unavailable"
		source.Count = 0
		source.Detail = followUpFirstLine(err.Error())
	}
	return source
}

func followUpFirstLine(value string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(value), "\n")
	return strings.TrimSpace(line)
}

type followUpMarkdown struct {
	Type        string `yaml:"type"`
	TypeHint    string `yaml:"type-hint"`
	Status      string `yaml:"status"`
	Created     string `yaml:"created"`
	CreatedAt   string `yaml:"created_at"`
	Origin      string `yaml:"origin"`
	title, text string
}

func readFollowUpMarkdown(path string) (followUpMarkdown, error) {
	var fields followUpMarkdown
	content, err := os.ReadFile(path)
	if err != nil {
		return fields, fmt.Errorf("read %s: %w", path, err)
	}
	body := strings.ReplaceAll(string(content), "\r\n", "\n")
	if strings.HasPrefix(body, "---\n") {
		front, rest, closed := strings.Cut(strings.TrimPrefix(body, "---\n"), "\n---")
		if !closed {
			return fields, fmt.Errorf("read %s: unterminated frontmatter", path)
		}
		if err := yaml.Unmarshal([]byte(front), &fields); err != nil {
			return fields, fmt.Errorf("read %s frontmatter: %w", path, err)
		}
		body = strings.TrimPrefix(rest, "\n")
	}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "# ") {
			fields.title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}
	if fields.title == "" {
		fields.title = strings.TrimSuffix(filepath.Base(path), ".md")
	}
	fields.text = strings.TrimSpace(body)
	return fields, nil
}

func followUpFilenameDate(path string) string {
	name := filepath.Base(path)
	if len(name) >= 10 && !followUpDate(name[:10]).IsZero() {
		return name[:10]
	}
	return ""
}

func followUpInbox(root, namespace string) ([]followUpItem, followUpSource) {
	items := []followUpItem{}
	source := followUpSourceResult("secondbrain-inbox", 0, nil)
	if root == "" {
		source.State = "missing"
		source.Detail = "secondbrain.root is not set in User Config"
		return items, source
	}
	inbox := filepath.Join(root, "inbox")
	info, err := os.Stat(inbox)
	if err == nil && !info.IsDir() {
		err = fmt.Errorf("%s is not a directory", inbox)
	}
	if err != nil {
		return items, followUpSourceResult(source.Name, 0, err)
	}
	dir := filepath.Join(inbox, namespace)
	// Do not traverse a symlink standing in for the namespace either.
	info, err = os.Lstat(dir)
	if os.IsNotExist(err) {
		return items, source
	}
	if err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		err = fmt.Errorf("%s is not a regular directory", dir)
	}
	if err != nil {
		return items, followUpSourceResult(source.Name, 0, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return items, followUpSourceResult(source.Name, 0, err)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		fields, err := readFollowUpMarkdown(filepath.Join(dir, entry.Name()))
		if err != nil {
			return []followUpItem{}, followUpSourceResult(source.Name, 0, err)
		}
		if fields.Status != "" && fields.Status != "pending" {
			continue
		}
		date := fields.CreatedAt
		if date == "" {
			date = followUpFilenameDate(entry.Name())
		}
		item := followUpItem{Kind: "inbox", Ref: filepath.ToSlash(filepath.Join("inbox", namespace, entry.Name())), Type: fields.TypeHint, Date: date, Title: fields.title, text: fields.text, origin: fields.Origin}
		if fields.Origin != namespace {
			item.RequestedBy = followUpString(fields.Origin)
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		a, b := followUpDate(items[i].Date), followUpDate(items[j].Date)
		if a.IsZero() != b.IsZero() {
			return !a.IsZero()
		}
		if !a.Equal(b) {
			return a.Before(b)
		}
		return items[i].Ref < items[j].Ref
	})
	source.Count = len(items)
	return items, source
}

func followUpMessages(path string) ([]followUpItem, followUpSource, error) {
	items := []followUpItem{}
	source := followUpSourceResult("ai-memory-messages", 0, nil)
	if path == "" {
		source.State = "missing"
		source.Detail = "no --messages file"
		return items, source, nil
	}
	fail := func(err error) ([]followUpItem, followUpSource, error) {
		return nil, source, fmt.Errorf("messages file %q: %w", path, err)
	}
	file, err := os.Open(path)
	if err != nil {
		return fail(err)
	}
	defer func() { _ = file.Close() }()
	const limit = 1 << 20
	content, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return fail(err)
	}
	if len(content) > limit {
		return fail(fmt.Errorf("exceeds 1 MiB"))
	}
	var document struct {
		Schema   string `json:"schema"`
		Messages []struct {
			ID        string `json:"id"`
			From      string `json:"from"`
			CreatedAt string `json:"created_at"`
			Subject   string `json:"subject"`
			Text      string `json:"text"`
			Type      string `json:"type"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(content, &document); err != nil {
		return fail(err)
	}
	if document.Schema != "roundfix/follow-up-messages/v1" || document.Messages == nil {
		return fail(fmt.Errorf("expected roundfix/follow-up-messages/v1 with a messages array"))
	}
	ids := make(map[string]bool)
	for _, message := range document.Messages {
		if strings.TrimSpace(message.ID) == "" || strings.TrimSpace(message.Text) == "" || followUpDate(message.CreatedAt).IsZero() {
			return fail(fmt.Errorf("each message requires id, created_at (RFC 3339 or YYYY-MM-DD) and text"))
		}
		if ids[message.ID] {
			return fail(fmt.Errorf("duplicate id %q", message.ID))
		}
		ids[message.ID] = true
		typ := message.Type
		switch typ {
		case "":
			typ = "message"
		case "feat", "fix", "perf", "refactor":
		default:
			return fail(fmt.Errorf("invalid type %q", typ))
		}
		title := followUpFirstLine(message.Subject)
		if title == "" {
			title = followUpFirstLine(message.Text)
			runes := []rune(title)
			if len(runes) > 80 {
				title = string(runes[:80])
			}
		}
		items = append(items, followUpItem{Kind: "message", Ref: "message:" + message.ID, Type: typ, Date: message.CreatedAt, Title: title, RequestedBy: followUpString(message.From), text: message.Text, origin: message.From})
	}
	source.Count = len(items)
	return items, source, nil
}

func followUpInventory(loaded roundconfig.Loaded, namespace, messages, workDir string) ([]followUpItem, []followUpSource, error) {
	items, inbox := followUpInbox(loaded.Config.Secondbrain.Root, namespace)
	messagePath := messages
	if messagePath != "" && !filepath.IsAbs(messagePath) {
		messagePath = filepath.Join(workDir, messagePath)
	}
	messageItems, messageSource, err := followUpMessages(messagePath)
	if err != nil {
		return nil, nil, err
	}
	items = append(items, messageItems...)
	sources := []followUpSource{inbox, messageSource}
	intents, intentErr := deliveryPlanIntent(loaded.GitRoot)
	repositoryItems := map[string][]followUpItem{"backlog": {}, "findings": {}, "repository-inbox": {}}
	if intentErr == nil {
		for _, intent := range intents {
			source, kind, typ, date := "", "", "", ""
			switch {
			case intent.Kind == "backlog" && intent.Status == "open":
				source, kind = "backlog", "backlog"
			case intent.Kind == "finding" && (intent.Status == "pending" || intent.Status == "partial"):
				source, kind, typ = "findings", "finding", "finding"
			case intent.Kind == "inbox":
				source, kind, typ = "repository-inbox", "note", "note"
			default:
				continue
			}
			fields, readErr := readFollowUpMarkdown(filepath.Join(loaded.GitRoot, filepath.FromSlash(intent.Path)))
			if readErr != nil {
				intentErr = readErr
				break
			}
			switch kind {
			case "backlog":
				typ, date = fields.Type, fields.Created
			case "finding":
				date = fields.CreatedAt
			case "note":
				date = followUpFilenameDate(intent.Path)
			}
			repositoryItems[source] = append(repositoryItems[source], followUpItem{Kind: kind, Ref: intent.Path, Type: typ, Date: date, Title: fields.title, text: fields.text, origin: namespace})
		}
	}
	for _, name := range []string{"backlog", "findings", "repository-inbox"} {
		found := repositoryItems[name]
		sources = append(sources, followUpSourceResult(name, len(found), intentErr))
		if intentErr == nil {
			items = append(items, found...)
		}
	}
	var specItems []followUpItem
	_, active, specErr := readFollowUpSpecs(loaded)
	if specErr == nil {
		for _, activeSpec := range active {
			fields, err := readFollowUpMarkdown(filepath.Join(activeSpec.Dir, "_prd.md"))
			if err != nil {
				specErr = err
				break
			}
			text := fields.text
			// The judge receives only the PRD's first paragraph, not its full body.
			lines := strings.Split(text, "\n")
			for len(lines) > 0 && (strings.TrimSpace(lines[0]) == "" || strings.HasPrefix(lines[0], "#")) {
				lines = lines[1:]
			}
			text, _, _ = strings.Cut(strings.Join(lines, "\n"), "\n\n")
			specItems = append(specItems, followUpItem{Kind: "spec", Ref: deliverPlanRelativePath(loaded.GitRoot, activeSpec.Dir), Type: "spec", Date: fields.Created, Title: fields.title, text: text, origin: namespace})
		}
	}
	sources = append(sources, followUpSourceResult("specs", len(specItems), specErr))
	if specErr == nil {
		items = append(items, specItems...)
	}
	return items, sources, nil
}

func readFollowUpSpecs(loaded roundconfig.Loaded) (string, []spec.Spec, error) {
	root, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return "", nil, err
	}
	active, _, err := spec.ListActiveDetailed(root.Path)
	return root.Path, active, err
}

type followUpRun struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Spec  string `json:"spec"`
}

type followUpQueueItem struct {
	Position  int    `json:"position"`
	Spec      string `json:"spec"`
	Stage     string `json:"stage"`
	Blocker   string `json:"blocker"`
	ParkClass string `json:"park_class,omitempty"`
	Next      string `json:"next,omitempty"`
}

type followUpPullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Branch  string `json:"headRefName"`
	URL     string `json:"url"`
	IsDraft bool   `json:"isDraft"`
}

type followUpChangedPath struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

type followUpGit struct {
	Branch          string                `json:"branch"`
	Head            string                `json:"head"`
	Upstream        string                `json:"upstream"`
	UnpushedCommits int                   `json:"unpushed_commits"`
	ChangedPaths    []followUpChangedPath `json:"changed_paths"`
}

func (git followUpGit) hasChanges() bool { return len(git.ChangedPaths) > 0 || git.UnpushedCommits > 0 }

type followUpProgress struct {
	queue         store.DeliveryQueue
	Runs          []followUpRun         `json:"runs"`
	DeliveryQueue []followUpQueueItem   `json:"delivery_queue"`
	PullRequests  []followUpPullRequest `json:"pull_requests"`
	Git           followUpGit           `json:"git"`
}

func (progress followUpProgress) count() int {
	count := len(progress.Runs) + len(progress.DeliveryQueue) + len(progress.PullRequests)
	if progress.Git.hasChanges() {
		count++
	}
	return count
}

func readFollowUpProgress(ctx context.Context, loaded roundconfig.Loaded, environment commandEnvironment) (progress followUpProgress, resultSources []followUpSource, resultErr error) {
	progress = followUpProgress{Runs: []followUpRun{}, DeliveryQueue: []followUpQueueItem{}, PullRequests: []followUpPullRequest{}, Git: followUpGit{ChangedPaths: []followUpChangedPath{}}}
	sources := []followUpSource{}
	reader, found, err := openFollowUpReader(ctx, loaded.HomeDir)
	if err != nil {
		return progress, nil, err
	}
	if found {
		defer func() {
			if err := reader.Close(); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("close Run Database reader: %w", err))
			}
		}()
		runs, err := reader.ListRuns(ctx, store.ListRunsQuery{GitRoot: loaded.GitRoot, States: store.StatesActive})
		if err != nil {
			return progress, nil, err
		}
		for _, run := range runs {
			progress.Runs = append(progress.Runs, followUpRun{run.ID, run.State, run.SpecSlug})
		}
		queue, _, err := reader.DeliveryQueue(ctx, loaded.GitRoot)
		if err != nil {
			return progress, nil, err
		}
		progress.queue = queue
		for _, item := range queue.Items {
			if item.Stage == store.DeliveryStageMerged {
				continue
			}
			row := followUpQueueItem{Position: item.Position, Spec: item.SpecSlug, Stage: string(item.Stage), Blocker: item.Blocker}
			if item.Stage == store.DeliveryStageParked {
				park := delivery.ClassifyPark(queue, item)
				row.ParkClass = string(park.Class)
				row.Next = park.Next
			}
			progress.DeliveryQueue = append(progress.DeliveryQueue, row)
		}
	}
	sources = append(sources, followUpSourceResult("runs", len(progress.Runs), nil), followUpSourceResult("delivery-queue", len(progress.DeliveryQueue), nil))
	gh := environment.dependencies.followUpGHRunner
	if gh == nil {
		gh = preflight.ExecGHRunner{}
	}
	raw, ghErr := gh.RunGH(ctx, loaded.GitRoot, "pr", "list", "--state", "open", "--limit", "100", "--json", "number,title,headRefName,url,isDraft")
	if ghErr == nil {
		ghErr = json.Unmarshal([]byte(raw), &progress.PullRequests)
	}
	if ghErr != nil || progress.PullRequests == nil {
		progress.PullRequests = []followUpPullRequest{}
	}
	sources = append(sources, followUpSourceResult("pull-requests", len(progress.PullRequests), ghErr))
	git, gitErr := preflight.InspectGit(ctx, environment.workDir, followUpGitRunner{})
	if gitErr == nil {
		progress.Git.Branch = git.Branch
		progress.Git.Head = git.HEAD
		if len(progress.Git.Head) > 7 {
			progress.Git.Head = progress.Git.Head[:7]
		}
		if git.UpstreamRemote != "" {
			progress.Git.Upstream = git.UpstreamRemote + "/" + git.UpstreamBranch
		}
		progress.Git.UnpushedCommits = git.UnpushedCommits
		for _, path := range git.Dirty {
			progress.Git.ChangedPaths = append(progress.Git.ChangedPaths, followUpChangedPath{path.Status, path.Path})
		}
	}
	gitCount := 0
	if progress.Git.hasChanges() {
		gitCount = 1
	}
	sources = append(sources, followUpSourceResult("git", gitCount, gitErr))
	return progress, sources, nil
}

// followUpReader owns a private live-WAL snapshot, when one is needed. SQLite
// may rebuild and coordinate SHM there without touching the operator's files.
type followUpReader struct {
	*store.Store
	temporaryHome string
}

func (reader *followUpReader) Close() error {
	err := reader.Store.Close()
	if reader.temporaryHome != "" {
		err = errors.Join(err, os.RemoveAll(reader.temporaryHome))
	}
	return err
}

func openFollowUpReader(ctx context.Context, home string) (*followUpReader, bool, error) {
	path := store.DatabasePath(home)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("inspect Run Database: %w", err)
	}
	wal, err := os.Stat(path + "-wal")
	if errors.Is(err, os.ErrNotExist) || (err == nil && wal.Size() <= 32) {
		reader, err := store.OpenStorageReader(ctx, home)
		if err != nil {
			return nil, false, err
		}
		return &followUpReader{Store: reader}, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("inspect Run Database WAL: %w", err)
	}
	return snapshotFollowUpReader(ctx, home)
}

func snapshotFollowUpReader(ctx context.Context, home string) (snapshot *followUpReader, found bool, resultErr error) {
	path := store.DatabasePath(home)
	// Open the existing lock without O_CREATE or write access. The shared CLI
	// lock primitive matches the Store's Unix flock / Windows byte-range lock.
	lock, err := os.Open(path + ".lock")
	if err != nil {
		return nil, false, fmt.Errorf("open Run Database snapshot lock: %w", err)
	}
	if err := lockReviewDispositionFile(ctx, lock); err != nil {
		return nil, false, errors.Join(fmt.Errorf("lock Run Database snapshot: %w", err), lock.Close())
	}
	defer func() {
		resultErr = errors.Join(resultErr, unlockReviewDispositionFile(lock), lock.Close())
		if resultErr != nil && snapshot != nil {
			resultErr = errors.Join(resultErr, snapshot.Close())
			snapshot, found = nil, false
		}
	}()
	temporaryHome, err := os.MkdirTemp("", "roundfix-follow-up-")
	if err != nil {
		return nil, false, fmt.Errorf("create Run Database snapshot: %w", err)
	}
	defer func() {
		if resultErr != nil {
			resultErr = errors.Join(resultErr, os.RemoveAll(temporaryHome))
		}
	}()
	target := store.DatabasePath(temporaryHome)
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return nil, false, fmt.Errorf("create Run Database snapshot directory: %w", err)
	}
	// Copy WAL first: a concurrent checkpoint may put these same pages in the
	// main file, but the writer lock prevents any newer transaction. Never copy
	// SHM: SQLite reconstructs that coordination file in the private directory.
	if err := copyFollowUpDatabaseFile(ctx, path+"-wal", target+"-wal", true); err != nil {
		return nil, false, err
	}
	if err := copyFollowUpDatabaseFile(ctx, path, target, false); err != nil {
		return nil, false, err
	}
	reader, found, err := openRunsListReader(ctx, temporaryHome)
	if err != nil {
		return nil, false, fmt.Errorf("open Run Database snapshot: %w", err)
	}
	if !found {
		return nil, false, fmt.Errorf("Run Database snapshot disappeared")
	}
	return &followUpReader{Store: reader, temporaryHome: temporaryHome}, true, nil
}

func copyFollowUpDatabaseFile(ctx context.Context, source, target string, optional bool) (resultErr error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	input, err := os.Open(source)
	if optional && errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read Run Database snapshot source: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, input.Close()) }()
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create Run Database snapshot file: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, output.Close()) }()
	buffer := make([]byte, 64*1024)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		n, readErr := input.Read(buffer)
		if n > 0 {
			if _, err := output.Write(buffer[:n]); err != nil {
				return fmt.Errorf("copy Run Database snapshot: %w", err)
			}
		}
		if errors.Is(readErr, io.EOF) {
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read Run Database snapshot: %w", readErr)
		}
	}
}

// Git status may refresh the index even during a read. Disable that optional
// write while retaining InspectGit's parsing and working-directory semantics.
type followUpGitRunner struct{}

func (followUpGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	return (preflight.ExecGitRunner{}).RunGit(ctx, workDir, append([]string{"--no-optional-locks"}, args...)...)
}
