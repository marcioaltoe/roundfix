package spec

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var lessonNamePattern = regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2})-[a-z0-9]+(?:-[a-z0-9]+)*\.md$`)

// LessonLearnedProblems reports departures from the Baseline's lesson form.
// It reads only the supplied name and bytes; a conforming lesson returns nil.
func LessonLearnedProblems(name string, content []byte) []string {
	var problems []string
	match := lessonNamePattern.FindStringSubmatch(filepath.Base(name))
	if len(match) != 2 || !lessonDate(match[1]) {
		problems = append(problems, "the name is not YYYY-MM-DD-<kebab-slug>.md")
	}
	header, body, err := splitFrontmatter(content)
	var fields struct {
		Created string `yaml:"created"`
	}
	if err != nil || yaml.Unmarshal(header, &fields) != nil || !lessonDate(fields.Created) {
		problems = append(problems, "the header has no created date")
	}
	if err != nil {
		body = content
	}
	lines := strings.Split(strings.ReplaceAll(string(body), "\r\n", "\n"), "\n")
	titles := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") && strings.TrimSpace(strings.TrimPrefix(line, "# ")) != "" {
			titles++
		}
	}
	if titles != 1 {
		problems = append(problems, "the title line is missing")
	}
	previous := -1
	for _, label := range []string{"Context", "Tried", "Went wrong", "Cause", "Solved by", "Applies when", "Does not apply when", "Links"} {
		count, position, empty := 0, -1, false
		for i, line := range lines {
			if rest, ok := strings.CutPrefix(line, "- "+label+":"); ok {
				count++
				if position == -1 {
					position = i
				}
				empty = empty || strings.TrimSpace(rest) == ""
			}
		}
		switch {
		case count == 0:
			problems = append(problems, fmt.Sprintf("the line %q is missing", label+":"))
		default:
			if empty {
				problems = append(problems, fmt.Sprintf("the line %q is empty", label+":"))
			}
			if count > 1 || position < previous {
				problems = append(problems, fmt.Sprintf("the line %q is out of order", label+":"))
			}
			if position > previous {
				previous = position
			}
		}
	}
	return problems
}

func lessonDate(value string) bool {
	_, err := time.Parse(time.DateOnly, value)
	return err == nil
}
