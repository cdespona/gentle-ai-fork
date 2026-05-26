package leanworkflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestInjectInstallsLeanWorkflowAgentsTemplatesAndGitignore(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	adapter := opencode.NewAdapter()

	result, err := Inject(home, workspace, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() changed = false, want true")
	}

	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}
	var settings struct {
		Agent map[string]struct {
			Mode       string         `json:"mode"`
			Prompt     string         `json:"prompt"`
			Permission map[string]any `json:"permission"`
		} `json:"agent"`
		Permission map[string]map[string]string `json:"permission"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("Unmarshal(opencode.json) error = %v", err)
	}

	for _, name := range []string{
		"lean-workflow-orchestrator",
		"requirements-griller",
		"planner",
		"test-guidance-planner",
		"todo-generator",
		"implementor",
	} {
		agent, ok := settings.Agent[name]
		if !ok {
			t.Fatalf("settings missing agent %q", name)
		}
		if !strings.HasPrefix(agent.Prompt, "{file:") {
			t.Fatalf("agent %q prompt = %q, want file ref", name, agent.Prompt)
		}
	}
	if settings.Agent["lean-workflow-orchestrator"].Mode != "primary" {
		t.Fatalf("orchestrator mode = %q, want primary", settings.Agent["lean-workflow-orchestrator"].Mode)
	}
	if got := settings.Permission["bash"]["*"]; got != "ask" {
		t.Fatalf("bash permission = %q, want ask", got)
	}

	readmePath := filepath.Join(home, ".config", "opencode", "lean-workflow", "README.md")
	readme := readTestFile(t, readmePath)
	if !strings.Contains(readme, "Acceptance Test Gate") {
		t.Fatalf("README missing acceptance-test guidance")
	}

	templatePath := filepath.Join(workspace, ".github", "lean-workflow-templates", "04-todo.md")
	template := readTestFile(t, templatePath)
	if !strings.Contains(template, "Acceptance test status") {
		t.Fatalf("todo template missing acceptance status input")
	}

	gitignore := readTestFile(t, filepath.Join(workspace, ".gitignore"))
	for _, entry := range []string{".github/plans/", ".github/lean-workflow-templates/"} {
		if !strings.Contains(gitignore, entry) {
			t.Fatalf(".gitignore missing %q", entry)
		}
	}
}

func TestInjectPreservesLocalTemplateEdits(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	templatePath := filepath.Join(workspace, ".github", "lean-workflow-templates", "01-requirements.md")
	if err := os.MkdirAll(filepath.Dir(templatePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(templatePath, []byte("local edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Inject(home, workspace, opencode.NewAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	got := readTestFile(t, templatePath)
	if got != "local edit\n" {
		t.Fatalf("template overwritten = %q", got)
	}
}

func TestInjectSkipsNonOpenCodeAdapters(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	adapter := fakeAdapter{}

	result, err := Inject(home, workspace, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed || len(result.Files) != 0 {
		t.Fatalf("Inject() = %#v, want no-op", result)
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

type fakeAdapter struct{ *opencode.Adapter }

func (fakeAdapter) Agent() model.AgentID { return "fake" }
