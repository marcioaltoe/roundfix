package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/config"
	"roundfix/internal/flowpermissions"
)

const permissionsUsage = `Usage:
  roundfix permissions [--harness <claude|codex|opencode>] [--remove] [--apply <digest>] [--format <text|json>]

Shows the Flow Command List in the rule syntax of an Agent harness and writes
it into this repository's permission file for that harness only when it is
given the digest a preview printed. The list holds the roundfix subcommands
the owned skills name and the Git and GitHub reads they name. It is not a
security boundary. With no --harness it previews every harness whose files
this repository holds. Nothing under your home directory is read or written.

Options:
  --harness  The harness: claude, codex or opencode (reported, never written)
  --remove   Preview, and with --apply perform, the removal of what Roundfix wrote
  --apply    Write exactly what the preview with this digest showed
  --format   Output format: text (default) or json

Exit codes:
  0  a preview was printed, or the plan was written, or nothing needed writing
  1  the permission file could not be read as the harness documents it, or not written
  2  invalid input
  3  the preview changed since the digest was printed; nothing was written
`

type permissionsHarnessReport struct {
	Harness    flowpermissions.Harness `json:"harness"`
	Action     flowpermissions.Action  `json:"action"`
	File       string                  `json:"file"`
	FileExists bool                    `json:"fileExists"`
	State      flowpermissions.State   `json:"state"`
	Add        []string                `json:"add"`
	AddGuards  []string                `json:"addGuards"`
	Remove     []string                `json:"remove"`
	Exclusions []permissionsExclusion  `json:"exclusions"`
	Digest     string                  `json:"digest"`
	Next       string                  `json:"next"`
	Written    bool                    `json:"written,omitempty"`
}

type permissionsExclusion struct {
	Command string `json:"command"`
	Reason  string `json:"reason"`
}

func runPermissionsCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, permissionsUsage)
		return exitOK
	}
	f := flag.NewFlagSet("permissions", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	harness := f.String("harness", "", "Agent harness")
	remove := f.Bool("remove", false, "Remove the list")
	digest := f.String("apply", "", "Preview digest")
	format := f.String("format", "text", "Output format")
	invalid := func(err error) int {
		fmt.Fprintf(stderr, "roundfix: permissions: %v\n", err)
		return exitPreflight
	}
	if err := f.Parse(args); err != nil {
		return invalid(err)
	}
	apply, harnessSet := false, false
	f.Visit(func(v *flag.Flag) {
		if v.Name == "apply" {
			apply = true
		}
		if v.Name == "harness" {
			harnessSet = true
		}
	})
	if f.NArg() > 0 {
		return invalid(fmt.Errorf("unexpected argument %q", f.Arg(0)))
	}
	if *format != "text" && *format != "json" {
		return invalid(fmt.Errorf("--format must be text or json"))
	}
	if harnessSet && *harness != "claude" && *harness != "codex" && *harness != "opencode" {
		return invalid(fmt.Errorf("--harness must be claude, codex or opencode"))
	}
	if apply && *harness == "" {
		return invalid(fmt.Errorf("--apply requires --harness"))
	}
	if apply && *harness == "opencode" {
		return invalid(fmt.Errorf("Roundfix writes nothing for OpenCode"))
	}
	workDir, err := environment.resolveWorkDir("permissions")
	if err != nil {
		return invalid(err)
	}
	root := config.FindGitRoot(workDir)
	if root == "" {
		return invalid(fmt.Errorf("not inside a Git repository"))
	}
	if err := ctx.Err(); err != nil {
		fmt.Fprintln(stderr, err)
		return exitRunFailed
	}
	harnesses := flowpermissions.Present(root)
	if *harness != "" {
		harnesses = []flowpermissions.Harness{flowpermissions.Harness(*harness)}
	}
	if len(harnesses) == 0 {
		fmt.Fprintln(stdout, "No harness files in this repository. Name one with --harness.")
		return exitOK
	}
	action := flowpermissions.ActionApply
	if *remove {
		action = flowpermissions.ActionRemove
	}
	plans := make([]flowpermissions.Plan, 0, len(harnesses))
	for _, h := range harnesses {
		var p flowpermissions.Plan
		if apply {
			p, err = flowpermissions.Apply(root, h, action, *digest)
		} else {
			p, err = flowpermissions.Preview(root, h, action)
		}
		if err != nil {
			switch {
			case errors.Is(err, flowpermissions.ErrOutsideRepository):
				fmt.Fprintf(stderr, "%s leaves the repository through a link. Nothing was written.\n", p.File)
			case errors.Is(err, flowpermissions.ErrUnreadable):
				fmt.Fprintf(stderr, "%s cannot be read as %s documents it: %s. Nothing was written.\n", p.File, permissionsHarnessName(h), p.Reason)
			case errors.Is(err, flowpermissions.ErrPreviewChanged), errors.Is(err, flowpermissions.ErrNotWritable):
				fmt.Fprintf(stderr, "The preview changed since digest %s was printed. Nothing was written.\n", *digest)
				fmt.Fprintf(stderr, "Run the preview again: %s\n", permissionsPreviewCommand(p))
				return exitUnverified
			default:
				fmt.Fprintf(stderr, "roundfix: permissions: %v. Nothing was written.\n", err)
			}
			return exitRunFailed
		}
		plans = append(plans, p)
	}
	if *format == "json" {
		report := struct {
			SchemaVersion string                     `json:"schemaVersion"`
			Roundfix      string                     `json:"roundfix"`
			Harnesses     []permissionsHarnessReport `json:"harnesses"`
		}{SchemaVersion: "roundfix/flow-permissions/v1", Roundfix: app.Version}
		for _, p := range plans {
			r := permissionsHarnessReport{
				Harness: p.Harness, Action: p.Action, File: p.File, FileExists: p.FileExists, State: p.State,
				Add: append([]string{}, p.Add...), AddGuards: append([]string{}, p.AddGuards...), Remove: append([]string{}, p.Remove...),
				Digest: p.Digest, Written: apply,
			}
			for _, ex := range flowpermissions.Exclusions() {
				r.Exclusions = append(r.Exclusions, permissionsExclusion{Command: strings.Join(ex.Words, " "), Reason: ex.Reason})
			}
			if p.Digest != "" {
				r.Next = permissionsPreviewCommand(p) + " --apply " + p.Digest
			}
			report.Harnesses = append(report.Harnesses, r)
		}
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			fmt.Fprintf(stderr, "roundfix: permissions: write JSON: %v\n", err)
			return exitRunFailed
		}
		return exitOK
	}
	var output strings.Builder
	for i, p := range plans {
		if i > 0 {
			fmt.Fprintln(&output)
		}
		if apply {
			printPermissionsWrite(&output, p)
		} else {
			printPermissionsPreview(&output, p)
		}
	}
	if _, err := io.WriteString(stdout, output.String()); err != nil {
		fmt.Fprintf(stderr, "roundfix: permissions: write output: %v\n", err)
		return exitRunFailed
	}
	return exitOK
}

