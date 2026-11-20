package spec

import (
	"reflect"
	"strings"
	"testing"
)

func TestLessonLearnedProblems(t *testing.T) {
	t.Parallel()
	valid := "---\ncreated: 2026-11-19\n---\n# A lesson\n\n- Context: work\n- Tried: attempt\n- Went wrong: symptom\n- Cause: proof\n- Solved by: solution\n- Applies when: here\n- Does not apply when: elsewhere\n- Links: evidence\n"
	for _, tc := range []struct {
		name, file, content string
		want                []string
	}{
		{"valid", "2026-11-19-a-lesson.md", valid, nil},
		{"name_not_dated", "a-lesson.md", valid, []string{"the name is not YYYY-MM-DD-<kebab-slug>.md"}},
		{"name_bad_date", "2026-02-30-a-lesson.md", valid, []string{"the name is not YYYY-MM-DD-<kebab-slug>.md"}},
		{"no_created", "2026-11-19-a-lesson.md", strings.Replace(valid, "created: 2026-11-19", "other: value", 1), []string{"the header has no created date"}},
		{"created_bad_date", "2026-11-19-a-lesson.md", strings.Replace(valid, "2026-11-19", "2026-02-30", 1), []string{"the header has no created date"}},
		{"no_title", "2026-11-19-a-lesson.md", strings.Replace(valid, "# A lesson\n", "", 1), []string{"the title line is missing"}},
		{"label_missing", "2026-11-19-a-lesson.md", strings.Replace(valid, "- Cause: proof\n", "", 1), []string{`the line "Cause:" is missing`}},
		{"label_empty", "2026-11-19-a-lesson.md", strings.Replace(valid, "- Cause: proof", "- Cause:   ", 1), []string{`the line "Cause:" is empty`}},
		{"label_out_of_order", "2026-11-19-a-lesson.md", strings.Replace(valid, "- Tried: attempt\n- Went wrong: symptom", "- Went wrong: symptom\n- Tried: attempt", 1), []string{`the line "Went wrong:" is out of order`}},
		{"label_duplicate", "2026-11-19-a-lesson.md", valid + "- Links: duplicate\n", []string{`the line "Links:" is out of order`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := LessonLearnedProblems(tc.file, []byte(tc.content)); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("problems=%q want=%q", got, tc.want)
			}
		})
	}
}
