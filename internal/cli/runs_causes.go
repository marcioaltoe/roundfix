package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
)

const runsCausesSynopsis = "roundfix runs causes [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]"

func runRunsCausesCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprintf(stdout, "Usage:\n  %s\n\nRead-only causes of failed Verification attempts and corrective Tasks.\nOptions:\n  --since   Inclusive UTC date\n  --until   Exclusive UTC date\n  --format  text (default) or json\n\nExits 0 for results, 2 for invalid usage, 1 for an unreadable Run Database.\n", runsCausesSynopsis)
		return exitOK
	}
	fail := func(err error, code int) int {
		fmt.Fprintf(stderr, "%s: runs causes failed: %v\nRun 'roundfix runs causes --help' for usage.\n", app.Name, err)
		return code
	}
	fs := flag.NewFlagSet("runs causes", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var since, until string
	var format string
	fs.StringVar(&since, "since", "", "Inclusive UTC date")
	fs.StringVar(&until, "until", "", "Exclusive UTC date")
	fs.StringVar(&format, "format", "text", "text or json")
	if err := fs.Parse(args); err != nil {
		return fail(err, exitPreflight)
	}
	if len(fs.Args()) > 0 {
		return fail(fmt.Errorf("unexpected argument %q", fs.Args()[0]), exitPreflight)
	}
	req := runcause.Request{}
	for _, d := range []struct {
		name, value string
		target      *time.Time
	}{{"since", since, &req.Since}, {"until", until, &req.Until}} {
		if d.value == "" {
			explicit := false
			fs.Visit(func(f *flag.Flag) {
				if f.Name == d.name {
					explicit = true
				}
			})
			if !explicit {
				continue
			}
		}
		parsed, err := time.Parse(time.DateOnly, d.value)
		if err != nil {
			return fail(fmt.Errorf("--%s must be a date in YYYY-MM-DD form, got %q", d.name, d.value), exitPreflight)
		}
		*d.target = parsed
	}
	if !req.Since.IsZero() && !req.Until.IsZero() && !req.Until.After(req.Since) {
		return fail(fmt.Errorf("--until must be after --since"), exitPreflight)
	}
	if format != "text" && format != "json" {
		return fail(fmt.Errorf("unknown --format %q; use text or json", format), exitPreflight)
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return fail(err, exitPreflight)
	}
	if strings.TrimSpace(loaded.GitRoot) == "" {
		return fail(fmt.Errorf("runs causes requires a Git repository"), exitPreflight)
	}
	req.RepositoryRoot, err = roundconfig.RepositoryRoot(loaded.GitRoot)
	if err != nil {
		return fail(err, exitPreflight)
	}
	// Historical inspection permits a missing Spec Root: archived graphs can
	// still exist, and otherwise Build reports specs_not_found.
	path := filepath.Clean(loaded.Config.Specs.Root)
	if !filepath.IsAbs(path) {
		path = filepath.Join(loaded.GitRoot, path)
	}
	root := roundconfig.SpecsRoot{Path: path, BuiltInRoot: path == filepath.Join(loaded.GitRoot, "docs", "specs")}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		root, err = roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
		if err != nil {
			return fail(err, exitPreflight)
		}
	}
	req.SpecsRoot = root.Path
	req.ArchiveRoot = spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot)
	reader, found, err := openRunsListReader(ctx, loaded.HomeDir)
	if err != nil {
		return fail(err, 1)
	}
	if found {
		req.Reader = reader
		defer func() { _ = reader.Close() }()
	}
	table, err := runcause.Load()
	if err != nil {
		return fail(err, 1)
	}
	report, err := runcause.Build(ctx, table, req)
	if err != nil {
		return fail(err, 1)
	}
	if format == "json" {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			return fail(err, 1)
		}
	} else {
		printRunsCauses(stdout, table, report)
	}
	return exitOK
}

func printRunsCauses(out io.Writer, table runcause.Table, report runcause.Report) {
	for _, item := range report.Items {
		attempt, signature := "-", "-"
		if item.Attempt != nil {
			attempt = fmt.Sprintf("attempt-%d", *item.Attempt)
		}
		if item.Signature != nil {
			signature = *item.Signature
		}
		check := strings.Join(strings.Fields(strings.SplitN(item.Check, "\n", 2)[0]), " ")
		chars := []rune(check)
		if len(chars) > 80 {
			check = string(chars[:80]) + "…"
		}
		fmt.Fprintf(out, "%s %s/%s %s %s %s %s: %s\n", item.Kind, item.Spec, item.Task, item.RunID, attempt, item.Class, signature, check)
	}
	since, until := "*", "*"
	if report.Window.Since != nil {
		since = *report.Window.Since
	}
	if report.Window.Until != nil {
		until = *report.Window.Until
	}
	classes := []string{}
	for _, class := range append(append([]string{}, table.Classes...), "unclassified") {
		classes = append(classes, fmt.Sprintf("%s %d", class, report.Summary[class]))
	}
	fmt.Fprintf(out, "Causes: %d item(s) from %d terminal Spec Run(s) in window %s..%s; %s; repository knowledge %d of %d classified; signatures %s\n", report.Summary["items"], report.Runs, since, until, strings.Join(classes, ", "), report.Summary["repository_knowledge"], report.Summary["classified"], report.SignaturesSHA256[:12])
}