func permissionsHarnessName(h flowpermissions.Harness) string {
	switch h {
	case flowpermissions.Claude:
		return "Claude Code"
	case flowpermissions.Codex:
		return "Codex"
	default:
		return "OpenCode"
	}
}

func permissionsPreviewCommand(p flowpermissions.Plan) string {
	command := "roundfix permissions --harness " + string(p.Harness)
	if p.Action == flowpermissions.ActionRemove {
		command += " --remove"
	}
	return command
}

func printPermissionsRules(w io.Writer, heading string, rules []string) {
	fmt.Fprintf(w, "%s (%d):\n", heading, len(rules))
	for _, rule := range rules {
		fmt.Fprintf(w, "  %s\n", rule)
	}
}

func printPermissionsPreview(w io.Writer, p flowpermissions.Plan) {
	fmt.Fprintf(w, "Flow Command List for %s, from Roundfix %s\n", permissionsHarnessName(p.Harness), app.Version)
	if p.Harness == flowpermissions.OpenCode {
		if p.State == flowpermissions.NotNeeded {
			fmt.Fprintln(w, "State: not needed")
			fmt.Fprintln(w, "OpenCode's base policy allows every shell command, and this repository holds no OpenCode configuration that asks about or denies one.")
		} else {
			fmt.Fprintln(w, "State: own rules")
			fmt.Fprintf(w, "This repository's OpenCode configuration (%s) holds rules that ask about or deny shell commands. Compare them with the list:\n", p.Reason)
			for _, entry := range flowpermissions.Entries() {
				fmt.Fprintf(w, "  %s\n", strings.Join(entry.Words, " "))
			}
		}
		fmt.Fprintln(w, "Roundfix writes nothing for OpenCode.")
		return
	}
	fmt.Fprintf(w, "File: %s", p.File)
	if !p.FileExists {
		fmt.Fprint(w, " (absent, it would be created)")
	}
	fmt.Fprintln(w)
	if p.Action == flowpermissions.ActionRemove && p.State != flowpermissions.Unreadable {
		if p.Digest == "" {
			fmt.Fprintln(w, "Nothing to remove: the file holds no entry of the list.")
			return
		}
		fmt.Fprintln(w)
		printPermissionsRules(w, "Would remove", p.Remove)
		if p.Harness == flowpermissions.Codex {
			fmt.Fprintf(w, "Would delete %s.\n", p.File)
		}
		fmt.Fprintln(w, "An entry with the same text that you wrote yourself is removed with them.")
		fmt.Fprintf(w, "\nNothing was written. To remove exactly this: %s --apply %s\n", permissionsPreviewCommand(p), p.Digest)
		return
	}
	switch p.State {
	case flowpermissions.NotApplied:
		fmt.Fprintln(w, "State: not applied")
	case flowpermissions.Behind:
		fmt.Fprintf(w, "State: behind, %d to add, %d to remove\n", len(p.Add)+len(p.AddGuards), len(p.Remove))
	case flowpermissions.Unreadable:
		fmt.Fprintf(w, "State: not readable: %s\n", p.Reason)
	default:
		fmt.Fprintln(w, "State: current")
	}
	fmt.Fprintln(w)
	if p.Harness == flowpermissions.Claude {
		printPermissionsRules(w, "Would add under permissions.allow", p.Add)
		fmt.Fprintln(w)
		printPermissionsRules(w, "Would add under permissions.ask, to keep these options asked about", p.AddGuards)
		fmt.Fprintln(w, "  These ask in every permission mode, a session that skips prompts included.")
	} else {
		printPermissionsRules(w, "Would allow", p.Add)
		fmt.Fprintln(w, "No option of a listed read can be kept asked about in Codex.")
	}
	if len(p.Remove) > 0 {
		fmt.Fprintln(w)
		printPermissionsRules(w, "Would remove, no longer on the list", p.Remove)
	}
	fmt.Fprintln(w, "\nLeft out, and still asked about:")
	printPermissionsExclusions(w)
	fmt.Fprintln(w, "  A commit, a push, a merge and a Pull Request stay asked about, and a delivery's push and merge stay behind the Spec's authorization.")
	fmt.Fprintln(w, "\nWhat this does not do:")
	fmt.Fprintln(w, "  It is not a security boundary: a roundfix command runs what this repository declares.")
	if p.Harness == flowpermissions.Claude {
		fmt.Fprintln(w, "  A listed read accepts any option; only the options listed above stay asked about.")
		fmt.Fprintln(w, "  The rules apply only after you trust this folder in Claude Code.")
	} else {
		fmt.Fprintln(w, "  In Codex an allowed command runs outside the sandbox with nobody asked.")
		fmt.Fprintln(w, "  Codex rules match the start of a command, so a listed read runs with any option: git show --output=<file> would write that file.")
		fmt.Fprintln(w, "  The rules apply only after you trust this project in Codex.")
	}
	switch {
	case p.State == flowpermissions.Unreadable:
		fmt.Fprintln(w, "\nNothing was written.")
	case p.Digest == "":
		fmt.Fprintln(w, "\nNothing to write: the file holds the installed list.")
	default:
		fmt.Fprintf(w, "\nNothing was written. To write exactly this: %s --apply %s\n", permissionsPreviewCommand(p), p.Digest)
	}
}

