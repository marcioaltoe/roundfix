package baseline

import (
	"regexp"
	"strings"
)

var makeRecipeInvocation = regexp.MustCompile(`(?:^|[;&|])\s*[@+\-]*\s*(?:rtk\s+)?(?:\$\(MAKE\)|\$\{MAKE\}|make)\s+([^;&|]+)`)

// makeTargetReach follows prerequisites and literal recursive Make invocations
// in one Makefile. It does not evaluate variables or included files.
func makeTargetReach(makefile []byte, target string) (targets map[string]struct{}, recipes []string) {
	prerequisites := make(map[string][]string)
	targetRecipes := make(map[string][]string)
	var current []string
	text := strings.ReplaceAll(string(makefile), "\\\n", " ")
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "\t") {
			for _, name := range current {
				targetRecipes[name] = append(targetRecipes[name], strings.TrimSpace(line))
			}
			continue
		}
		line, _, _ = strings.Cut(line, "#")
		if strings.TrimSpace(line) == "" {
			continue
		}
		current = nil
		names, rest, rule := strings.Cut(line, ":")
		if !rule || strings.Contains(names, "=") || strings.HasPrefix(rest, "=") {
			continue
		}
		for _, name := range strings.Fields(names) {
			if makeTargetName.MatchString(name) {
				current = append(current, name)
			}
		}
		deps, recipe, inline := strings.Cut(strings.TrimPrefix(rest, ":"), ";")
		for _, name := range current {
			for _, dep := range strings.Fields(deps) {
				if makeTargetName.MatchString(dep) {
					prerequisites[name] = append(prerequisites[name], dep)
				}
			}
			if inline {
				targetRecipes[name] = append(targetRecipes[name], strings.TrimSpace(recipe))
			}
		}
	}
	targets = make(map[string]struct{})
	pending := []string{target}
	for len(pending) != 0 {
		name := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if _, visited := targets[name]; visited {
			continue
		}
		targets[name] = struct{}{}
		pending = append(pending, prerequisites[name]...)
		for _, recipe := range targetRecipes[name] {
			recipes = append(recipes, recipe)
			for _, match := range makeRecipeInvocation.FindAllStringSubmatch(recipe, -1) {
				arguments, _, _ := strings.Cut(match[1], "#")
				for _, argument := range strings.Fields(arguments) {
					if makeTargetName.MatchString(argument) {
						pending = append(pending, argument)
					}
				}
			}
		}
	}
	return targets, recipes
}

func gatePartReached(command string, targets map[string]struct{}, recipes []string) bool {
	command = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(command), "rtk "))
	fields := strings.Fields(command)
	if len(fields) == 2 && fields[0] == "make" {
		_, reached := targets[fields[1]]
		return reached
	}
	for _, recipe := range recipes {
		recipe = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(recipe), "@-+"))
		recipe = strings.TrimPrefix(recipe, "rtk ")
		if command != "" && strings.Contains(recipe, command) {
			return true
		}
	}
	return false
}
