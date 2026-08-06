package skills

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ReadinessState is the outcome of comparing a skill's declared compatibility
// identity with Roundfix's independently declared minimum.
type ReadinessState string

const (
	ReadinessSatisfies   ReadinessState = "satisfies"
	ReadinessBelow       ReadinessState = "below"
	ReadinessUnversioned ReadinessState = "unversioned"
)

// ErrSkillVersionUnresolvable distinguishes a source that cannot provide a
// comparable declaration from a skill that is absent.
var ErrSkillVersionUnresolvable = errors.New("skill version unresolvable")

// SkillVersion is a skill's declared compatibility identity. Roundfix reads
// it and never invents one.
type SkillVersion struct {
	Declared string
	Source   string
}

const OwnedSkillsUpgradeAction = "roundfix skills install --target project"

// OwnedSkillReadiness records the facts every owned-skill surface renders.
type OwnedSkillReadiness struct {
	Skill   string
	Minimum string
	Found   string
	Source  string
	State   ReadinessState
	Detail  string
}

// ComparisonDetail names the complete blocking comparison.
func (readiness OwnedSkillReadiness) ComparisonDetail() string {
	return fmt.Sprintf(
		"below minimum: skill %q requires %s, found %s",
		readiness.Skill,
		readiness.Minimum,
		readiness.Found,
	)
}

// Diagnostic adds the owned-skill remediation to the comparison facts.
func (readiness OwnedSkillReadiness) Diagnostic() string {
	return readiness.ComparisonDetail() + "; upgrade: " + OwnedSkillsUpgradeAction
}

// BundleReadiness is the embedded owned-skill contract used by skills check.
type BundleReadiness struct {
	Owned       []OwnedSkillReadiness
	Diagnostics []Diagnostic
}

// Readiness compares a declared version against Roundfix's declared minimum.
func Readiness(declared SkillVersion, minimum string) (ReadinessState, error) {
	if strings.TrimSpace(declared.Source) == "" {
		return ReadinessUnversioned, fmt.Errorf("%w: source is unreachable", ErrSkillVersionUnresolvable)
	}
	if strings.TrimSpace(declared.Declared) == "" {
		return ReadinessUnversioned, nil
	}

	found, ok := parseSkillVersion(declared.Declared)
	if !ok {
		return ReadinessUnversioned, fmt.Errorf(
			"%w: source %q declares %q",
			ErrSkillVersionUnresolvable,
			declared.Source,
			declared.Declared,
		)
	}
	required, ok := parseSkillVersion(minimum)
	if !ok {
		return "", fmt.Errorf("invalid minimum skill version %q", minimum)
	}
	for index := range found {
		if found[index] > required[index] {
			return ReadinessSatisfies, nil
		}
		if found[index] < required[index] {
			return ReadinessBelow, nil
		}
	}
	return ReadinessSatisfies, nil
}

// ValidVersion reports whether version is a strict three-part skill version.
func ValidVersion(version string) bool {
	_, ok := parseSkillVersion(version)
	return ok
}

func parseSkillVersion(version string) ([3]uint64, bool) {
	var parsed [3]uint64
	parts := strings.Split(strings.TrimSpace(version), ".")
	if len(parts) != len(parsed) {
		return parsed, false
	}
	for index, part := range parts {
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return parsed, false
		}
		parsed[index] = value
	}
	return parsed, true
}
