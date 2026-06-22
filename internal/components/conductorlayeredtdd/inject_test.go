package conductorlayeredtdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestInjectInstallsWorkflowPromptsSkillsAndGitignore(t *testing.T) {
	workspace := t.TempDir()

	result, err := Inject(workspace, InjectOptions{IncludeMemorySkills: true})
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() changed = false, want true")
	}

	workflow := readTestFile(t, filepath.Join(workspace, "workflows", "conductor", "layered-tdd.yaml"))
	if !strings.Contains(workflow, "entry_point: preflight_tests") {
		t.Fatalf("workflow missing preflight entry point")
	}
	if !strings.Contains(workflow, "default: make audit") {
		t.Fatalf("workflow missing Makefile audit default")
	}

	prompt := readTestFile(t, filepath.Join(workspace, "workflows", "conductor", "prompts", "implementor.md"))
	if !strings.Contains(prompt, "approved layer boundary") {
		t.Fatalf("implementor prompt missing layer boundary rule")
	}

	skill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-tdd", "SKILL.md"))
	if !strings.Contains(skill, "Test-Driven Development") {
		t.Fatalf("TDD skill was not installed into .github/skills/conductor-tdd")
	}
	memorySkill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-memory-capture", "SKILL.md"))
	if !strings.Contains(memorySkill, "machine/agent-memory/projects/<project>/") {
		t.Fatalf("Markdown memory-capture skill was not installed into .github/skills/conductor-memory-capture")
	}
	recallSkill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-memory-recall", "SKILL.md"))
	if !strings.Contains(recallSkill, "without scanning the whole vault") {
		t.Fatalf("Markdown memory-recall skill was not installed into .github/skills/conductor-memory-recall")
	}
	consolidateSkill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-memory-consolidate", "SKILL.md"))
	if !strings.Contains(consolidateSkill, "canonical project memory files") {
		t.Fatalf("Markdown memory-consolidate skill was not installed into .github/skills/conductor-memory-consolidate")
	}

	instructionsPath := filepath.Join(workspace, ".github", "copilot-instructions.md")
	if _, err := os.Stat(instructionsPath); !os.IsNotExist(err) {
		t.Fatalf("Copilot instructions should not be created by this component; stat error = %v", err)
	}

	gitignore := readTestFile(t, filepath.Join(workspace, ".gitignore"))
	for _, entry := range []string{"workflows/conductor/", ".github/plans/", ".github/skills/conductor-*/"} {
		if !strings.Contains(gitignore, entry) {
			t.Fatalf(".gitignore missing %q", entry)
		}
	}
	if strings.Contains(gitignore, ".github/skills/\n") {
		t.Fatalf(".github/skills root should remain trackable")
	}
}

func TestInjectSkipsMemorySkillsWhenMemoryIsNotSelected(t *testing.T) {
	workspace := t.TempDir()

	if _, err := Inject(workspace, InjectOptions{}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(workspace, ".github", "skills", "conductor-memory-capture", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("memory skill should not be installed without markdown memory selection; stat error = %v", err)
	}
	baseSkill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-tdd", "SKILL.md"))
	if !strings.Contains(baseSkill, "Test-Driven Development") {
		t.Fatalf("base TDD skill was not installed")
	}
}

func TestInjectInstallsRequestedProjectSkillsWithConductorPrefix(t *testing.T) {
	workspace := t.TempDir()

	if _, err := Inject(workspace, InjectOptions{
		ProjectSkillIDs: []model.SkillID{model.SkillJavaDevelopment},
	}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	skill := readTestFile(t, filepath.Join(workspace, ".github", "skills", "conductor-java-development", "SKILL.md"))
	if !strings.Contains(skill, "Java Development") {
		t.Fatalf("Java development skill was not installed into .github/skills/conductor-java-development")
	}
}

func TestInjectPreservesLocalWorkflowAndSkillEdits(t *testing.T) {
	workspace := t.TempDir()
	workflowPath := filepath.Join(workspace, "workflows", "conductor", "layered-tdd.yaml")
	skillPath := filepath.Join(workspace, ".github", "skills", "conductor-tdd", "SKILL.md")
	for _, path := range []string{workflowPath, skillPath} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(workflowPath, []byte("local workflow edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skillPath, []byte("local skill edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Inject(workspace, InjectOptions{}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	if got := readTestFile(t, workflowPath); got != "local workflow edit\n" {
		t.Fatalf("workflow overwritten = %q", got)
	}
	if got := readTestFile(t, skillPath); got != "local skill edit\n" {
		t.Fatalf("skill overwritten = %q", got)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return string(data)
}
