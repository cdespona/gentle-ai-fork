package markdownmemory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/agents/codex"
	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
)

func TestEnsureTemplatesCreatesMissingFilesWithoutOverwriting(t *testing.T) {
	vault := t.TempDir()
	cfg := Config{VaultRoot: vault, Namespace: "machine/agent-memory", Project: "gentle-ai-fork"}
	readmePath := filepath.Join(vault, "machine", "agent-memory", "README.md")
	if err := os.MkdirAll(filepath.Dir(readmePath), 0o755); err != nil {
		t.Fatal(err)
	}
	userContent := []byte("user-authored memory contract\n")
	if err := os.WriteFile(readmePath, userContent, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := EnsureTemplates(cfg)
	if err != nil {
		t.Fatalf("EnsureTemplates() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("EnsureTemplates() Changed = false, want true for newly created missing files")
	}

	after, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(userContent) {
		t.Fatalf("README.md was overwritten:\n%s", string(after))
	}
	if _, err := os.Stat(filepath.Join(vault, "machine", "agent-memory", "projects", "gentle-ai-fork", "index.md")); err != nil {
		t.Fatalf("project index was not created: %v", err)
	}
}

func TestInjectOpenCodeAndCodexReceiveProtocol(t *testing.T) {
	for _, tt := range []struct {
		name       string
		promptPath func(string) string
		adapter    agents.Adapter
	}{
		{
			name:       "opencode",
			promptPath: func(home string) string { return filepath.Join(home, ".config", "opencode", "AGENTS.md") },
			adapter:    opencode.NewAdapter(),
		},
		{
			name:       "codex",
			promptPath: func(home string) string { return filepath.Join(home, ".codex", "agents.md") },
			adapter:    codex.NewAdapter(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			vault := t.TempDir()
			cfg := Config{VaultRoot: vault, Namespace: "machine/agent-memory", Project: "gentle-ai-fork"}
			if _, err := Inject(home, tt.adapter, cfg); err != nil {
				t.Fatalf("Inject() error = %v", err)
			}
			body, err := os.ReadFile(tt.promptPath(home))
			if err != nil {
				t.Fatal(err)
			}
			text := string(body)
			if !strings.Contains(text, "<!-- gentle-ai:markdown-memory-protocol -->") {
				t.Fatalf("prompt missing markdown memory marker:\n%s", text)
			}
			if !strings.Contains(text, "machine/agent-memory/projects/gentle-ai-fork/index.md") {
				t.Fatalf("prompt missing rendered project recall path:\n%s", text)
			}
		})
	}
}