func printPermissionsExclusions(w io.Writer) {
	// Keep one line per reason, with the applying command first.
	exclusions := flowpermissions.Exclusions()
	var reasons []string
	commands := map[string][]string{}
	for _, ex := range exclusions {
		if strings.Join(ex.Words, " ") == "roundfix permissions" {
			reasons = append(reasons, ex.Reason)
			break
		}
	}
	for _, ex := range exclusions {
		if _, ok := commands[ex.Reason]; !ok && (len(reasons) == 0 || reasons[0] != ex.Reason) {
			reasons = append(reasons, ex.Reason)
		}
		words := ex.Words
		if len(commands[ex.Reason]) > 0 && words[0] == strings.Fields(commands[ex.Reason][0])[0] {
			words = words[1:]
		}
		commands[ex.Reason] = append(commands[ex.Reason], strings.Join(words, " "))
	}
	for _, reason := range reasons {
		fmt.Fprintf(w, "  %s: %s\n", strings.Join(commands[reason], ", "), reason)
	}
}

func printPermissionsWrite(w io.Writer, p flowpermissions.Plan) {
	if p.Action == flowpermissions.ActionRemove {
		if p.Harness == flowpermissions.Codex {
			fmt.Fprintf(w, "Deleted %s.\n", p.File)
		} else {
			fmt.Fprintf(w, "Removed %d entries from %s.\n", len(p.Remove), p.File)
		}
		return
	}
	if p.Harness == flowpermissions.Claude {
		fmt.Fprintf(w, "Wrote %d entries and %d Guard Entries to %s.\n", len(p.Add), len(p.AddGuards), p.File)
		fmt.Fprintln(w, "The rules apply only after you trust this folder in Claude Code.")
	} else {
		fmt.Fprintf(w, "Wrote %s with %d entries.\n", p.File, len(flowpermissions.Entries()))
		fmt.Fprintln(w, "The rules apply only after you trust this project in Codex.")
	}
	fmt.Fprintf(w, "To remove them: roundfix permissions --harness %s --remove\n", p.Harness)
}
