"""Distribution and protected-version contract tests.

Suite: Roundfix-owned distribution identity
Invariant: every authoritative Roundfix distribution surface reports 0.0.1
while operational, external, third-party, and upstream versions stay unchanged.
Boundary IN: checked-in Go, npm, setup, skill, changelog, and release artifacts.
Boundary OUT: remote tags and GitHub Releases, which the Release Plan only reads.
"""

import json
import re
import unittest
from pathlib import Path


OWNED_VERSION = "0.0.1"
RELEASE_PLAN_SCHEMA = "roundfix.release-plan/0.0.1"
OWNED_SKILLS = {
    "archive-spec",
    "brainstorming",
    "business-analyst",
    "council",
    "evidence-gate",
    "implement-spec",
    "implement-task",
    "qa-gate",
    "roundfix",
    "setup-context-driven",
    "write-idea",
    "write-prd",
    "write-tasks",
    "write-techspec",
}
PLATFORM_PACKAGES = {
    "@roundfix/cli-darwin-arm64": "cli-darwin-arm64",
    "@roundfix/cli-darwin-x64": "cli-darwin-x64",
    "@roundfix/cli-linux-arm64": "cli-linux-arm64",
    "@roundfix/cli-linux-x64": "cli-linux-x64",
    "@roundfix/cli-win32-x64": "cli-win32-x64",
}
MAINTAINED_PROFILES = {
    "go-cli-tui": "go-cli-tui.json",
    "rust-cli": "rust-cli.json",
    "standard-typescript-monorepo": "standard-typescript-monorepo.json",
}
SETUP_SNAPSHOTS = {
    "go-cli": "go-cli.json",
    "rust-cli": "rust-cli.json",
    "typescript-bun": "typescript-bun.json",
}


def find_repository_root():
    skill_root = Path(__file__).resolve().parents[1]
    for candidate in (skill_root, *skill_root.parents):
        if (candidate / "go.mod").is_file() and (
            candidate / "dist" / "npm" / "roundfix" / "package.json"
        ).is_file():
            return candidate
    raise RuntimeError("repository root is not discoverable from the skill tree")


REPOSITORY_ROOT = find_repository_root()
SKILL_ROOT = Path(__file__).resolve().parents[1]


def read_json(path):
    return json.loads(path.read_text(encoding="utf-8"))


def require_match(test, path, pattern):
    text = path.read_text(encoding="utf-8")
    match = re.search(pattern, text, flags=re.MULTILINE | re.DOTALL)
    test.assertIsNotNone(match, f"{path} does not match {pattern!r}")
    return match.group(1)


def skill_version(test, path):
    return require_match(
        test,
        path,
        r"^---\n.*?^metadata:\n.*?^[ \t]+version:[ \t]+[\"']?([^\"'\n]+)",
    ).strip()


def assert_owned_version_matrix(observed, expected):
    mismatches = [
        f"{surface}: got {observed.get(surface)!r}, want {version!r}"
        for surface, version in sorted(expected.items())
        if observed.get(surface) != version
    ]
    if mismatches:
        raise AssertionError("owned version mismatch: " + "; ".join(mismatches))


