package releaseplan

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"sort"
	"strconv"
	"strings"
)

const resetDigestSchema = "roundfix.release-reset-plan/v1"

// BuildReset inventories every stable tag and GitHub Release through read-only
// providers and returns a deterministic approval-required reset plan.
func BuildReset(ctx context.Context, request ResetRequest, source ResetInventorySource) (ResetPlan, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, err := ParseStableVersion(request.TargetVersion); err != nil {
		return ResetPlan{}, fmt.Errorf("parse release reset target version: %w", err)
	}
	if strings.TrimSpace(request.Target.Name) == "" || strings.TrimSpace(request.Target.CommitSHA) == "" {
		return ResetPlan{}, resetInventoryError(
			"validate release reset target",
			request.Target.Name,
			"resolve a clean committed target before planning the release history reset",
			errors.New("target revision or commit is empty"),
		)
	}
	if source == nil {
		return ResetPlan{}, resetInventoryError(
			"load release reset inventory",
			"",
			"configure read-only Git and GitHub inventory providers",
			errors.New("missing inventory source"),
		)
	}

	tags, err := source.Tags(ctx)
	if err != nil {
		return ResetPlan{}, resetInventoryError(
			"inventory local and remote stable tags",
			"",
			"restore complete Git tag access and rerun the Release Plan Command",
			err,
		)
	}
	tags, err = normalizeResetTags(tags)
	if err != nil {
		return ResetPlan{}, err
	}

	releases, err := source.Releases(ctx)
	if err != nil {
		return ResetPlan{}, resetInventoryError(
			"inventory GitHub Releases",
			"",
			"restore complete paginated GitHub Release access and rerun the Release Plan Command",
			err,
		)
	}
	releases, err = normalizeResetReleases(releases)
	if err != nil {
		return ResetPlan{}, err
	}
	resolveReleaseTargetCommits(releases, tags)

	plan := ResetPlan{
		SchemaVersion:  SchemaVersion,
		State:          StateApprovalRequired,
		TargetVersion:  request.TargetVersion,
		TargetRevision: request.Target.Name,
		TargetCommit:   request.Target.CommitSHA,
		Tags:           tags,
		Releases:       releases,
	}
	plan.PlanDigest = resetPlanDigest(plan)
	plan.Approval = Approval{
		Required:        true,
		Increment:       IncrementNone,
		ProposedVersion: request.TargetVersion,
		Question:        fmt.Sprintf("Approve a separate release history reset to %s for plan %s?", request.TargetVersion, plan.PlanDigest),
	}
	return plan, nil
}

func normalizeResetTags(input []TagRef) ([]TagRef, error) {
	tags := append([]TagRef(nil), input...)
	identities := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		if _, err := ParseStableVersion(tag.Name); err != nil {
			return nil, resetInventoryError(
				"validate release reset tag",
				tag.Name,
				"repair the stable tag inventory and rerun the Release Plan Command",
				err,
			)
		}
		if tag.Source != TagSourceLocal && tag.Source != TagSourceRemote {
			return nil, resetInventoryError(
				"validate release reset tag source",
				string(tag.Source),
				"provide local or remote as the tag source",
				errors.New("unknown tag source"),
			)
		}
		if tag.Source == TagSourceRemote && strings.TrimSpace(tag.Remote) == "" {
			return nil, resetInventoryError(
				"validate release reset remote tag",
				tag.Name,
				"include the remote name for every remote tag",
				errors.New("remote name is empty"),
			)
		}
		if strings.TrimSpace(tag.Ref) == "" || strings.TrimSpace(tag.TargetCommit) == "" || strings.TrimSpace(tag.ImmutableID) == "" {
			return nil, resetInventoryError(
				"validate release reset tag",
				tag.Name,
				"include the ref, immutable identity, and target commit for every stable tag",
				errors.New("tag identity is incomplete"),
			)
		}
		if _, exists := identities[tag.ImmutableID]; exists {
			return nil, duplicateResetInventoryIdentity("tag", tag.ImmutableID)
		}
		identities[tag.ImmutableID] = struct{}{}
	}
	sort.Slice(tags, func(left, right int) bool {
		leftTag := tags[left]
		rightTag := tags[right]
		if leftTag.Name != rightTag.Name {
			return leftTag.Name < rightTag.Name
		}
		if leftTag.Source != rightTag.Source {
			return leftTag.Source < rightTag.Source
		}
		if leftTag.Remote != rightTag.Remote {
			return leftTag.Remote < rightTag.Remote
		}
		if leftTag.Ref != rightTag.Ref {
			return leftTag.Ref < rightTag.Ref
		}
		return leftTag.ImmutableID < rightTag.ImmutableID
	})
	return tags, nil
}

