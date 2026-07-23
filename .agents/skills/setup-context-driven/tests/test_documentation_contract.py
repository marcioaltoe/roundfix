"""0.0.1 documentation contract tests.

Suite: Context-Driven Baseline operating documentation
Invariant: shipped guides describe the complete project-agnostic 0.0.1 setup, Readoption, and release-reset contracts.
Boundary IN: canonical glossary, user guides, and canonical/distributed setup and Roundfix skills.
Boundary OUT: CLI behavior and asset semantics, which their focused Go and Python suites own.
"""

import unittest
from pathlib import Path


CANONICAL_TERMS = (
    "Source Baseline",
    "Normative Clause Manifest",
    "Operational Contract",
    "Repository Capability",
    "Skill Activation",
    "HTTP Contract Decision",
    "Baseline Readoption",
    "Repository-Specific Normative Rules",
)
REQUIRED_STACK = (
    "TypeScript",
    "Bun",
    "Turborepo",
    "Vite",
    "React",
    "Hono",
    "Drizzle",
    "Zod",
    "Tailwind",
    "shadcn",
    "TanStack Query",
    "TanStack Router",
    "Better Auth",
    "PostgreSQL",
    "LogTape",
    "Oxlint",
    "Oxfmt",
    "Vitest",
)
ACTIVATION_BUNDLES = (
    "bundle.production-code",
    "bundle.frontend-react",
    "bundle.frontend-ui-quality",
    "bundle.hono-endpoint",
    "bundle.hono-endpoint-persistence",
    "bundle.testing",
    "bundle.debugging",
    "bundle.security",
    "bundle.qa",
    "bundle.delivery",
)


def find_repository_root():
    skill_root = Path(__file__).resolve().parents[1]
    for candidate in (skill_root, *skill_root.parents):
        if (candidate / "go.mod").is_file() and (candidate / "CONTEXT.md").is_file():
            return candidate
    raise RuntimeError("repository root is not discoverable from the skill tree")


REPOSITORY_ROOT = find_repository_root()
CONTEXT_GUIDE = REPOSITORY_ROOT / "docs/user-guide/context-driven-development.md"
CONFIGURATION_GUIDE = REPOSITORY_ROOT / "docs/user-guide/configuration.md"
RELEASE_RUNBOOK = REPOSITORY_ROOT / "docs/user-guide/release-runbook.md"
SETUP_SKILL = REPOSITORY_ROOT / ".agents/skills/setup-context-driven/SKILL.md"
ROUNDFIX_SKILL = REPOSITORY_ROOT / ".agents/skills/roundfix/SKILL.md"


def read(path):
    return path.read_text(encoding="utf-8")


def tree_bytes(root):
    return {
        path.relative_to(root).as_posix(): path.read_bytes()
        for path in sorted(root.rglob("*"))
        if path.is_file()
    }