class VersionContractTests(unittest.TestCase):
    def test_authoritative_distribution_surfaces_report_0_0_1(self):
        owned_versions = {
            "application": require_match(
                self,
                REPOSITORY_ROOT / "internal" / "app" / "version.go",
                r'^var Version = "([^"]+)"$',
            ),
            "setup-generation": require_match(
                self,
                SKILL_ROOT / "scripts" / "context_setup.py",
                r'^OWNED_VERSION_0_0_1 = "([^"]+)"$',
            ),
        }

        launcher = read_json(
            REPOSITORY_ROOT / "dist" / "npm" / "roundfix" / "package.json"
        )
        owned_versions["npm:roundfix"] = launcher["version"]
        platform_root = REPOSITORY_ROOT / "dist" / "npm" / "packages"
        self.assertEqual(
            {path.name for path in platform_root.iterdir() if path.is_dir()},
            set(PLATFORM_PACKAGES.values()),
        )
        self.assertEqual(
            set(launcher["optionalDependencies"]),
            set(PLATFORM_PACKAGES),
        )
        for package_name, directory in PLATFORM_PACKAGES.items():
            owned_versions[f"npm:{package_name}"] = read_json(
                REPOSITORY_ROOT
                / "dist"
                / "npm"
                / "packages"
                / directory
                / "package.json"
            )["version"]
            owned_versions[f"npm-dependency:{package_name}"] = launcher[
                "optionalDependencies"
            ][package_name]

        for tree in (".agents/skills", "skills"):
            for name in OWNED_SKILLS:
                owned_versions[f"{tree}:{name}"] = skill_version(
                    self,
                    REPOSITORY_ROOT / tree / name / "SKILL.md",
                )

        schema_expectations = {}
        profile_root = SKILL_ROOT / "assets" / "profiles"
        self.assertEqual(
            {path.name for path in profile_root.glob("*.json")},
            set(MAINTAINED_PROFILES.values()),
        )
        for profile_id, filename in MAINTAINED_PROFILES.items():
            profile = read_json(profile_root / filename)
            self.assertEqual(profile["id"], profile_id)
            profile_schema = "setup-context-driven/profile/0.0.1"
            owned_versions[f"profile-schema:{profile_id}"] = profile[
                "schemaVersion"
            ]
            schema_expectations[f"profile-schema:{profile_id}"] = profile_schema
            owned_versions[f"profile:{profile_id}"] = profile["version"]
            owned_versions[f"profile-marker:{profile_id}"] = profile["markerVersion"]

        setup_root = SKILL_ROOT / "assets" / "setups"
        self.assertEqual(
            {path.name for path in setup_root.glob("*.json")},
            set(SETUP_SNAPSHOTS.values()),
        )
        for setup_id, filename in SETUP_SNAPSHOTS.items():
            setup = read_json(setup_root / filename)
            self.assertEqual(setup["id"], setup_id)
            setup_schema = "setup-context-driven/setup-snapshot/0.0.1"
            owned_versions[f"setup-snapshot-schema:{setup_id}"] = setup[
                "schemaVersion"
            ]
            schema_expectations[
                f"setup-snapshot-schema:{setup_id}"
            ] = setup_schema
            owned_versions[f"setup-snapshot:{setup_id}"] = setup["version"]

        baseline_root = (
            SKILL_ROOT
            / "assets"
            / "source-baselines"
            / "baseline.standard-typescript-monorepo-0.0.1"
        )
        baseline_documents = {
            "source-baseline-index": (
                SKILL_ROOT / "assets" / "source-baselines" / "index.json",
                "setup-context-driven/source-baseline-index/0.0.1",
            ),
            "source-baseline": (
                baseline_root / "baseline.json",
                "setup-context-driven/source-baseline/0.0.1",
            ),
            "source-baseline-manifest": (
                baseline_root / "manifest.json",
                "setup-context-driven/source-baseline-manifest/0.0.1",
            ),
            "source-accounting": (
                baseline_root / "accounting.json",
                "setup-context-driven/source-accounting/0.0.1",
            ),
        }
        for name, (path, schema) in baseline_documents.items():
            document = read_json(path)
            owned_versions[f"{name}-schema"] = document["schemaVersion"]
            schema_expectations[f"{name}-schema"] = schema
            owned_versions[name] = document["version"]

        release_schema = require_match(
            self,
            REPOSITORY_ROOT / "internal" / "releaseplan" / "model.go",
            r'^const SchemaVersion = "([^"]+)"$',
        )
        owned_versions["schema:release-plan"] = release_schema
        schema_expectations["schema:release-plan"] = RELEASE_PLAN_SCHEMA
        changelog = (REPOSITORY_ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
        headings = re.findall(r"^## \[([^\]]+)\]", changelog, flags=re.MULTILINE)
        self.assertEqual(headings, [OWNED_VERSION])
        owned_versions["changelog"] = headings[0]

        expected_versions = {
            surface: schema_expectations.get(surface, OWNED_VERSION)
            for surface in owned_versions
        }
        self.assertEqual(len(owned_versions), 66)
        assert_owned_version_matrix(owned_versions, expected_versions)
        for surface in sorted(owned_versions):
            with self.subTest(mutated_surface=surface):
                mutated = dict(owned_versions)
                mutated[surface] = mutated[surface] + "-mutated"
                with self.assertRaisesRegex(AssertionError, re.escape(surface)):
                    assert_owned_version_matrix(mutated, expected_versions)

        self.assertNotEqual(release_schema, "roundfix.release-plan/v1")

    def test_changelog_restarts_at_only_0_0_1(self):
        changelog = (REPOSITORY_ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
        headings = re.findall(r"^## \[([^\]]+)\]", changelog, flags=re.MULTILINE)

        self.assertEqual(headings, [OWNED_VERSION])
        self.assertIn("[0.0.1]:", changelog)
        self.assertNotRegex(changelog, r"^## \[(?:0\.1\.0|0\.2\.0|0\.3\.0)\]")

    def test_protected_versions_match_the_operational_and_upstream_fixture(self):
        fixture = read_json(
            SKILL_ROOT
            / "tests"
            / "fixtures"
            / "version-contract"
            / "protected-versions.json"
        )

        run_database_schema = int(
            require_match(
                self,
                REPOSITORY_ROOT / "internal" / "store" / "store.go",
                r"^const schemaVersion = ([0-9]+)$",
            )
        )
        self.assertEqual(run_database_schema, fixture["runDatabaseSchema"])

        skills_lock = read_json(REPOSITORY_ROOT / "skills-lock.json")
        self.assertEqual(skills_lock["version"], fixture["skillsLockSchema"])

        external_lock = read_json(
            SKILL_ROOT / "assets" / "lock-hash-compatibility-v1.json"
        )
        self.assertEqual(
            {
                "schemaVersion": external_lock["schemaVersion"],
                "version": external_lock["version"],
            },
            fixture["externalLockHashCompatibility"],
        )

        for protocol in fixture["thirdPartyProtocols"]:
            with self.subTest(protocol=protocol["fragment"]):
                source = (
                    REPOSITORY_ROOT / protocol["path"]
                ).read_text(encoding="utf-8")
                self.assertIn(protocol["fragment"], source)

        self.assertTrue(
            OWNED_SKILLS.isdisjoint(fixture["upstreamSkills"]),
            "upstream fixture must not contain Roundfix-owned skills",
        )
        for name, version in sorted(fixture["upstreamSkills"].items()):
            with self.subTest(upstream_skill=name):
                self.assertEqual(
                    skill_version(
                        self,
                        REPOSITORY_ROOT / ".agents" / "skills" / name / "SKILL.md",
                    ),
                    version,
                )

    def test_release_workflow_validates_identity_without_history_deletion(self):
        workflow = (
            REPOSITORY_ROOT / ".github" / "workflows" / "release.yml"
        ).read_text(encoding="utf-8")
        self.assertIn(
            'OWNED_VERSION="$(jq -r .version dist/npm/roundfix/package.json)"',
            workflow,
        )
        self.assertIn('if [ "$TAG" != "$OWNED_VERSION" ]; then', workflow)

        release_sources = [
            REPOSITORY_ROOT / ".github" / "workflows" / "release.yml",
            REPOSITORY_ROOT / "internal" / "releaseplan" / "model.go",
            REPOSITORY_ROOT / "internal" / "releaseplan" / "reset.go",
            REPOSITORY_ROOT / "internal" / "cli" / "releaseplan_command.go",
            REPOSITORY_ROOT / "internal" / "cli" / "releaseplan_git_source.go",
        ]
        release_text = "\n".join(
            path.read_text(encoding="utf-8") for path in release_sources
        ).lower()
        for forbidden in (
            "gh release delete",
            "git tag -d",
            "git push --delete",
            "deletetag(",
            "deleterelease(",
        ):
            with self.subTest(forbidden=forbidden):
                self.assertNotIn(forbidden, release_text)


if __name__ == "__main__":
    unittest.main()