func normalizeResetReleases(input []ReleaseRef) ([]ReleaseRef, error) {
	releases := append([]ReleaseRef(nil), input...)
	immutableIDs := make(map[string]struct{}, len(releases))
	nodeIDs := make(map[string]struct{}, len(releases))
	databaseIDs := make(map[int64]struct{}, len(releases))
	for _, release := range releases {
		if release.ID <= 0 || strings.TrimSpace(release.NodeID) == "" ||
			strings.TrimSpace(release.TagName) == "" || strings.TrimSpace(release.ImmutableID) == "" {
			return nil, resetInventoryError(
				"validate GitHub Release inventory",
				release.TagName,
				"include the database ID, node ID, tag name, and immutable identity for every GitHub Release",
				errors.New("GitHub Release identity is incomplete"),
			)
		}
		if _, exists := immutableIDs[release.ImmutableID]; exists {
			return nil, duplicateResetInventoryIdentity("GitHub Release", release.ImmutableID)
		}
		if _, exists := nodeIDs[release.NodeID]; exists {
			return nil, duplicateResetInventoryIdentity("GitHub Release node", release.NodeID)
		}
		if _, exists := databaseIDs[release.ID]; exists {
			return nil, duplicateResetInventoryIdentity("GitHub Release database", strconv.FormatInt(release.ID, 10))
		}
		immutableIDs[release.ImmutableID] = struct{}{}
		nodeIDs[release.NodeID] = struct{}{}
		databaseIDs[release.ID] = struct{}{}
	}
	sort.Slice(releases, func(left, right int) bool {
		leftRelease := releases[left]
		rightRelease := releases[right]
		if leftRelease.TagName != rightRelease.TagName {
			return leftRelease.TagName < rightRelease.TagName
		}
		if leftRelease.ID != rightRelease.ID {
			return leftRelease.ID < rightRelease.ID
		}
		return leftRelease.ImmutableID < rightRelease.ImmutableID
	})
	return releases, nil
}

func resolveReleaseTargetCommits(releases []ReleaseRef, tags []TagRef) {
	tagTargets := make(map[string]map[string]struct{})
	for _, tag := range tags {
		targets := tagTargets[tag.Name]
		if targets == nil {
			targets = map[string]struct{}{}
			tagTargets[tag.Name] = targets
		}
		targets[tag.TargetCommit] = struct{}{}
	}
	for index := range releases {
		if releases[index].TargetCommit != "" {
			continue
		}
		targets := tagTargets[releases[index].TagName]
		if len(targets) != 1 {
			continue
		}
		for target := range targets {
			releases[index].TargetCommit = target
		}
	}
}

func resetPlanDigest(plan ResetPlan) string {
	digest := sha256.New()
	writeResetDigestField(digest, "schema", resetDigestSchema)
	writeResetDigestField(digest, "targetVersion", plan.TargetVersion)
	writeResetDigestField(digest, "targetRevision", plan.TargetRevision)
	writeResetDigestField(digest, "targetCommit", plan.TargetCommit)
	writeResetDigestField(digest, "tagCount", strconv.Itoa(len(plan.Tags)))
	for _, tag := range plan.Tags {
		writeResetDigestField(digest, "tag.name", tag.Name)
		writeResetDigestField(digest, "tag.source", string(tag.Source))
		writeResetDigestField(digest, "tag.remote", tag.Remote)
		writeResetDigestField(digest, "tag.ref", tag.Ref)
		writeResetDigestField(digest, "tag.immutableID", tag.ImmutableID)
		writeResetDigestField(digest, "tag.targetCommit", tag.TargetCommit)
	}
	writeResetDigestField(digest, "releaseCount", strconv.Itoa(len(plan.Releases)))
	for _, release := range plan.Releases {
		writeResetDigestField(digest, "release.id", strconv.FormatInt(release.ID, 10))
		writeResetDigestField(digest, "release.nodeID", release.NodeID)
		writeResetDigestField(digest, "release.name", release.Name)
		writeResetDigestField(digest, "release.tagName", release.TagName)
		writeResetDigestField(digest, "release.targetCommitish", release.TargetCommitish)
		writeResetDigestField(digest, "release.targetCommit", release.TargetCommit)
		writeResetDigestField(digest, "release.immutableID", release.ImmutableID)
	}
	return fmt.Sprintf("sha256:%x", digest.Sum(nil))
}

func writeResetDigestField(digest hash.Hash, name string, value string) {
	fmt.Fprintf(digest, "%d:%s%d:%s", len(name), name, len(value), value)
}

func resetInventoryError(operation string, identity string, nextAction string, err error) ResetInventoryError {
	return ResetInventoryError{
		Operation:  operation,
		Identity:   identity,
		NextAction: nextAction,
		Err:        err,
		Kinds:      []error{ErrIncompleteResetInventory},
	}
}

func duplicateResetInventoryIdentity(kind string, identity string) ResetInventoryError {
	return ResetInventoryError{
		Operation:  "validate unique " + kind + " identity",
		Identity:   identity,
		NextAction: "repair the duplicate inventory result and rerun the Release Plan Command",
		Err:        ErrDuplicateInventoryIdentity,
		Kinds:      []error{ErrIncompleteResetInventory},
	}
}