class DocumentationContractTests(unittest.TestCase):
    def test_canonical_terms_and_readoption_workflow_are_actionable(self):
        glossary = read(REPOSITORY_ROOT / "CONTEXT.md")
        guide = read(CONTEXT_GUIDE)

        for term in CANONICAL_TERMS:
            with self.subTest(term=term):
                self.assertIn(f"**{term}**:", glossary)
                self.assertIn(term, guide)

        for fragment in (
            "setup-context-driven/decisions/0.0.1",
            '"sourceBaseline"',
            '"dispositions"',
            '"entryDigest"',
            "managed-entry",
            "repository-document",
            "repository-rules",
            "rejected",
            "--decision-file <decision-file>",
            "--confirm-plan <planDigest>",
            "formatter → Verification → audit → reapply",
        ):
            with self.subTest(fragment=fragment):
                self.assertIn(fragment, guide)

    def test_standard_typescript_profile_and_capability_policy_are_exact(self):
        guide = read(CONTEXT_GUIDE)
        setup_skill = read(SETUP_SKILL)
        combined = " ".join((guide + "\n" + setup_skill).split())
        normalized_guide = " ".join(guide.split())

        for fragment in (
            "standard-typescript-monorepo",
            "packages/frontend",
            "packages/backend",
            "public system boundary",
            "`domain`, `application`, and `infrastructure`",
            "thin HTTP handlers",
            "HTTP-independent",
            "Drizzle-owned",
            "`REST`",
            "`Post-only`",
            "`scope`",
            "`methods`",
            "`owner`",
            "`reason`",
            "Inngest",
            "Docker",
        ):
            with self.subTest(fragment=fragment):
                self.assertIn(" ".join(fragment.split()), combined)

        for stack_entry in REQUIRED_STACK:
            with self.subTest(stack_entry=stack_entry):
                self.assertIn(stack_entry, guide)
        for bundle in ACTIVATION_BUNDLES:
            with self.subTest(bundle=bundle):
                self.assertIn(bundle, guide)

        self.assertIn("Context7 and Exa are required", normalized_guide)
        self.assertIn(
            "Firecrawl, `rtk`, and `rg` are recommended",
            normalized_guide,
        )
        self.assertIn("three to seven varied searches", normalized_guide)
        self.assertLess(
            normalized_guide.index("Search local repository code"),
            normalized_guide.index("Use Context7"),
        )
        self.assertIn("capability.recommended.missing", normalized_guide)
        self.assertIn("do not block", setup_skill)

    def test_owned_and_protected_version_surfaces_are_unambiguous(self):
        guide = read(CONTEXT_GUIDE)
        configuration = read(CONFIGURATION_GUIDE)
        setup_skill = read(SETUP_SKILL)
        combined = guide + "\n" + configuration + "\n" + setup_skill

        for owned in (
            "Roundfix CLI",
            "npm",
            "Context-Driven Baseline",
            "Source Baseline",
            "schemas",
            "manifests",
            "profiles",
            "Roundfix-owned distributed skills",
            "Release Plan JSON schema",
            "changelog",
        ):
            with self.subTest(owned=owned):
                self.assertIn(owned, combined)
        for protected in (
            "User Config",
            "Project Config",
            "Runs",
            "Run Database",
            "PRAGMA user_version",
            "skills-lock.json",
            "upstream-managed skill",
            "third-party protocol",
            "Git history",
            "Specs",
            "ADRs",
        ):
            with self.subTest(protected=protected):
                self.assertIn(protected, combined)

    def test_release_reset_guidance_stops_at_a_fresh_read_only_plan(self):
        documents = (
            (RELEASE_RUNBOOK.name, read(RELEASE_RUNBOOK)),
            (ROUNDFIX_SKILL.name, read(ROUNDFIX_SKILL)),
        )

        for document_name, document in documents:
            normalized = " ".join(document.split())
            for fragment in (
                "roundfix release plan --reset-to v0.0.1",
                "`--from`",
                "`--to`",
                "`--impact`",
                "`--reason`",
                "local and remote stable tag",
                "GitHub Release",
                "complete pagination",
                "planDigest",
                "approval_required",
                "read-only",
                "fresh inventory",
                "separate destructive release operation",
                "explicit human approval",
                "no tag or GitHub Release deletion action",
            ):
                with self.subTest(document=document_name, fragment=fragment):
                    self.assertIn(" ".join(fragment.split()), normalized)
            self.assertRegex(normalized, r"(?:exit|exits) `3`")

    def test_examples_are_project_agnostic_and_skill_trees_are_identical(self):
        documented_surfaces = "\n".join(
            read(path)
            for path in (
                CONTEXT_GUIDE,
                CONFIGURATION_GUIDE,
                RELEASE_RUNBOOK,
                SETUP_SKILL,
                ROUNDFIX_SKILL,
            )
        )
        for forbidden in (
            "/Users/",
            "/home/",
            "typescript-bun-monorepo",
            "Gesttione",
            "Fluxus",
            "Tax POC",
        ):
            with self.subTest(forbidden=forbidden):
                self.assertNotIn(forbidden, documented_surfaces)

        for skill_name in ("setup-context-driven", "roundfix"):
            with self.subTest(skill=skill_name):
                canonical = REPOSITORY_ROOT / ".agents/skills" / skill_name
                distributed = REPOSITORY_ROOT / "skills" / skill_name
                self.assertEqual(tree_bytes(canonical), tree_bytes(distributed))


if __name__ == "__main__":
    unittest.main()
