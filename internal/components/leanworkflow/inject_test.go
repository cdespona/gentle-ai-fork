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
			Model      string         `json:"model"`
			Variant    string         `json:"variant"`
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

	orchestratorPrompt := readTestFile(t, filepath.Join(home, ".config", "opencode", "lean-workflow", "prompts", "lean-workflow-orchestrator.md"))
	if !strings.Contains(orchestratorPrompt, "Delegate by phase objective and artifact path") {
		t.Fatalf("orchestrator prompt missing delegation contract")
	}
	if !strings.Contains(orchestratorPrompt, "Do not prescribe replacement headings for `01-requirements.md`") {
		t.Fatalf("orchestrator prompt missing requirements-template guard")
	}

	grillerPrompt := readTestFile(t, filepath.Join(home, ".config", "opencode", "lean-workflow", "prompts", "requirements-griller.md"))
	if !strings.Contains(grillerPrompt, "Treat the orchestrator handoff as context") {
		t.Fatalf("requirements griller prompt missing handoff guard")
	}
	if !strings.Contains(grillerPrompt, "Put user-facing questions in `## Blocking Questions`") {
		t.Fatalf("requirements griller prompt missing blocking-question instruction")
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

func TestInjectAppliesLeanWorkflowModelAssignments(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()

	assignments := map[string]model.ModelAssignment{
		"lean-workflow-orchestrator": {ProviderID: "anthropic", ModelID: "claude-sonnet-4", Effort: "medium"},
		"implementor":                {ProviderID: "openai", ModelID: "gpt-5"},
	}

	if _, err := Inject(home, workspace, opencode.NewAdapter(), InjectOptions{OpenCodeModelAssignments: assignments}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	agents := readAgentSettings(t, filepath.Join(home, ".config", "opencode", "opencode.json"))
	if got := agents["lean-workflow-orchestrator"].Model; got != "anthropic/claude-sonnet-4" {
		t.Fatalf("orchestrator model = %q, want anthropic/claude-sonnet-4", got)
	}
	if got := agents["lean-workflow-orchestrator"].Variant; got != "medium" {
		t.Fatalf("orchestrator variant = %q, want medium", got)
	}
	if got := agents["implementor"].Model; got != "openai/gpt-5" {
		t.Fatalf("implementor model = %q, want openai/gpt-5", got)
	}
	if got := agents["implementor"].Variant; got != "" {
		t.Fatalf("implementor variant = %q, want empty", got)
	}
	if got := agents["planner"].Model; got != "" {
		t.Fatalf("unassigned planner model = %q, want empty", got)
	}
}

func TestInjectAppliesRootModelToNewLeanWorkflowAgents(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, []byte(`{"model":"anthropic/claude-haiku-3-5"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Inject(home, workspace, opencode.NewAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	agents := readAgentSettings(t, settingsPath)
	for _, name := range []string{
		"lean-workflow-orchestrator",
		"requirements-griller",
		"planner",
		"test-guidance-planner",
		"todo-generator",
		"implementor",
	} {
		if got := agents[name].Model; got != "anthropic/claude-haiku-3-5" {
			t.Fatalf("%s model = %q, want root model", name, got)
		}
	}
}

func TestInjectPreservesExistingLeanWorkflowAgentModel(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, []byte(`{
  "model": "anthropic/claude-haiku-3-5",
  "agent": {
    "implementor": {"model": "anthropic/claude-opus-4", "variant": "high"}
  }
}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Inject(home, workspace, opencode.NewAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	agents := readAgentSettings(t, settingsPath)
	if got := agents["implementor"].Model; got != "anthropic/claude-opus-4" {
		t.Fatalf("existing implementor model = %q, want preserved model", got)
	}
	if got := agents["implementor"].Variant; got != "high" {
		t.Fatalf("existing implementor variant = %q, want preserved variant", got)
	}
	if got := agents["planner"].Model; got != "anthropic/claude-haiku-3-5" {
		t.Fatalf("new planner model = %q, want root fallback", got)
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

type testAgentSettings struct {
	Model   string `json:"model"`
	Variant string `json:"variant"`
}

func readAgentSettings(t *testing.T, settingsPath string) map[string]testAgentSettings {
	t.Helper()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}
	var settings struct {
		Agent map[string]testAgentSettings `json:"agent"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("Unmarshal(%q) error = %v", settingsPath, err)
	}
	return settings.Agent
}

type fakeAdapter struct{ *opencode.Adapter }

func (fakeAdapter) Agent() model.AgentID { return "fake" }
