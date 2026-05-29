package markdownmemory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

type Config struct {
	VaultRoot string
	Namespace string
	Project   string
}

type InjectionResult struct {
	Changed bool
	Files   []string
}

type templateFile struct {
	Path  string
	Asset string
}

func Inject(homeDir string, adapter agents.Adapter, cfg Config) (InjectionResult, error) {
	if err := ValidateConfig(cfg); err != nil {
		return InjectionResult{}, err
	}

	files := []string{}
	changed := false

	templateResult, err := EnsureTemplates(cfg)
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || templateResult.Changed
	files = append(files, templateResult.Files...)

	if !supportsProtocolInjection(adapter.Agent()) || !adapter.SupportsSystemPrompt() {
		return InjectionResult{Changed: changed, Files: files}, nil
	}

	promptPath := adapter.SystemPromptFile(homeDir)
	existing, err := readFileOrEmpty(promptPath)
	if err != nil {
		return InjectionResult{}, err
	}
	protocol := renderTemplate(assets.MustRead("markdown-memory/protocol.md"), cfg)
	updated := filemerge.InjectMarkdownSection(existing, "markdown-memory-protocol", protocol)
	writeResult, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || writeResult.Changed
	files = append(files, promptPath)

	return InjectionResult{Changed: changed, Files: files}, nil
}

func EnsureTemplates(cfg Config) (InjectionResult, error) {
	if err := ValidateConfig(cfg); err != nil {
		return InjectionResult{}, err
	}

	root := MemoryRoot(cfg)
	projectRoot := filepath.Join(root, "projects", cfg.Project)
	date := time.Now().Format("2006-01-02")
	files := []templateFile{
		{Path: filepath.Join(root, "README.md"), Asset: "markdown-memory/README.md"},
		{Path: filepath.Join(root, "hot.md"), Asset: "markdown-memory/hot.md"},
		{Path: filepath.Join(root, "index.md"), Asset: "markdown-memory/index.md"},
		{Path: filepath.Join(projectRoot, "index.md"), Asset: "markdown-memory/project-index.md"},
		{Path: filepath.Join(projectRoot, "current-state.md"), Asset: "markdown-memory/current-state.md"},
		{Path: filepath.Join(projectRoot, "decisions.md"), Asset: "markdown-memory/decisions.md"},
		{Path: filepath.Join(projectRoot, "architecture.md"), Asset: "markdown-memory/architecture.md"},
		{Path: filepath.Join(projectRoot, "risks.md"), Asset: "markdown-memory/risks.md"},
		{Path: filepath.Join(projectRoot, "open-questions.md"), Asset: "markdown-memory/open-questions.md"},
		{Path: filepath.Join(projectRoot, "handoffs", "initial-setup.md"), Asset: "markdown-memory/handoff.md"},
		{Path: filepath.Join(projectRoot, "sessions", date+"-initial-setup.active.md"), Asset: "markdown-memory/session.md"},
		{Path: filepath.Join(projectRoot, "inbox", "staged-observations.md"), Asset: "markdown-memory/staged-observations.md"},
		{Path: filepath.Join(root, "shared", "coding-preferences.md"), Asset: "markdown-memory/shared-coding-preferences.md"},
		{Path: filepath.Join(root, "shared", "security-standards.md"), Asset: "markdown-memory/shared-security-standards.md"},
		{Path: filepath.Join(root, "shared", "agent-protocol.md"), Asset: "markdown-memory/shared-agent-protocol.md"},
	}

	result := InjectionResult{Files: make([]string, 0, len(files))}
	for _, file := range files {
		content := renderTemplate(assets.MustRead(file.Asset), cfg)
		writeResult, err := writeFileIfMissing(file.Path, []byte(content))
		if err != nil {
			return InjectionResult{}, err
		}
		result.Changed = result.Changed || writeResult.Changed
		result.Files = append(result.Files, file.Path)
	}
	return result, nil
}

func Paths(cfg Config) []string {
	root := MemoryRoot(cfg)
	projectRoot := filepath.Join(root, "projects", cfg.Project)
	date := time.Now().Format("2006-01-02")
	return []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "hot.md"),
		filepath.Join(root, "index.md"),
		filepath.Join(projectRoot, "index.md"),
		filepath.Join(projectRoot, "current-state.md"),
		filepath.Join(projectRoot, "decisions.md"),
		filepath.Join(projectRoot, "architecture.md"),
		filepath.Join(projectRoot, "risks.md"),
		filepath.Join(projectRoot, "open-questions.md"),
		filepath.Join(projectRoot, "handoffs", "initial-setup.md"),
		filepath.Join(projectRoot, "sessions", date+"-initial-setup.active.md"),
		filepath.Join(projectRoot, "inbox", "staged-observations.md"),
		filepath.Join(root, "shared", "coding-preferences.md"),
		filepath.Join(root, "shared", "security-standards.md"),
		filepath.Join(root, "shared", "agent-protocol.md"),
	}
}

func MemoryRoot(cfg Config) string {
	return filepath.Join(cfg.VaultRoot, filepath.FromSlash(cfg.Namespace))
}

func ValidateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.VaultRoot) == "" {
		return fmt.Errorf("memory-vault is required for markdown memory")
	}
	if !filepath.IsAbs(cfg.VaultRoot) {
		return fmt.Errorf("memory-vault must be an absolute path")
	}
	namespace := filepath.ToSlash(filepath.Clean(strings.TrimSpace(cfg.Namespace)))
	if namespace == "." || namespace == "" || filepath.IsAbs(cfg.Namespace) {
		return fmt.Errorf("memory-namespace must be relative to memory-vault")
	}
	if namespace != "machine/agent-memory" && !strings.HasPrefix(namespace, "machine/agent-memory/") {
		return fmt.Errorf("memory-namespace must stay under machine/agent-memory")
	}
	if strings.TrimSpace(cfg.Project) == "" {
		return fmt.Errorf("memory-project is required for markdown memory")
	}
	return nil
}

func supportsProtocolInjection(agent model.AgentID) bool {
	return agent == model.AgentOpenCode || agent == model.AgentCodex
}

func renderTemplate(content string, cfg Config) string {
	replacer := strings.NewReplacer(
		"{{DATE}}", time.Now().Format("2006-01-02"),
		"{{VAULT}}", cfg.VaultRoot,
		"{{NAMESPACE}}", filepath.ToSlash(filepath.Clean(cfg.Namespace)),
		"{{PROJECT}}", cfg.Project,
	)
	return replacer.Replace(content)
}

func writeFileIfMissing(path string, content []byte) (filemerge.WriteResult, error) {
	if _, err := os.Lstat(path); err == nil {
		return filemerge.WriteResult{}, nil
	} else if !os.IsNotExist(err) {
		return filemerge.WriteResult{}, fmt.Errorf("stat memory file %q: %w", path, err)
	}
	return filemerge.WriteFileAtomic(path, content, 0o644)
}

func readFileOrEmpty(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(data), nil
}
