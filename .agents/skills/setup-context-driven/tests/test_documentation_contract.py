"""Public Baseline documentation contract tests.

Suite: Context-Driven Baseline operating documentation
Invariant: public guides and owned skills describe the shipped Go CLI contract with canonical terms and portable examples.
Boundary IN: glossary, README, user guides, and canonical/distributed setup and Roundfix skills.
Boundary OUT: command parsing and strict JSON parsing, which the Go documentation contract suite owns.
"""

import unittest
from pathlib import Path


CANONICAL_TERMS = (
    "Context-Driven Baseline",
    "Baseline Profile",
    "Baseline Command",
    "Source Baseline",
    "Source Baseline Entry",
    "Normative Clause Manifest",
    "Operational Contract",
    "Repository Capability",
    "Skill Activation",
    "Decision Plan",
    "Change Plan",
    "Setup Manifest",
    "Baseline Readoption",
    "Repository-Specific Normative Rules",
    "Repository Skill Set",
)


def find_repository_root():
    skill_root = Path(__file__).resolve().parents[1]
    for candidate in (skill_root, *skill_root.parents):
        if (candidate / "go.mod").is_file() and (candidate / "CONTEXT.md").is_file():
            return candidate
    raise RuntimeError("repository root is not discoverable from the skill tree")


REPOSITORY_ROOT = find_repository_root()
README = REPOSITORY_ROOT / "README.md"
CONTEXT_GUIDE = REPOSITORY_ROOT / "docs/user-guide/context-driven-development.md"
COMMAND_GUIDE = REPOSITORY_ROOT / "docs/user-guide/commands.md"
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
    def test_canonical_terms_and_decision_input_are_actionable(self):
        glossary = read(REPOSITORY_ROOT / "CONTEXT.md")
        guide = read(CONTEXT_GUIDE)

        for term in CANONICAL_TERMS:
            with self.subTest(term=term):
                self.assertIn(f"**{term}**:", glossary)
                self.assertIn(term, guide)

        for fragment in (
            "setup-context-driven/decisions/0.0.1",
            '"schemaVersion"',
            '"decisions"',
            '"preservation.mode"',
            "--decision-file baseline-decisions.json",
            "--confirm-plan sha256:<64-lowercase-hex>",
            "roundfix/baseline-plan/v1",
            "roundfix/baseline-result/v1",
        ):
            with self.subTest(fragment=fragment):
                self.assertIn(fragment, guide)

    def test_public_command_family_is_the_only_documented_runtime(self):
        documents = {
            "README": read(README),
            "Context-Driven user guide": read(CONTEXT_GUIDE),
            "command reference": read(COMMAND_GUIDE),
            "setup skill": read(SETUP_SKILL),
            "Roundfix skill": read(ROUNDFIX_SKILL),
        }

        for document_name, document in documents.items():
            with self.subTest(document=document_name):
                self.assertIn("roundfix baseline", document)
                self.assertNotIn("context_setup.py", document)
                self.assertNotIn("python3", document.lower())

        combined = " ".join("\n".join(documents.values()).split())
        for fragment in (
            "roundfix baseline plan",
            "roundfix baseline apply",
            "roundfix baseline profile init",
            "roundfix baseline skills restore",
            "roundfix baseline assets sync",
            "stdout",
            "stderr",
            "profile expectation",
            "repository command",
            "recommendation",
            "no partial plan",
            "no independent setup engine",
        ):
            with self.subTest(fragment=fragment):
                self.assertIn(" ".join(fragment.split()), combined)

    def test_recovery_security_and_migration_contracts_are_complete(self):
        guide = read(CONTEXT_GUIDE)

        for fragment in (
            "### First adoption",
            "Greenfield",
            "Preservation",
            "profile change",
            "rejected plans",
            "Cross-clone",
            "Stale plan",
            "Unsafe carrier",
            "Interrupted transaction",
            "Incomplete rollback",
            "Repository Skill Set restoration",
            "Canonical asset synchronization",
            "2 MiB",
            "256 entries",
            "512 KiB",
            "never:",
            "Migrate from the script-backed setup skill",
            "There is no independent skill",
        ):
            with self.subTest(fragment=fragment):
                self.assertIn(fragment, guide)

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
                README,
                CONTEXT_GUIDE,
                COMMAND_GUIDE,
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
