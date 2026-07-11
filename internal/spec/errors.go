package spec

import (
	"fmt"
	"strings"
)

// SpecNotFoundError reports that the requested Spec has no directory with a
// _prd.md under docs/specs/.
type SpecNotFoundError struct {
	Slug string
	Dir  string
}

func (err SpecNotFoundError) Error() string {
	return fmt.Sprintf("Spec %q not found: no _prd.md under %q; check the slug or run the spec workflow first", err.Slug, err.Dir)
}

// InactiveSpecError reports that the Spec exists but its PRD frontmatter
// status is not "active".
type InactiveSpecError struct {
	Slug   string
	Status string
}

func (err InactiveSpecError) Error() string {
	return fmt.Sprintf("Spec %q is not active: _prd.md frontmatter status is %q; expected %q", err.Slug, err.Status, prdStatusActive)
}

// ManifestError reports a Task Graph manifest that is missing or fails to
// parse or validate, naming the failed check.
type ManifestError struct {
	Path   string
	Reason string
	Err    error
}

func (err ManifestError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("Task Graph manifest %q: %s: %v", err.Path, err.Reason, err.Err)
	}
	return fmt.Sprintf("Task Graph manifest %q: %s", err.Path, err.Reason)
}

func (err ManifestError) Unwrap() error {
	return err.Err
}

// ManifestSchemaError reports a Task Graph manifest whose schema is not the
// supported spec-tasks version.
type ManifestSchemaError struct {
	Path   string
	Schema string
}

func (err ManifestSchemaError) Error() string {
	return fmt.Sprintf("Task Graph manifest %q has schema %q; expected %q", err.Path, err.Schema, manifestSchema)
}

// UnknownNeedError reports a Task whose needs entry names no Task Graph node.
type UnknownNeedError struct {
	TaskID string
	Need   string
}

func (err UnknownNeedError) Error() string {
	return fmt.Sprintf("Task %q needs %q, which is not a Task Graph node id", err.TaskID, err.Need)
}

// CycleError reports a dependency cycle in the Task Graph, listing the Tasks
// that could not be ordered, in manifest order.
type CycleError struct {
	TaskIDs []string
}

func (err CycleError) Error() string {
	return fmt.Sprintf("Task Graph has a dependency cycle among: %s", strings.Join(err.TaskIDs, ", "))
}

// MissingTaskFileError reports a Task Graph node whose task file does not
// exist.
type MissingTaskFileError struct {
	TaskID string
	Path   string
}

func (err MissingTaskFileError) Error() string {
	return fmt.Sprintf("Task %q file %q does not exist", err.TaskID, err.Path)
}

// TaskFileError reports a task file that exists but cannot be parsed.
type TaskFileError struct {
	TaskID string
	Path   string
	Err    error
}

func (err TaskFileError) Error() string {
	return fmt.Sprintf("Task %q file %q: %v", err.TaskID, err.Path, err.Err)
}

func (err TaskFileError) Unwrap() error {
	return err.Err
}

// TaskContextError reports an invalid Task-authored ## Context entry.
type TaskContextError struct {
	Kind   string
	Path   string
	Reason string
}

func (err TaskContextError) Error() string {
	parts := []string{"Task Context entry"}
	if err.Kind != "" {
		parts = append(parts, fmt.Sprintf("kind %q", err.Kind))
	}
	if err.Path != "" {
		parts = append(parts, fmt.Sprintf("path %q", err.Path))
	}
	if err.Reason != "" {
		parts = append(parts, err.Reason)
	}
	return strings.Join(parts, ": ")
}

// MissingVerificationError reports a task file with no parseable Verification
// command.
type MissingVerificationError struct {
	TaskID string
	Path   string
}

func (err MissingVerificationError) Error() string {
	return fmt.Sprintf("Task %q has no Verification commands: add a backticked command bullet to the ## Verification section in %q", err.TaskID, err.Path)
}

// QAReportError reports a QA Report that exists but whose verdict cannot be
// read.
type QAReportError struct {
	Path string
	Err  error
}

func (err QAReportError) Error() string {
	return fmt.Sprintf("QA Report %q is unreadable: %v", err.Path, err.Err)
}

func (err QAReportError) Unwrap() error {
	return err.Err
}
