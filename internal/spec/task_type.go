package spec

import (
	"fmt"
	"strings"
)

const (
	TaskTypeBackend  TaskType = "backend"
	TaskTypeFrontend TaskType = "frontend"
	TaskTypeData     TaskType = "data"
	TaskTypeInfra    TaskType = "infra"
	TaskTypeDocs     TaskType = "docs"
	TaskTypeTest     TaskType = "test"
	TaskTypeChore    TaskType = "chore"
	TaskTypeQA       TaskType = "qa"
)

// TaskType is the author-declared routing category stored in task frontmatter.
type TaskType string

var canonicalTaskTypes = []TaskType{
	TaskTypeBackend,
	TaskTypeFrontend,
	TaskTypeData,
	TaskTypeInfra,
	TaskTypeDocs,
	TaskTypeTest,
	TaskTypeChore,
	TaskTypeQA,
}

// TaskTypeError reports an invalid task frontmatter type value.
type TaskTypeError struct {
	Path  string
	Value string
}

func (err TaskTypeError) Error() string {
	return fmt.Sprintf("Task Type %q in task file %q is invalid; allowed values: %s; update the task frontmatter type to one allowed value", err.Value, err.Path, allowedTaskTypeValues())
}

// AllowedTaskTypes returns the closed set of Task Type values in authoring order.
func AllowedTaskTypes() []TaskType {
	return append([]TaskType(nil), canonicalTaskTypes...)
}

// ParseTaskType validates the exact frontmatter Task Type value.
func ParseTaskType(taskPath string, raw string) (TaskType, error) {
	taskType := TaskType(raw)
	if AllowedTaskType(taskType) && raw == strings.TrimSpace(raw) {
		return taskType, nil
	}
	return "", TaskTypeError{Path: taskPath, Value: raw}
}

// AllowedTaskType reports whether taskType is one of the closed routing values.
func AllowedTaskType(taskType TaskType) bool {
	switch taskType {
	case TaskTypeBackend, TaskTypeFrontend, TaskTypeData, TaskTypeInfra, TaskTypeDocs, TaskTypeTest, TaskTypeChore, TaskTypeQA:
		return true
	default:
		return false
	}
}

func allowedTaskTypeValues() string {
	values := make([]string, 0, len(canonicalTaskTypes))
	for _, taskType := range canonicalTaskTypes {
		values = append(values, string(taskType))
	}
	return strings.Join(values, ", ")
}
