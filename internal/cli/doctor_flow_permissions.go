package cli

import (
	"fmt"
	"strings"

	"roundfix/internal/flowpermissions"
)

// flowPermissionPlans reads only the detected project files. An unsafe target
// is advisory too; Preview refuses its link before opening anything outside Git.
func flowPermissionPlans(repository string) []flowpermissions.Plan {
	if repository == "" {
		return nil
	}
	var plans []flowpermissions.Plan
	for _, harness := range flowpermissions.Present(repository) {
		plan, err := flowpermissions.Preview(repository, harness, flowpermissions.ActionApply)
		if err != nil {
			plan.State = flowpermissions.Unreadable
		}
		plans = append(plans, plan)
	}
	return plans
}

func doctorFlowPermissionsResult(repository string) CheckResult {
	result := CheckResult{Name: HealthCheckFlowPermissions, Status: CheckStatusSkipped}
	if repository == "" {
		result.Detail = "outside a Git repository"
		return result
	}
	plans := flowPermissionPlans(repository)
	if len(plans) == 0 {
		result.Detail = "no harness files in this repository"
		return result
	}
	var parts []string
	code := ""
	for _, plan := range plans {
		state := strings.ReplaceAll(string(plan.State), "_", " ")
		switch plan.State {
		case flowpermissions.Behind:
			state = fmt.Sprintf("behind, %d to add, %d to remove", len(plan.Add)+len(plan.AddGuards), len(plan.Remove))
			code = "DR-FLOW-PERMISSIONS-BEHIND"
		case flowpermissions.Unreadable:
			state = "not readable"
			if code == "" {
				code = "DR-FLOW-PERMISSIONS-UNREADABLE"
			}
		case flowpermissions.Current, flowpermissions.NotNeeded, flowpermissions.OwnRules:
			result.Status = CheckStatusOK
		}
		if (plan.State == flowpermissions.Behind || plan.State == flowpermissions.Unreadable) && result.NextAction == "" {
			result.NextAction = permissionsPreviewCommand(plan)
		}
		parts = append(parts, string(plan.Harness)+": "+state)
	}
	result.Detail = strings.Join(parts, "; ")
	if code != "" {
		result.Status = CheckStatusWarn
		result.Detail = code + ": " + result.Detail
	} else if result.Status == CheckStatusSkipped {
		result.Detail += "; to see the list: roundfix permissions"
	}
	return result
}
